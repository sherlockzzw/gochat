package balance

import (
	"gochat/api/api/balance"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// Withdraw 提现申请
func (h *BalanceHandler) Withdraw(ctx *gin.Context) {
	req, err := analysis.BindParameter[balance.WithdrawRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.withdrawLogic(ctx, &req)
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

func (h *BalanceHandler) withdrawLogic(ctx *gin.Context, req *balance.WithdrawRequest) (resp *balance.WithdrawResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 验证金额范围：1-2000元 (100-200000分)
	amount := req.GetAmount()
	if amount < 100 || amount > 200000 {
		return nil, code_msg.BadRequest, nil
	}

	// 检查今日提现次数（每日限3次）
	todayCount, err := h.dao.GetTodayWithdrawCount(userID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if todayCount >= 3 {
		return nil, code_msg.BadRequest, nil // TODO: 定义专门的错误码
	}

	// 检查余额是否充足
	userBalance, err := h.dao.GetBalance(userID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if userBalance.Balance < amount {
		return nil, code_msg.BadRequest, nil // TODO: 定义余额不足错误码
	}

	// 计算手续费（这里先预留，实际应该从配置中读取）
	fee := int64(0) // TODO: 从系统配置中读取手续费规则

	// 创建提现申请
	now := time.Now().Unix()
	withdrawReq := &models.WithdrawRequest{
		UserID:     userID,
		Amount:     amount,
		Fee:        fee,
		AccountInfo: req.GetAccountInfo(), // TODO: 加密存储
		Status:     models.StatusPending,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	err = h.dao.CreateWithdrawRequest(withdrawReq)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return &balance.WithdrawResponse{
		RequestId: withdrawReq.ID,
		Success:   true,
	}, 0, nil
}

// GetWithdrawRequests 获取提现申请列表
func (h *BalanceHandler) GetWithdrawRequests(ctx *gin.Context) {
	req, err := analysis.BindQuery[balance.GetWithdrawRequestsRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.getWithdrawRequestsLogic(ctx, &req)
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

func (h *BalanceHandler) getWithdrawRequestsLogic(ctx *gin.Context, req *balance.GetWithdrawRequestsRequest) (resp *balance.GetWithdrawRequestsResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 获取提现申请列表
	withdrawRequests, total, err := h.dao.GetWithdrawRequests(userID, int(req.GetPage()), int(req.GetPageSize()))
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	requestInfos := make([]*balance.WithdrawRequestInfo, 0, len(withdrawRequests))
	for _, r := range withdrawRequests {
		requestInfos = append(requestInfos, &balance.WithdrawRequestInfo{
			Id:        r.ID,
			Amount:    r.Amount,
			Fee:       r.Fee,
			Status:    r.Status,
			AuditorId: r.AuditorID,
			AuditTime: r.AuditTime,
			AuditNote: r.AuditNote,
			CreatedAt: r.CreatedAt,
		})
	}

	return &balance.GetWithdrawRequestsResponse{
		Requests:    requestInfos,
		TotalCount:  int32(total),
		CurrentPage: req.GetPage(),
		PageSize:    req.GetPageSize(),
	}, 0, nil
}

