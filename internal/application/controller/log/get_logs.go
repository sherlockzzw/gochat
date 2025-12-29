package log

import (
	"gochat/api/admin/log"
	"gochat/internal/pkg/analysis"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetLogs 获取操作日志
func (c *LogController) GetLogs(ctx *gin.Context) {
	req, err := analysis.BindQuery[log.GetLogsRequest](ctx, c.response)
	if err != nil {
		return
	}

	page := int(req.GetPage())
	if page <= 0 {
		page = 1
	}
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 20
	}

	logs, total, err := c.adminLogDao.GetLogs(
		page,
		pageSize,
		req.GetActionType(),
		req.GetTargetType(),
		req.GetStartTime(),
		req.GetEndTime(),
	)
	if err != nil {
		c.response.JsonError(ctx, err, err.Error())
		return
	}

	var logInfos []*log.LogInfo
	for _, l := range logs {
		logInfos = append(logInfos, &log.LogInfo{
			Id:         l.ID,
			AdminId:    l.AdminID,
			AdminName:  l.AdminName,
			ActionType: l.ActionType,
			Description: l.Description,
			TargetType: l.TargetType,
			TargetId:   l.TargetID,
			IpAddress:  l.IPAddress,
			CreatedAt:  timestamppb.New(time.Unix(l.CreatedAt, 0)),
		})
	}

	resp := &log.GetLogsResponse{
		Logs: logInfos,
		Total: int32(total),
	}

	c.response.JsonSuccess(ctx, resp)
}

