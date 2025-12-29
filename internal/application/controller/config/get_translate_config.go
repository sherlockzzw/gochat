package config

import (
	"gochat/api/admin/config"
	"gochat/internal/infrastructure/models"

	"github.com/gin-gonic/gin"
)

// GetTranslateConfig 获取易翻译配置
func (c *ConfigController) GetTranslateConfig(ctx *gin.Context) {
	var enabled bool
	var dailyLimit int32
	var apiKey string
	var languages []string

	c.configDao.GetConfigValue(models.ConfigKeyTranslateEnabled, &enabled)
	c.configDao.GetConfigValue(models.ConfigKeyTranslateDailyLimit, &dailyLimit)
	
	cfg, _ := c.configDao.GetConfig(models.ConfigKeyTranslateApiKey)
	if cfg != nil {
		apiKey = cfg.ConfigValue
	}
	
	c.configDao.GetConfigValue(models.ConfigKeyTranslateLanguages, &languages)

	resp := &config.GetTranslateConfigResponse{
		Enabled:    enabled,
		DailyLimit: dailyLimit,
		ApiKey:     apiKey,
		Languages:  languages,
	}

	c.response.JsonSuccess(ctx, resp)
}

