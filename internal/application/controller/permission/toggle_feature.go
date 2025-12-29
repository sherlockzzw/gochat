package permission

import (
	"gochat/api/admin/permission"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	utils_jwt "gochat/internal/pkg/utils"
	"gochat/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// ToggleFeature 功能开关
func (c *PermissionController) ToggleFeature(ctx *gin.Context) {
	req, err := analysis.BindParameter[permission.ToggleFeatureRequest](ctx, c.response)
	if err != nil {
		return
	}

	// 使用SystemConfigDao更新配置
	configDao := dao.NewSystemConfigDao(utils.DB)
	enabled := "false"
	if req.GetEnabled() {
		enabled = "true"
	}
	configDao.SetConfig(req.GetFeatureKey(), enabled, "功能开关")

	// 记录日志
	adminID, err := utils_jwt.GetCurrentAdminID(ctx)
	if err != nil {
		// 如果获取管理员ID失败，记录日志但不影响主流程
		adminID = 0
	}
	c.adminLogDao.CreateLog(&models.AdminLog{
		AdminID:     adminID,
		ActionType:  models.ActionTypeFeatureToggle,
		Description: "功能开关: " + req.GetFeatureKey() + " = " + enabled,
		TargetType:  "config",
		CreatedAt:   time.Now().Unix(),
	})

	c.response.JsonSuccess(ctx, &permission.ToggleFeatureResponse{})
}

