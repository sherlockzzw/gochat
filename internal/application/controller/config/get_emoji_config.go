package config

import (
	"gochat/api/admin/config"
	"gochat/internal/infrastructure/models"

	"github.com/gin-gonic/gin"
)

// GetEmojiConfig 获取表情包配置
func (c *ConfigController) GetEmojiConfig(ctx *gin.Context) {
	var allowCustom bool
	var maxSize, maxCount int32

	c.configDao.GetConfigValue(models.ConfigKeyEmojiAllowCustom, &allowCustom)
	c.configDao.GetConfigValue(models.ConfigKeyEmojiMaxSize, &maxSize)
	c.configDao.GetConfigValue(models.ConfigKeyEmojiMaxCount, &maxCount)

	resp := &config.GetEmojiConfigResponse{
		AllowCustom: allowCustom,
		MaxSize:     maxSize,
		MaxCount:    maxCount,
	}

	c.response.JsonSuccess(ctx, resp)
}

