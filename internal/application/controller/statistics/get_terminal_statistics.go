package statistics

import (
	"gochat/api/admin/statistics"
	"gochat/internal/pkg/analysis"

	"github.com/gin-gonic/gin"
)

// GetTerminalStatistics 获取终端统计
func (c *StatisticsController) GetTerminalStatistics(ctx *gin.Context) {
	_, err := analysis.BindQuery[statistics.GetTerminalStatisticsRequest](ctx, c.response)
	if err != nil {
		return
	}

	// TODO: 实现终端统计
	// 可以从UserDevice表中统计各终端活跃用户数

	resp := &statistics.GetTerminalStatisticsResponse{
		TerminalData: make(map[string]int64),
	}

	c.response.JsonSuccess(ctx, resp)
}

