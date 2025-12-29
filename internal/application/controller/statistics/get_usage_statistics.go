package statistics

import (
	"gochat/api/admin/statistics"
	"gochat/internal/pkg/analysis"

	"github.com/gin-gonic/gin"
)

// GetUsageStatistics 获取功能使用统计
func (c *StatisticsController) GetUsageStatistics(ctx *gin.Context) {
	_, err := analysis.BindQuery[statistics.GetUsageStatisticsRequest](ctx, c.response)
	if err != nil {
		return
	}

	// TODO: 实现功能使用统计
	// 可以从各个业务表中统计使用次数

	resp := &statistics.GetUsageStatisticsResponse{
		UsageData: make(map[string]int64),
	}

	c.response.JsonSuccess(ctx, resp)
}

