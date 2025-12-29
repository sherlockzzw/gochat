package permission

import (
	"gochat/api/admin/permission"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// AssignUserToGroup 分配用户到权限组
func (c *PermissionController) AssignUserToGroup(ctx *gin.Context) {
	req, err := analysis.BindParameter[permission.AssignUserToGroupRequest](ctx, c.response)
	if err != nil {
		return
	}

	if err := c.permissionDao.AssignUserToGroup(req.GetUserId(), req.GetGroupId()); err != nil {
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
		ActionType:  models.ActionTypePermissionAssign,
		Description: "分配用户到权限组",
		TargetType:  "user",
		TargetID:    req.GetUserId(),
		CreatedAt:   time.Now().Unix(),
	})

	c.response.JsonSuccess(ctx, &permission.AssignUserToGroupResponse{})
}

