package permission

import (
	"encoding/json"
	"gochat/api/admin/permission"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// UpdatePermissionGroup 更新权限组
func (c *PermissionController) UpdatePermissionGroup(ctx *gin.Context) {
	req, err := analysis.BindParameter[permission.UpdatePermissionGroupRequest](ctx, c.response)
	if err != nil {
		return
	}

	// 获取权限组
	group, err := c.permissionDao.GetPermissionGroupByID(req.GetId())
	if err != nil {
		c.response.JsonError(ctx, err, err.Error())
		return
	}
	if group == nil {
		c.response.JsonErrorFixation(ctx, code_msg.NotFound)
		return
	}

	// 更新字段
	permissionsJSON, _ := json.Marshal(req.GetPermissions())
	group.Name = req.GetName()
	group.Description = req.GetDescription()
	group.Permissions = string(permissionsJSON)
	group.UpdatedAt = time.Now().Unix()

	// 保存更新
	if err := c.permissionDao.UpdatePermissionGroup(group); err != nil {
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
		ActionType:  models.ActionTypePermissionUpdate,
		Description: "更新权限组: " + group.Name,
		TargetType:  "permission",
		TargetID:    group.ID,
		CreatedAt:   time.Now().Unix(),
	})

	c.response.JsonSuccess(ctx, &permission.UpdatePermissionGroupResponse{})
}
