package finance

import (
	"gochat/api/admin/finance"
	"gochat/internal/pkg/analysis"

	"github.com/gin-gonic/gin"
)

// AdjustBalance 调整用户余额
func (c *FinanceController) AdjustBalance(ctx *gin.Context) {
	req, err := analysis.BindParameter[finance.AdjustBalanceRequest](ctx, c.response)
	if err != nil {
		return
	}

	// TODO: 实现余额调整逻辑
	// 1. 验证用户存在
	// 2. 更新余额
	// 3. 创建资金流水
	// 4. 记录日志
	_ = req // 暂时未使用

	c.response.JsonSuccess(ctx, &finance.AdjustBalanceResponse{
		Code:    0,
		Message: "调整成功",
	})
}

