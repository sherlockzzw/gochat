package statistics

import (
	"gochat/api/admin/statistics"
	"gochat/internal/pkg/analysis"

	"github.com/gin-gonic/gin"
)

// GetFinanceStatistics 获取资金流水统计
func (c *StatisticsController) GetFinanceStatistics(ctx *gin.Context) {
	_, err := analysis.BindQuery[statistics.GetFinanceStatisticsRequest](ctx, c.response)
	if err != nil {
		return
	}

	// TODO: 实现资金流水统计
	// 从BalanceFlow表中统计各类资金流水

	resp := &statistics.GetFinanceStatisticsResponse{
		TotalRecharge: 0,
		TotalWithdraw: 0,
		TotalRedpacket: 0,
		TotalTransfer: 0,
	}

	c.response.JsonSuccess(ctx, resp)
}

