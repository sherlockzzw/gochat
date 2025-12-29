package finance

import (
	"gochat/api/admin/finance"
	"gochat/internal/pkg/analysis"

	"github.com/gin-gonic/gin"
)

// AuditWithdraw 审核提现
func (c *FinanceController) AuditWithdraw(ctx *gin.Context) {
	req, err := analysis.BindParameter[finance.AuditWithdrawRequest](ctx, c.response)
	if err != nil {
		return
	}

	// TODO: 实现审核逻辑
	_ = req // 暂时未使用

	c.response.JsonSuccess(ctx, &finance.AuditWithdrawResponse{
		Code:    0,
		Message: "审核成功",
	})
}

