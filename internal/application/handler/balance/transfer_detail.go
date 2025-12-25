package balance

import (
	"gochat/api/api/balance"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"

	"github.com/gin-gonic/gin"
)

// GetTransferDetail 获取转账详情
func (h *BalanceHandler) GetTransferDetail(ctx *gin.Context) {
	req, err := analysis.BindQuery[balance.GetTransferDetailRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.getTransferDetailLogic(ctx, &req)
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

func (h *BalanceHandler) getTransferDetailLogic(ctx *gin.Context, req *balance.GetTransferDetailRequest) (resp *balance.GetTransferDetailResponse, errCode code_msg.BusinessCode, err error) {
	transferID := req.GetTransferId()

	// 获取转账信息
	transfer, err := h.dao.GetTransfer(transferID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if transfer == nil {
		return nil, code_msg.NotFound, nil
	}

	detail := &balance.TransferDetail{
		Id:         transfer.ID,
		Type:       transfer.Type,
		SenderId:   transfer.SenderID,
		ReceiverId: transfer.ReceiverID,
		GroupId:    transfer.GroupID,
		Amount:     transfer.Amount,
		Status:     transfer.Status,
		Remark:     transfer.Remark,
		CreatedAt:  transfer.CreatedAt,
	}

	return &balance.GetTransferDetailResponse{
		Detail: detail,
	}, 0, nil
}

