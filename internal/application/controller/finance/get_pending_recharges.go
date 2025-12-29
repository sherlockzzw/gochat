package finance

import (
	"gochat/api/admin/finance"
	"gochat/internal/pkg/analysis"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetPendingRecharges 获取待审核充值
func (c *FinanceController) GetPendingRecharges(ctx *gin.Context) {
	req, err := analysis.BindQuery[finance.GetPendingRechargesRequest](ctx, c.response)
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

	recharges, total, err := c.balanceDao.GetPendingRecharges(page, pageSize)
	if err != nil {
		c.response.JsonError(ctx, err, err.Error())
		return
	}

	var rechargeInfos []*finance.RechargeInfo
	for _, r := range recharges {
		user, _ := c.userDao.GetUserByID(r.UserID)
		userName := ""
		if user != nil {
			userName = user.Name
		}

		rechargeInfos = append(rechargeInfos, &finance.RechargeInfo{
			Id:        r.ID,
			UserId:    r.UserID,
			UserName:  userName,
			Amount:    r.Amount,
			Status:    r.Status,
			Remark:    r.Remark,
			CreatedAt: timestamppb.New(time.Unix(r.CreatedAt, 0)),
		})
	}

	resp := &finance.GetPendingRechargesResponse{
		Recharges: rechargeInfos,
		Total:     int32(total),
	}

	c.response.JsonSuccess(ctx, resp)
}

