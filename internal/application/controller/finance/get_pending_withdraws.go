package finance

import (
	"gochat/api/admin/finance"
	"gochat/internal/pkg/analysis"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetPendingWithdraws 获取待审核提现
func (c *FinanceController) GetPendingWithdraws(ctx *gin.Context) {
	req, err := analysis.BindQuery[finance.GetPendingWithdrawsRequest](ctx, c.response)
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

	withdraws, total, err := c.balanceDao.GetPendingWithdraws(page, pageSize)
	if err != nil {
		c.response.JsonError(ctx, err, err.Error())
		return
	}

	var withdrawInfos []*finance.WithdrawInfo
	for _, w := range withdraws {
		user, _ := c.userDao.GetUserByID(w.UserID)
		userName := ""
		if user != nil {
			userName = user.Name
		}

		withdrawInfos = append(withdrawInfos, &finance.WithdrawInfo{
			Id:        w.ID,
			UserId:    w.UserID,
			UserName:  userName,
			Amount:    w.Amount,
			Fee:       w.Fee,
			Status:    w.Status,
			CreatedAt: timestamppb.New(time.Unix(w.CreatedAt, 0)),
		})
	}

	resp := &finance.GetPendingWithdrawsResponse{
		Withdraws: withdrawInfos,
		Total:     int32(total),
	}

	c.response.JsonSuccess(ctx, resp)
}

