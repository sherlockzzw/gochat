package finance

import (
	"errors"
	"fmt"
	"gochat/api/admin/finance"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdjustBalance 调整用户余额
func (c *FinanceController) AdjustBalance(ctx *gin.Context) {
	req, err := analysis.BindParameter[finance.AdjustBalanceRequest](ctx, c.response)
	if err != nil {
		return
	}

	adminID, err := utils.GetCurrentAdminID(ctx)
	if err != nil {
		c.response.JsonError(ctx, err, err.Error())
		return
	}

	if req.GetUserId() <= 0 {
		c.response.JsonSuccess(ctx, &finance.AdjustBalanceResponse{Code: 400, Message: "参数错误"})
		return
	}
	if req.GetAmount() == 0 {
		c.response.JsonSuccess(ctx, &finance.AdjustBalanceResponse{Code: 400, Message: "调整金额不能为0"})
		return
	}
	if req.GetReason() == "" {
		c.response.JsonSuccess(ctx, &finance.AdjustBalanceResponse{Code: 400, Message: "请填写调整原因"})
		return
	}

	now := time.Now().Unix()

	user, err := c.userDao.GetUserByID(req.GetUserId())
	if err != nil {
		c.response.JsonError(ctx, err, err.Error())
		return
	}
	if user == nil {
		c.response.JsonSuccess(ctx, &finance.AdjustBalanceResponse{Code: 404, Message: "用户不存在"})
		return
	}

	err = c.balanceDao.GetDB().Transaction(func(tx *gorm.DB) error {
		b, err := c.balanceDao.GetBalanceForUpdate(tx, req.GetUserId())
		if err != nil {
			return err
		}

		newBalance := b.Balance + req.GetAmount()
		if newBalance < 0 {
			return fmt.Errorf("余额不足")
		}

		if err := tx.Model(&models.UserBalance{}).Where("user_id = ?", req.GetUserId()).Update("balance", newBalance).Error; err != nil {
			return err
		}

		flowType := ""
		if req.GetAmount() > 0 {
			flowType = models.FlowTypeRecharge
		} else {
			flowType = models.FlowTypeWithdraw
		}

		flow := &models.BalanceFlow{
			UserID:    req.GetUserId(),
			Type:      flowType,
			Amount:    req.GetAmount(),
			Balance:   newBalance,
			RelatedID: 0,
			Remark:    fmt.Sprintf("管理员调整余额: %s", req.GetReason()),
			CreatedAt: now,
		}
		return tx.Create(flow).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.response.JsonSuccess(ctx, &finance.AdjustBalanceResponse{Code: 404, Message: "用户余额记录不存在"})
			return
		}
		c.response.JsonSuccess(ctx, &finance.AdjustBalanceResponse{Code: 500, Message: err.Error()})
		return
	}

	_ = c.adminLogDao.CreateLog(&models.AdminLog{
		AdminID:     adminID,
		ActionType:  models.ActionTypeFinanceAdjust,
		Description: fmt.Sprintf("调整用户余额 user_id=%d amount=%d reason=%s", req.GetUserId(), req.GetAmount(), req.GetReason()),
		TargetType:  "user",
		TargetID:    req.GetUserId(),
		IPAddress:   ctx.ClientIP(),
		UserAgent:   ctx.GetHeader("User-Agent"),
		CreatedAt:   now,
	})

	c.response.JsonSuccess(ctx, &finance.AdjustBalanceResponse{Code: 0, Message: "调整成功"})
}

