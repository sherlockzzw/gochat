package statistics

import (
	"gochat/api/admin/statistics"
	"gochat/internal/pkg/analysis"

	"github.com/gin-gonic/gin"
)

// GetUsageStatistics 获取功能使用统计
// 功能: 统计各业务功能的使用次数，包括发送消息、创建群组、语音通话、视频通话、发送红包、转账等
// 参数: 通过 Query 参数传递 start_time 和 end_time（Unix 时间戳，秒）
// 返回: 功能使用统计数据，key 为功能名称，value 为使用次数
func (c *StatisticsController) GetUsageStatistics(ctx *gin.Context) {
	req, err := analysis.BindQuery[statistics.GetUsageStatisticsRequest](ctx, c.response)
	if err != nil {
		return
	}

	startTime := req.GetStartTime()
	endTime := req.GetEndTime()

	usageData, err := c.balanceDao.GetUsageStatistics(startTime, endTime)
	if err != nil {
		c.response.JsonError(ctx, err, "获取功能使用统计失败")
		return
	}

	resp := &statistics.GetUsageStatisticsResponse{
		UsageData: usageData,
	}

	c.response.JsonSuccess(ctx, resp)
}

