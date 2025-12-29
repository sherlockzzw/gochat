package finance

import (
	"gochat/api/admin/finance"
	"gochat/internal/infrastructure/models"

	"github.com/gin-gonic/gin"
)

// GetFinanceConfig 获取资金配置
func (c *FinanceController) GetFinanceConfig(ctx *gin.Context) {
	var rechargeEnabled, withdrawEnabled bool
	var rechargeDailyLimit, withdrawDailyLimit, withdrawFeeRate int64
	var withdrawDailyCount int32

	c.configDao.GetConfigValue(models.ConfigKeyFinanceRechargeEnabled, &rechargeEnabled)
	c.configDao.GetConfigValue(models.ConfigKeyFinanceWithdrawEnabled, &withdrawEnabled)
	c.configDao.GetConfigValue(models.ConfigKeyFinanceRechargeDailyLimit, &rechargeDailyLimit)
	c.configDao.GetConfigValue(models.ConfigKeyFinanceWithdrawDailyLimit, &withdrawDailyLimit)
	c.configDao.GetConfigValue(models.ConfigKeyFinanceWithdrawDailyCount, &withdrawDailyCount)
	c.configDao.GetConfigValue(models.ConfigKeyFinanceWithdrawFeeRate, &withdrawFeeRate)

	resp := &finance.GetFinanceConfigResponse{
		RechargeEnabled:    rechargeEnabled,
		WithdrawEnabled:    withdrawEnabled,
		RechargeDailyLimit: rechargeDailyLimit,
		WithdrawDailyLimit: withdrawDailyLimit,
		WithdrawDailyCount: withdrawDailyCount,
		WithdrawFeeRate:    withdrawFeeRate,
	}

	c.response.JsonSuccess(ctx, resp)
}

