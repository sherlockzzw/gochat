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

// Recharge 充值申请
func (h *BalanceHandler) Recharge(ctx *gin.Context) {
	req, err := analysis.BindParameter[balance.RechargeRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.rechargeLogic(ctx, &req)
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

func (h *BalanceHandler) rechargeLogic(ctx *gin.Context, req *balance.RechargeRequest) (resp *balance.RechargeResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 验证金额范围：0.01-5000元 (1-500000分)
	amount := req.GetAmount()
	if amount < 1 || amount > 500000 {
		return nil, code_msg.BadRequest, nil
	}

	// 创建充值申请
	now := time.Now().Unix()
	rechargeReq := &models.RechargeRequest{
		UserID:    userID,
		Amount:    amount,
		Status:    models.StatusPending,
		Remark:    req.GetRemark(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = h.dao.CreateRechargeRequest(rechargeReq)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return &balance.RechargeResponse{
		RequestId: rechargeReq.ID,
		Success:   true,
	}, 0, nil
}

// GetRechargeRequests 获取充值申请列表
func (h *BalanceHandler) GetRechargeRequests(ctx *gin.Context) {
	req, err := analysis.BindQuery[balance.GetRechargeRequestsRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.getRechargeRequestsLogic(ctx, &req)
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

func (h *BalanceHandler) getRechargeRequestsLogic(ctx *gin.Context, req *balance.GetRechargeRequestsRequest) (resp *balance.GetRechargeRequestsResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 获取充值申请列表
	requests, total, err := h.dao.GetRechargeRequests(userID, int(req.GetPage()), int(req.GetPageSize()))
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 转换为响应格式
	requestInfos := make([]*balance.RechargeRequestInfo, 0, len(requests))
	for _, r := range requests {
		requestInfos = append(requestInfos, &balance.RechargeRequestInfo{
			Id:         r.ID,
			Amount:     r.Amount,
			Status:     r.Status,
			Remark:     r.Remark,
			AuditorId:  r.AuditorID,
			AuditTime:  r.AuditTime,
			AuditNote:  r.AuditNote,
			CreatedAt:  r.CreatedAt,
		})
	}

	return &balance.GetRechargeRequestsResponse{
		Requests:    requestInfos,
		TotalCount:  int32(total),
		CurrentPage: req.GetPage(),
		PageSize:    req.GetPageSize(),
	}, 0, nil
}




