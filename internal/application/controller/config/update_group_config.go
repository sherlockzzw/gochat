package config

import (
	"gochat/api/admin/config"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"time"

	"github.com/gin-gonic/gin"
)

// UpdateGroupConfig 更新群组配置
func (c *ConfigController) UpdateGroupConfig(ctx *gin.Context) {
	req, err := analysis.BindParameter[config.UpdateGroupConfigRequest](ctx, c.response)
	if err != nil {
		return
	}

	c.configDao.SetConfigValue(models.ConfigKeyGroupAllowCreate, req.GetAllowCreate(), "允许创建群组")
	c.configDao.SetConfigValue(models.ConfigKeyGroupRedPacketEnabled, req.GetRedpacketEnabled(), "群红包启用")
	c.configDao.SetConfigValue(models.ConfigKeyGroupRedPacketLimit, req.GetRedpacketLimit(), "群红包每日次数上限")

	// 记录日志
	c.adminLogDao.CreateLog(&models.AdminLog{
		AdminID:     0, // TODO
		ActionType:  models.ActionTypeConfigUpdate,
		Description: "更新群组配置",
		TargetType:  "config",
		CreatedAt:   time.Now().Unix(),
	})

	c.response.JsonSuccess(ctx, &config.UpdateGroupConfigResponse{
		Code:    0,
		Message: "更新成功",
	})
}

