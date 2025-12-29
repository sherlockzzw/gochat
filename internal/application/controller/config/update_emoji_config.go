package config

import (
	"gochat/api/admin/config"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"time"

	"github.com/gin-gonic/gin"
)

// UpdateEmojiConfig 更新表情包配置
func (c *ConfigController) UpdateEmojiConfig(ctx *gin.Context) {
	req, err := analysis.BindParameter[config.UpdateEmojiConfigRequest](ctx, c.response)
	if err != nil {
		return
	}

	c.configDao.SetConfigValue(models.ConfigKeyEmojiAllowCustom, req.GetAllowCustom(), "允许用户自定义表情包")
	c.configDao.SetConfigValue(models.ConfigKeyEmojiMaxSize, req.GetMaxSize(), "表情包最大文件大小(MB)")
	c.configDao.SetConfigValue(models.ConfigKeyEmojiMaxCount, req.GetMaxCount(), "用户最大表情包数量")

	// 记录日志
	c.adminLogDao.CreateLog(&models.AdminLog{
		AdminID:     0, // TODO
		ActionType:  models.ActionTypeConfigUpdate,
		Description: "更新表情包配置",
		TargetType:  "config",
		CreatedAt:   time.Now().Unix(),
	})

	c.response.JsonSuccess(ctx, &config.UpdateEmojiConfigResponse{
		Code:    0,
		Message: "更新成功",
	})
}

