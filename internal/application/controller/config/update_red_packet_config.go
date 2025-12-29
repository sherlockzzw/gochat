package config

import (
	"gochat/api/admin/config"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"time"

	"github.com/gin-gonic/gin"
)

// UpdateRedPacketConfig 更新红包配置
func (c *ConfigController) UpdateRedPacketConfig(ctx *gin.Context) {
	req, err := analysis.BindParameter[config.UpdateRedPacketConfigRequest](ctx, c.response)
	if err != nil {
		return
	}

	c.configDao.SetConfigValue(models.ConfigKeyRedPacketPrivateMin, req.GetPrivateMin(), "私聊红包最小金额")
	c.configDao.SetConfigValue(models.ConfigKeyRedPacketPrivateMax, req.GetPrivateMax(), "私聊红包最大金额")
	c.configDao.SetConfigValue(models.ConfigKeyRedPacketPrivateCount, req.GetPrivateCount(), "私聊红包个数上限")
	c.configDao.SetConfigValue(models.ConfigKeyRedPacketGroupMin, req.GetGroupMin(), "群聊红包最小金额")
	c.configDao.SetConfigValue(models.ConfigKeyRedPacketGroupMax, req.GetGroupMax(), "群聊红包最大金额")
	c.configDao.SetConfigValue(models.ConfigKeyRedPacketGroupCount, req.GetGroupCount(), "群聊红包个数上限")

	// 记录日志
	c.adminLogDao.CreateLog(&models.AdminLog{
		AdminID:     0, // TODO
		ActionType:  models.ActionTypeConfigUpdate,
		Description: "更新红包配置",
		TargetType:  "config",
		CreatedAt:   time.Now().Unix(),
	})

	c.response.JsonSuccess(ctx, &config.UpdateRedPacketConfigResponse{
		Code:    0,
		Message: "更新成功",
	})
}

