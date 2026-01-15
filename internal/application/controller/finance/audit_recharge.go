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

// AuditRecharge 审核充值
func (c *FinanceController) AuditRecharge(ctx *gin.Context) {
	req, err := analysis.BindParameter[finance.AuditRechargeRequest](ctx, c.response)
	if err != nil {
		return
	}

	adminID, err := utils.GetCurrentAdminID(ctx)
	if err != nil {
		c.response.JsonError(ctx, err, err.Error())
		return
	}

	if req.GetId() <= 0 {
		c.response.JsonSuccess(ctx, &finance.AuditRechargeResponse{Code: 400, Message: "参数错误"})
		return
	}
	if req.GetAction() != "approve" && req.GetAction() != "reject" {
		c.response.JsonSuccess(ctx, &finance.AuditRechargeResponse{Code: 400, Message: "参数错误"})
		return
	}
	if req.GetAction() == "reject" && req.GetNote() == "" {
		c.response.JsonSuccess(ctx, &finance.AuditRechargeResponse{Code: 400, Message: "请填写驳回原因"})
		return
	}

	now := time.Now().Unix()

	err = c.balanceDao.GetDB().Transaction(func(tx *gorm.DB) error {
		var r models.RechargeRequest
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", req.GetId()).First(&r).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("充值申请不存在")
			}
			return err
		}

		if r.Status != models.StatusPending {
			return fmt.Errorf("该充值申请已处理")
		}

		updates := map[string]interface{}{
			"auditor_id": adminID,
			"audit_time":  now,
			"audit_note":  req.GetNote(),
			"updated_at":  now,
		}

		if req.GetAction() == "reject" {
			updates["status"] = models.StatusRejected
			return tx.Model(&models.RechargeRequest{}).Where("id = ?", r.ID).Updates(updates).Error
		}

		updates["status"] = models.StatusApproved
		if err := tx.Model(&models.RechargeRequest{}).Where("id = ?", r.ID).Updates(updates).Error; err != nil {
			return err
		}

		userBalance, err := c.balanceDao.GetBalanceForUpdate(tx, r.UserID)
		if err != nil {
			return err
		}
		newBalance := userBalance.Balance + r.Amount
		if err := tx.Model(&models.UserBalance{}).Where("user_id = ?", r.UserID).Update("balance", newBalance).Error; err != nil {
			return err
		}

		flow := &models.BalanceFlow{
			UserID:    r.UserID,
			Type:      models.FlowTypeRecharge,
			Amount:    r.Amount,
			Balance:   newBalance,
			RelatedID: r.ID,
			Remark:    "充值审核通过",
			CreatedAt: now,
		}
		if err := tx.Create(flow).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		c.response.JsonSuccess(ctx, &finance.AuditRechargeResponse{Code: 500, Message: err.Error()})
		return
	}

	_ = c.adminLogDao.CreateLog(&models.AdminLog{
		AdminID:     adminID,
		ActionType:  models.ActionTypeFinanceAudit,
		Description: fmt.Sprintf("审核充值申请 id=%d action=%s", req.GetId(), req.GetAction()),
		TargetType:  "finance",
		TargetID:    req.GetId(),
		IPAddress:   ctx.ClientIP(),
		UserAgent:   ctx.GetHeader("User-Agent"),
		CreatedAt:   now,
	})

	c.response.JsonSuccess(ctx, &finance.AuditRechargeResponse{Code: 0, Message: "审核成功"})
}

