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

// AuditWithdraw 审核提现
func (c *FinanceController) AuditWithdraw(ctx *gin.Context) {
	req, err := analysis.BindParameter[finance.AuditWithdrawRequest](ctx, c.response)
	if err != nil {
		return
	}

	adminID, err := utils.GetCurrentAdminID(ctx)
	if err != nil {
		c.response.JsonError(ctx, err, err.Error())
		return
	}

	if req.GetId() <= 0 {
		c.response.JsonSuccess(ctx, &finance.AuditWithdrawResponse{Code: 400, Message: "参数错误"})
		return
	}
	if req.GetAction() != "approve" && req.GetAction() != "reject" {
		c.response.JsonSuccess(ctx, &finance.AuditWithdrawResponse{Code: 400, Message: "参数错误"})
		return
	}
	if req.GetAction() == "reject" && req.GetNote() == "" {
		c.response.JsonSuccess(ctx, &finance.AuditWithdrawResponse{Code: 400, Message: "请填写驳回原因"})
		return
	}

	now := time.Now().Unix()

	err = c.balanceDao.GetDB().Transaction(func(tx *gorm.DB) error {
		var w models.WithdrawRequest
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", req.GetId()).First(&w).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("提现申请不存在")
			}
			return err
		}

		if w.Status != models.StatusPending {
			return fmt.Errorf("该提现申请已处理")
		}

		updates := map[string]interface{}{
			"auditor_id": adminID,
			"audit_time":  now,
			"audit_note":  req.GetNote(),
			"updated_at":  now,
		}

		if req.GetAction() == "approve" {
			updates["status"] = models.StatusApproved
			if err := tx.Model(&models.WithdrawRequest{}).Where("id = ?", w.ID).Updates(updates).Error; err != nil {
				return err
			}
			return nil
		}

		updates["status"] = models.StatusRejected
		if err := tx.Model(&models.WithdrawRequest{}).Where("id = ?", w.ID).Updates(updates).Error; err != nil {
			return err
		}

		userBalance, err := c.balanceDao.GetBalanceForUpdate(tx, w.UserID)
		if err != nil {
			return err
		}
		newBalance := userBalance.Balance + w.Amount
		if err := tx.Model(&models.UserBalance{}).Where("user_id = ?", w.UserID).Update("balance", newBalance).Error; err != nil {
			return err
		}

		flow := &models.BalanceFlow{
			UserID:    w.UserID,
			Type:      models.FlowTypeRefund,
			Amount:    w.Amount,
			Balance:   newBalance,
			RelatedID: w.ID,
			Remark:    "提现驳回退回",
			CreatedAt: now,
		}
		if err := tx.Create(flow).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		c.response.JsonSuccess(ctx, &finance.AuditWithdrawResponse{Code: 500, Message: err.Error()})
		return
	}

	_ = c.adminLogDao.CreateLog(&models.AdminLog{
		AdminID:     adminID,
		ActionType:  models.ActionTypeFinanceAudit,
		Description: fmt.Sprintf("审核提现申请 id=%d action=%s", req.GetId(), req.GetAction()),
		TargetType:  "finance",
		TargetID:    req.GetId(),
		IPAddress:   ctx.ClientIP(),
		UserAgent:   ctx.GetHeader("User-Agent"),
		CreatedAt:   now,
	})

	c.response.JsonSuccess(ctx, &finance.AuditWithdrawResponse{Code: 0, Message: "审核成功"})
}

