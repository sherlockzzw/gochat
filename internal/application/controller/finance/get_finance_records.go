package finance

import (
	"gochat/api/admin/finance"
	"gochat/internal/pkg/analysis"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetFinanceRecords 获取资金流水
func (c *FinanceController) GetFinanceRecords(ctx *gin.Context) {
	req, err := analysis.BindQuery[finance.GetFinanceRecordsRequest](ctx, c.response)
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

	flows, total, err := c.balanceDao.GetAllBalanceFlows(
		req.GetUserId(),
		req.GetType(),
		req.GetStartTime(),
		req.GetEndTime(),
		page,
		pageSize,
	)
	if err != nil {
		c.response.JsonError(ctx, err, err.Error())
		return
	}

	var records []*finance.FinanceRecord
	for _, f := range flows {
		records = append(records, &finance.FinanceRecord{
			Id:        f.ID,
			UserId:    f.UserID,
			Type:      f.Type,
			Amount:    f.Amount,
			Balance:   f.Balance,
			RelatedId: f.RelatedID,
			Remark:    f.Remark,
			CreatedAt: timestamppb.New(time.Unix(f.CreatedAt, 0)),
		})
	}

	resp := &finance.GetFinanceRecordsResponse{
		Records: records,
		Total:   int32(total),
	}

	c.response.JsonSuccess(ctx, resp)
}

