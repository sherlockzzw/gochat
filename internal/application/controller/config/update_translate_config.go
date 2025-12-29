package config

import (
	"gochat/api/admin/config"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"time"

	"github.com/gin-gonic/gin"
)

// UpdateTranslateConfig 更新易翻译配置
func (c *ConfigController) UpdateTranslateConfig(ctx *gin.Context) {
	req, err := analysis.BindParameter[config.UpdateTranslateConfigRequest](ctx, c.response)
	if err != nil {
		return
	}

	c.configDao.SetConfigValue(models.ConfigKeyTranslateEnabled, req.GetEnabled(), "易翻译全局启用")
	c.configDao.SetConfigValue(models.ConfigKeyTranslateDailyLimit, req.GetDailyLimit(), "每日调用次数限制")
	c.configDao.SetConfig(models.ConfigKeyTranslateApiKey, req.GetApiKey(), "API密钥")
	c.configDao.SetConfigValue(models.ConfigKeyTranslateLanguages, req.GetLanguages(), "支持的语言列表")

	// 记录日志
	c.adminLogDao.CreateLog(&models.AdminLog{
		AdminID:     0, // TODO
		ActionType:  models.ActionTypeConfigUpdate,
		Description: "更新易翻译配置",
		TargetType:  "config",
		CreatedAt:   time.Now().Unix(),
	})

	c.response.JsonSuccess(ctx, &config.UpdateTranslateConfigResponse{
		Code:    0,
		Message: "更新成功",
	})
}

