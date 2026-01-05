package balance

import (
	"gochat/api/api/balance"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetBalanceFlows 获取资金流水
func (h *BalanceHandler) GetBalanceFlows(ctx *gin.Context) {
	req, err := analysis.BindQuery[balance.GetBalanceFlowsRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.getBalanceFlowsLogic(ctx, &req)
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

func (h *BalanceHandler) getBalanceFlowsLogic(ctx *gin.Context, req *balance.GetBalanceFlowsRequest) (resp *balance.GetBalanceFlowsResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	flowType := req.GetFlowType()

	// 获取资金流水
	flows, total, err := h.dao.GetBalanceFlows(userID, flowType, int(req.GetPage()), int(req.GetPageSize()))
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 转换为响应格式
	flowInfos := make([]*balance.BalanceFlowInfo, 0, len(flows))
	for _, f := range flows {
		flowInfos = append(flowInfos, &balance.BalanceFlowInfo{
			Id:        f.ID,
			Type:      f.Type,
			Amount:    f.Amount,
			Balance:   f.Balance,
			RelatedId: f.RelatedID,
			Remark:    f.Remark,
			CreatedAt: f.CreatedAt,
		})
	}

	return &balance.GetBalanceFlowsResponse{
		Flows:       flowInfos,
		TotalCount:  int32(total),
		CurrentPage: req.GetPage(),
		PageSize:    req.GetPageSize(),
	}, 0, nil
}




