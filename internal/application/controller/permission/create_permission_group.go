package permission

import (
	"encoding/json"
	"gochat/api/admin/permission"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// CreatePermissionGroup 创建权限组
func (c *PermissionController) CreatePermissionGroup(ctx *gin.Context) {
	req, err := analysis.BindParameter[permission.CreatePermissionGroupRequest](ctx, c.response)
	if err != nil {
		return
	}

	permissionsJSON, _ := json.Marshal(req.GetPermissions())

	now := time.Now().Unix()
	group := &models.PermissionGroup{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Permissions: string(permissionsJSON),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := c.permissionDao.CreatePermissionGroup(group); err != nil {
		c.response.JsonError(ctx, err, err.Error())
		return
	}

	// 记录日志
	adminID, err := utils.GetCurrentAdminID(ctx)
	if err != nil {
		// 如果获取管理员ID失败，记录日志但不影响主流程
		adminID = 0
	}
	c.adminLogDao.CreateLog(&models.AdminLog{
		AdminID:     adminID,
		ActionType:  models.ActionTypePermissionCreate,
		Description: "创建权限组: " + group.Name,
		TargetType:  "permission",
		TargetID:    group.ID,
		CreatedAt:   now,
	})

	var permissions []string
	json.Unmarshal([]byte(group.Permissions), &permissions)

	resp := &permission.CreatePermissionGroupResponse{
		Group: &permission.PermissionGroupInfo{
			Id:          group.ID,
			Name:        group.Name,
			Description: group.Description,
			Permissions: permissions,
			IsDefault:   group.IsDefault,
			CreatedAt:   group.CreatedAt,
			UpdatedAt:   group.UpdatedAt,
		},
	}

	c.response.JsonSuccess(ctx, resp)
}

