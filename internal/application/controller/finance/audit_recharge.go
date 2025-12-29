package finance

import (
	"gochat/api/admin/finance"
	"gochat/internal/pkg/analysis"

	"github.com/gin-gonic/gin"
)

// AuditRecharge 审核充值
func (c *FinanceController) AuditRecharge(ctx *gin.Context) {
	req, err := analysis.BindParameter[finance.AuditRechargeRequest](ctx, c.response)
	if err != nil {
		return
	}

	// TODO: 实现审核逻辑
	// 1. 获取充值申请
	// 2. 更新状态
	// 3. 如果通过，增加用户余额
	// 4. 创建资金流水
	// 5. 记录日志
	_ = req // 暂时未使用

	c.response.JsonSuccess(ctx, &finance.AuditRechargeResponse{
		Code:    0,
		Message: "审核成功",
	})
}

