package statistics

import (
	"gochat/api/admin/statistics"
	"gochat/internal/pkg/analysis"

	"github.com/gin-gonic/gin"
)

// GetTerminalStatistics 获取终端统计
// 功能: 统计各终端类型的活跃用户数，包括 iOS、Android、手机网页、PC桌面、PC网页等
// 参数: 通过 Query 参数传递 start_time 和 end_time（Unix 时间戳，秒）
// 返回: 终端统计数据，key 为终端类型名称，value 为活跃用户数
func (c *StatisticsController) GetTerminalStatistics(ctx *gin.Context) {
	req, err := analysis.BindQuery[statistics.GetTerminalStatisticsRequest](ctx, c.response)
	if err != nil {
		return
	}

	startTime := req.GetStartTime()
	endTime := req.GetEndTime()

	terminalData, err := c.deviceDao.GetTerminalStatistics(startTime, endTime)
	if err != nil {
		c.response.JsonError(ctx, err, "获取终端统计失败")
		return
	}

	resp := &statistics.GetTerminalStatisticsResponse{
		TerminalData: terminalData,
	}

	c.response.JsonSuccess(ctx, resp)
}

