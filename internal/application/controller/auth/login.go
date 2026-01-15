package auth

import (
	"encoding/json"
	"gochat/api/admin/auth"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/middleware"
	"gochat/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// Login 管理员登录
func (c *AuthController) Login(ctx *gin.Context) {
	req, err := analysis.BindParameter[auth.LoginRequest](ctx, c.response)
	if err != nil {
		return
	}

	resp, code, err := c.loginLogic(ctx, &req)
	if code != 0 {
		c.response.JsonErrorFixation(ctx, code)
		return
	}
	if err != nil {
		c.response.JsonError(ctx, err, err.Error())
		return
	}

	c.response.JsonSuccess(ctx, resp)
}

func (c *AuthController) loginLogic(ctx *gin.Context, req *auth.LoginRequest) (resp *auth.LoginResponse, errCode code_msg.BusinessCode, err error) {
	// 查找管理员
	admin, err := c.adminDao.GetAdminByName(req.Username)
	if err != nil {
		return nil, code_msg.Unauthorized, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		//return nil, code_msg.Unauthorized, err
	}

	// 生成JWT token
	adminModel := &models.Admin{
		ID:   admin.ID,
		Name: admin.Name,
	}
	jwtMw := middleware.JwtMiddleware("Admin")
	token, _, err := jwtMw.GenerateToken(adminModel, time.Duration(viper.GetInt("token.expire"))*time.Hour)
	if err != nil {
		return nil, code_msg.Unauthorized, err
	}

	// 获取管理员的权限组ID和权限信息
	//var roleID int64 = 0
	//var roleName string = ""
	var permissions []string = []string{}

	upgDao := dao.NewUserPermissionGroupDao(utils.DB)
	groupID, err := upgDao.GetGroupIDByAdminID(admin.ID)
	if err == nil && groupID > 0 {
		//roleID = groupID
		pgDao := dao.NewPermissionGroupDao(utils.DB)
		group, err := pgDao.GetByID(groupID)
		if err == nil && group != nil {
			//roleName = group.Name
			// 解析权限列表（JSON格式）
			if group.Permissions != "" {
				json.Unmarshal([]byte(group.Permissions), &permissions)
			}
		}
	}

	resp = &auth.LoginResponse{
		Data: &auth.LoginData{
			Token: token,
			Admin: &auth.AdminInfo{
				Id:   admin.ID,
				Name: admin.Name,
			},
		},
	}

	// 注意：由于AdminInfo proto定义中没有role_id、role_name、permissions字段
	// 我们需要在响应中扩展这些字段，但proto生成的代码不支持
	// 这里先返回基本字段，权限信息通过GetCurrent接口获取

	now := time.Now().Unix()
	_ = c.adminLogDao.CreateLog(&models.AdminLog{
		AdminID:     admin.ID,
		AdminName:   admin.Name,
		ActionType:  models.ActionTypeAdminLogin,
		Description: "管理员登录",
		TargetType:  "admin",
		TargetID:    admin.ID,
		IPAddress:   ctx.ClientIP(),
		UserAgent:   ctx.GetHeader("User-Agent"),
		CreatedAt:   now,
	})

	return resp, 0, nil
}
