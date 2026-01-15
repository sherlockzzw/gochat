package auth

import (
	"encoding/json"
	"gochat/internal/infrastructure/dao"
	utils_jwt "gochat/internal/pkg/utils"
	"gochat/utils"

	"github.com/gin-gonic/gin"
)

// GetCurrentAdminResponse 获取当前管理员信息响应结构
type GetCurrentAdminResponse struct {
	Admin struct {
		ID          int64    `json:"id"`
		Name        string   `json:"name"`
		RoleID      int64    `json:"role_id,omitempty"`
		RoleName    string   `json:"role_name,omitempty"`
		Permissions []string `json:"permissions,omitempty"`
	} `json:"admin"`
}

// GetCurrent 获取当前登录管理员信息（含权限信息）
func (c *AuthController) GetCurrent(ctx *gin.Context) {
	// 从JWT token中获取管理员ID
	adminID, err := utils_jwt.GetCurrentAdminID(ctx)
	if err != nil {
		c.response.JsonError(ctx, err, "获取管理员信息失败")
		return
	}

	// 获取管理员基本信息
	admin, err := c.adminDao.GetAdminByID(adminID)
	if err != nil {
		c.response.JsonError(ctx, err, "管理员不存在")
		return
	}

	// 获取管理员的权限组ID
	upgDao := dao.NewUserPermissionGroupDao(utils.DB)
	groupID, err := upgDao.GetGroupIDByAdminID(adminID)
	if err != nil {
		// 如果没有权限组，也不报错，只是groupID为0
		groupID = 0
	}

	var roleID int64 = 0
	var roleName string = ""
	var permissions []string = []string{}

	// 如果管理员有权限组，获取权限组信息和权限列表
	if groupID > 0 {
		roleID = groupID
		pgDao := dao.NewPermissionGroupDao(utils.DB)
		group, err := pgDao.GetByID(groupID)
		if err == nil && group != nil {
			roleName = group.Name
			// 解析权限列表（JSON格式）
			if group.Permissions != "" {
				json.Unmarshal([]byte(group.Permissions), &permissions)
			}
		}
	}

	// 构建响应
	var resp GetCurrentAdminResponse
	resp.Admin.ID = admin.ID
	resp.Admin.Name = admin.Name
	resp.Admin.RoleID = roleID
	resp.Admin.RoleName = roleName
	resp.Admin.Permissions = permissions

	c.response.JsonSuccess(ctx, resp)
}
