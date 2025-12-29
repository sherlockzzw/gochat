package config

import (
	"gochat/api/admin/config"
	"gochat/internal/infrastructure/models"

	"github.com/gin-gonic/gin"
)

// GetRedPacketConfig 获取红包配置
func (c *ConfigController) GetRedPacketConfig(ctx *gin.Context) {
	var privateMin, privateMax, groupMin, groupMax int64
	var privateCount, groupCount int32

	c.configDao.GetConfigValue(models.ConfigKeyRedPacketPrivateMin, &privateMin)
	c.configDao.GetConfigValue(models.ConfigKeyRedPacketPrivateMax, &privateMax)
	c.configDao.GetConfigValue(models.ConfigKeyRedPacketPrivateCount, &privateCount)
	c.configDao.GetConfigValue(models.ConfigKeyRedPacketGroupMin, &groupMin)
	c.configDao.GetConfigValue(models.ConfigKeyRedPacketGroupMax, &groupMax)
	c.configDao.GetConfigValue(models.ConfigKeyRedPacketGroupCount, &groupCount)

	resp := &config.GetRedPacketConfigResponse{
		PrivateMin:   privateMin,
		PrivateMax:   privateMax,
		PrivateCount: privateCount,
		GroupMin:     groupMin,
		GroupMax:     groupMax,
		GroupCount:   groupCount,
	}

	c.response.JsonSuccess(ctx, resp)
}

