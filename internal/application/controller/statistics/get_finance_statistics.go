package statistics

import (
	"gochat/api/admin/statistics"
	"gochat/internal/pkg/analysis"

	"github.com/gin-gonic/gin"
)

// GetFinanceStatistics 获取资金流水统计
func (c *StatisticsController) GetFinanceStatistics(ctx *gin.Context) {
	req, err := analysis.BindQuery[statistics.GetFinanceStatisticsRequest](ctx, c.response)
	if err != nil {
		return
	}

	startTime := req.GetStartTime()
	endTime := req.GetEndTime()

	totalRecharge, totalWithdraw, totalRedpacket, totalTransfer, err := c.balanceDao.GetFinanceStatistics(startTime, endTime)
	if err != nil {
		c.response.JsonError(ctx, err, "获取资金流水统计失败")
		return
	}

	resp := &statistics.GetFinanceStatisticsResponse{
		TotalRecharge:  totalRecharge,
		TotalWithdraw:  totalWithdraw,
		TotalRedpacket: totalRedpacket,
		TotalTransfer:  totalTransfer,
	}

	c.response.JsonSuccess(ctx, resp)
}

