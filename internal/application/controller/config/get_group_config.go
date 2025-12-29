package config

import (
	"gochat/api/admin/config"
	"gochat/internal/infrastructure/models"

	"github.com/gin-gonic/gin"
)

// GetGroupConfig 获取群组配置
func (c *ConfigController) GetGroupConfig(ctx *gin.Context) {
	var allowCreate, redpacketEnabled bool
	var redpacketLimit int32

	c.configDao.GetConfigValue(models.ConfigKeyGroupAllowCreate, &allowCreate)
	c.configDao.GetConfigValue(models.ConfigKeyGroupRedPacketEnabled, &redpacketEnabled)
	c.configDao.GetConfigValue(models.ConfigKeyGroupRedPacketLimit, &redpacketLimit)

	resp := &config.GetGroupConfigResponse{
		AllowCreate:    allowCreate,
		RedpacketEnabled: redpacketEnabled,
		RedpacketLimit: redpacketLimit,
	}

	c.response.JsonSuccess(ctx, resp)
}

