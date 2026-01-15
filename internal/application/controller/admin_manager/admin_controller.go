package admin_manager

import (
    admin_admin "gochat/api/admin/admin"
    "gochat/internal/infrastructure/dao"
    "gochat/internal/infrastructure/models"
    "gochat/internal/pkg/analysis"
    "gochat/internal/pkg/response"
    utils_jwt "gochat/internal/pkg/utils"
    "gochat/utils"
    "time"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

// AdminController 提供管理员相关 API
// 注意：管理员与权限组通过 user_permission_groups 进行关联，不在 admin 表做冗余 role_id

type AdminController struct {
    resp               *response.SvcRequest
    adminDao           *dao.AdminDao
    permissionGroupDao *dao.PermissionGroupDao
    upgDao             *dao.UserPermissionGroupDao
    adminLogDao        *dao.AdminLogDao
}

func NewAdminController() *AdminController {
    db := utils.DB
    return &AdminController{
        resp:               response.NewSvcRequest(),
        adminDao:           dao.NewAdminDao(db),
        permissionGroupDao: dao.NewPermissionGroupDao(db),
        upgDao:             dao.NewUserPermissionGroupDao(db),
        adminLogDao:        dao.NewAdminLogDao(db),
    }
}

// GetAdmins 列表
func (c *AdminController) GetAdmins(ctx *gin.Context) {
    req, err := analysis.BindQuery[admin_admin.GetAdminsRequest](ctx, c.resp)
    if err != nil {
        return
    }

    list, total, err := c.adminDao.ListWithPagination(int(req.GetPage()), int(req.GetPageSize()))
    if err != nil {
        c.resp.JsonError(ctx, err, "获取管理员失败")
        return
    }

    // 绑定权限组信息
    adminIDs := make([]int64, 0, len(list))
    for _, a := range list {
        adminIDs = append(adminIDs, a.ID)
    }
    idGroupMap, err := c.upgDao.GetGroupIDsByAdminIDs(adminIDs)
    if err != nil {
        c.resp.JsonError(ctx, err, "获取权限组失败")
        return
    }
    // 批量获取组名
    uniqueGroupIDs := make([]int64, 0)
    seen := map[int64]struct{}{}
    for _, gid := range idGroupMap {
        if gid == 0 {
            continue
        }
        if _, ok := seen[gid]; !ok {
            seen[gid] = struct{}{}
            uniqueGroupIDs = append(uniqueGroupIDs, gid)
        }
    }
    groups, _ := c.permissionGroupDao.List(uniqueGroupIDs)
    gidName := map[int64]string{}
    for _, g := range groups {
        gidName[g.ID] = g.Name
    }

    var adminInfos []*admin_admin.AdminInfo
    for _, a := range list {
        gid := idGroupMap[a.ID]
        adminInfos = append(adminInfos, &admin_admin.AdminInfo{
            Id:        a.ID,
            Name:      a.Name,
            RoleId:    gid,
            RoleName:  gidName[gid],
            CreatedAt: a.CreatedAt,
            UpdatedAt: a.UpdatedAt,
        })
    }

    resp := &admin_admin.GetAdminsResponse{
        Admins: adminInfos,
        Total:  int32(total),
    }
    c.resp.JsonSuccess(ctx, resp)
}

// CreateAdmin 新建
func (c *AdminController) CreateAdmin(ctx *gin.Context) {
    req, err := analysis.BindParameter[admin_admin.CreateAdminRequest](ctx, c.resp)
    if err != nil {
        return
    }
    if req.GetName() == "" || req.GetPassword() == "" {
        c.resp.JsonParamError(ctx, "name/password不能为空")
        return
    }
    // 重名检查
    if exist, _ := c.adminDao.GetAdminByName(req.GetName()); exist != nil {
        c.resp.JsonParamError(ctx, "管理员已存在")
        return
    }

    hashed, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
    if err != nil {
        c.resp.JsonError(ctx, err, "密码加密失败")
        return
    }

    now := time.Now().Unix()
    admin := &models.Admin{
        Name:      req.GetName(),
        Password:  string(hashed),
        CreatedAt: now,
        UpdatedAt: now,
    }

    if err := utils.DB.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(admin).Error; err != nil {
            return err
        }
        if req.GetRoleId() > 0 {
            upgDao := dao.NewUserPermissionGroupDao(tx)
            if err := upgDao.Bind(admin.ID, req.GetRoleId()); err != nil {
                return err
            }
        }
        return nil
    }); err != nil {
        c.resp.JsonError(ctx, err, "创建失败")
        return
    }

    roleName := ""
    if req.GetRoleId() > 0 {
        if pg, _ := c.permissionGroupDao.GetByID(req.GetRoleId()); pg != nil {
            roleName = pg.Name
        }
    }

    resp := &admin_admin.CreateAdminResponse{
        Code: 0,
        Message: "success",
        Admin: &admin_admin.AdminInfo{
            Id: admin.ID,
            Name: admin.Name,
            RoleId: req.GetRoleId(),
            RoleName: roleName,
            CreatedAt: admin.CreatedAt,
            UpdatedAt: admin.UpdatedAt,
        },
    }

    operatorID, err := utils_jwt.GetCurrentAdminID(ctx)
    if err != nil {
        operatorID = 0
    }
    _ = c.adminLogDao.CreateLog(&models.AdminLog{
        AdminID:     operatorID,
        ActionType:  models.ActionTypeAdminCreate,
        Description: "创建管理员: " + admin.Name,
        TargetType:  "admin",
        TargetID:    admin.ID,
        IPAddress:   ctx.ClientIP(),
        UserAgent:   ctx.GetHeader("User-Agent"),
        CreatedAt:   time.Now().Unix(),
    })

    c.resp.JsonSuccess(ctx, resp)
}

// UpdateAdmin 修改
func (c *AdminController) UpdateAdmin(ctx *gin.Context) {
    req, err := analysis.BindParameter[admin_admin.UpdateAdminRequest](ctx, c.resp)
    if err != nil {
        return
    }
    if req.GetId() <= 0 {
        c.resp.JsonParamError(ctx, "id无效")
        return
    }

    updates := map[string]interface{}{}
    if req.GetPassword() != "" {
        hashed, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
        if err != nil {
            c.resp.JsonError(ctx, err, "密码加密失败")
            return
        }
        updates["password"] = string(hashed)
    }
    if len(updates) == 0 && req.GetRoleId() == 0 {
        c.resp.JsonSuccess(ctx, &admin_admin.UpdateAdminResponse{Code:0, Message:"success"})
        return
    }
    updates["updated_at"] = time.Now().Unix()

    // transaction: update admin + permission group binding
    if err := utils.DB.Transaction(func(tx *gorm.DB) error {
        if len(updates) > 0 {
            if err := dao.NewAdminDao(tx).UpdateAdmin(req.GetId(), updates); err != nil {
                return err
            }
        }
        if req.GetRoleId() > 0 {
            if err := dao.NewUserPermissionGroupDao(tx).Bind(req.GetId(), req.GetRoleId()); err != nil {
                return err
            }
        }
        return nil
    }); err != nil {
        c.resp.JsonError(ctx, err, "更新失败")
        return
    }

    operatorID, err := utils_jwt.GetCurrentAdminID(ctx)
    if err != nil {
        operatorID = 0
    }
    _ = c.adminLogDao.CreateLog(&models.AdminLog{
        AdminID:     operatorID,
        ActionType:  models.ActionTypeAdminUpdate,
        Description: "更新管理员",
        TargetType:  "admin",
        TargetID:    req.GetId(),
        IPAddress:   ctx.ClientIP(),
        UserAgent:   ctx.GetHeader("User-Agent"),
        CreatedAt:   time.Now().Unix(),
    })

    c.resp.JsonSuccess(ctx, &admin_admin.UpdateAdminResponse{Code:0, Message:"success"})
}

// DeleteAdmin 删除
func (c *AdminController) DeleteAdmin(ctx *gin.Context) {
    req, err := analysis.BindParameter[admin_admin.DeleteAdminRequest](ctx, c.resp)
    if err != nil {
        return
    }
    if req.GetId() <= 0 {
        c.resp.JsonParamError(ctx, "invalid id")
        return
    }

    if err := utils.DB.Transaction(func(tx *gorm.DB) error {
        if err := dao.NewAdminDao(tx).DeleteAdmin(req.GetId()); err != nil {
            return err
        }
        if err := dao.NewUserPermissionGroupDao(tx).DeleteByAdminID(req.GetId()); err != nil {
            return err
        }
        return nil
    }); err != nil {
        c.resp.JsonError(ctx, err, "删除失败")
        return
    }

    operatorID, err := utils_jwt.GetCurrentAdminID(ctx)
    if err != nil {
        operatorID = 0
    }
    _ = c.adminLogDao.CreateLog(&models.AdminLog{
        AdminID:     operatorID,
        ActionType:  models.ActionTypeAdminDelete,
        Description: "删除管理员",
        TargetType:  "admin",
        TargetID:    req.GetId(),
        IPAddress:   ctx.ClientIP(),
        UserAgent:   ctx.GetHeader("User-Agent"),
        CreatedAt:   time.Now().Unix(),
    })

    c.resp.JsonSuccess(ctx, &admin_admin.DeleteAdminResponse{Code:0, Message:"success"})
}
