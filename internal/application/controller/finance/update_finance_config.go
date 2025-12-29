package finance

import (
	"gochat/api/admin/finance"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"time"

	"github.com/gin-gonic/gin"
)

// UpdateFinanceConfig 更新资金配置
func (c *FinanceController) UpdateFinanceConfig(ctx *gin.Context) {
	req, err := analysis.BindParameter[finance.UpdateFinanceConfigRequest](ctx, c.response)
	if err != nil {
		return
	}

	c.configDao.SetConfigValue(models.ConfigKeyFinanceRechargeEnabled, req.GetRechargeEnabled(), "充值启用")
	c.configDao.SetConfigValue(models.ConfigKeyFinanceWithdrawEnabled, req.GetWithdrawEnabled(), "提现启用")
	c.configDao.SetConfigValue(models.ConfigKeyFinanceRechargeDailyLimit, req.GetRechargeDailyLimit(), "单用户单日充值上限")
	c.configDao.SetConfigValue(models.ConfigKeyFinanceWithdrawDailyLimit, req.GetWithdrawDailyLimit(), "单用户单日提现上限")
	c.configDao.SetConfigValue(models.ConfigKeyFinanceWithdrawDailyCount, req.GetWithdrawDailyCount(), "每日提现次数限制")
	c.configDao.SetConfigValue(models.ConfigKeyFinanceWithdrawFeeRate, req.GetWithdrawFeeRate(), "提现手续费率")

	// 记录日志
	c.adminLogDao.CreateLog(&models.AdminLog{
		AdminID:     0, // TODO
		ActionType:  models.ActionTypeConfigUpdate,
		Description: "更新资金配置",
		TargetType:  "config",
		CreatedAt:   time.Now().Unix(),
	})

	c.response.JsonSuccess(ctx, &finance.UpdateFinanceConfigResponse{
		Code:    0,
		Message: "更新成功",
	})
}

