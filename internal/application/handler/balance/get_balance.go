package balance

import (
	"gochat/api/api/balance"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetBalance 获取余额
func (h *BalanceHandler) GetBalance(ctx *gin.Context) {
	resp, code, err := h.getBalanceLogic(ctx)
	if code != 0 {
		h.response.JsonErrorFixation(ctx, code)
		return
	}
	if err != nil {
		h.response.JsonError(ctx, err, err.Error())
		return
	}

	h.response.JsonSuccess(ctx, resp)
}

func (h *BalanceHandler) getBalanceLogic(ctx *gin.Context) (resp *balance.GetBalanceResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 获取余额
	userBalance, err := h.dao.GetBalance(userID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return &balance.GetBalanceResponse{
		Balance: userBalance.Balance,
	}, 0, nil
}

