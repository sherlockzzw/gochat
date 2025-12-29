package permission

import (
	"gochat/api/admin/permission"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// DeletePermissionGroup 删除权限组
func (c *PermissionController) DeletePermissionGroup(ctx *gin.Context) {
	req, err := analysis.BindParameter[permission.DeletePermissionGroupRequest](ctx, c.response)
	if err != nil {
		return
	}

	// 获取权限组信息（用于日志）
	group, err := c.permissionDao.GetPermissionGroupByID(req.GetId())
	if err != nil {
		c.response.JsonError(ctx, err, err.Error())
		return
	}
	if group == nil {
		c.response.JsonErrorFixation(ctx, code_msg.NotFound)
		return
	}

	// 检查是否是默认组（默认组不允许删除）
	if group.IsDefault {
		c.response.JsonErrorFixation(ctx, code_msg.ParameterError)
		return
	}

	// 检查是否有用户在使用该权限组
	userCount, err := c.permissionDao.GetGroupUserCount(req.GetId())
	if err != nil {
		c.response.JsonError(ctx, err, err.Error())
		return
	}
	if userCount > 0 {
		c.response.JsonErrorFixation(ctx, code_msg.ParameterError)
		return
	}

	// 删除权限组
	if err := c.permissionDao.DeletePermissionGroup(req.GetId()); err != nil {
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
		ActionType:  models.ActionTypePermissionDelete,
		Description: "删除权限组: " + group.Name,
		TargetType:  "permission",
		TargetID:    req.GetId(),
		CreatedAt:   time.Now().Unix(),
	})

	c.response.JsonSuccess(ctx, &permission.DeletePermissionGroupResponse{})
}

