package balance

import (
	"gochat/api/api/balance"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"

	"github.com/gin-gonic/gin"
)

// GetRedPacketDetail 获取红包详情
func (h *BalanceHandler) GetRedPacketDetail(ctx *gin.Context) {
	req, err := analysis.BindQuery[balance.GetRedPacketDetailRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.getRedPacketDetailLogic(ctx, &req)
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

func (h *BalanceHandler) getRedPacketDetailLogic(ctx *gin.Context, req *balance.GetRedPacketDetailRequest) (resp *balance.GetRedPacketDetailResponse, errCode code_msg.BusinessCode, err error) {
	redPacketID := req.GetRedPacketId()

	// 获取红包信息
	redPacket, err := h.dao.GetRedPacket(redPacketID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if redPacket == nil {
		return nil, code_msg.NotFound, nil
	}

	// 获取领取记录（群聊时使用）
	var receiveInfos []*balance.RedPacketReceiveInfo
	if redPacket.Type == "group" {
		receives, err := h.dao.GetRedPacketReceives(redPacketID)
		if err != nil {
			return nil, code_msg.ServerError, err
		}
		receiveInfos = make([]*balance.RedPacketReceiveInfo, 0, len(receives))
		for _, r := range receives {
			receiveInfos = append(receiveInfos, &balance.RedPacketReceiveInfo{
				UserId:     r.UserID,
				Amount:     r.Amount,
				ReceivedAt: r.ReceivedAt,
			})
		}
	}

	// 转换红包分配类型
	var redPacketType balance.RedPacketType
	if redPacket.RedPacketType == models.RedPacketAllocTypeNormal {
		redPacketType = balance.RedPacketType_NORMAL
	} else {
		redPacketType = balance.RedPacketType_LUCKY
	}

	detail := &balance.RedPacketDetail{
		Id:             redPacket.ID,
		Type:           redPacket.Type,
		RedPacketType:  redPacketType,
		SenderId:       redPacket.SenderID,
		ReceiverId:     redPacket.ReceiverID,
		GroupId:        redPacket.GroupID,
		Amount:         redPacket.Amount,
		TotalAmount:    redPacket.TotalAmount,
		Count:          int32(redPacket.Count),
		RemainingAmount: redPacket.RemainingAmount,
		RemainingCount:  int32(redPacket.RemainingCount),
		Status:         redPacket.Status,
		Message:        redPacket.Message,
		ExpiredAt:      redPacket.ExpiredAt,
		CreatedAt:      redPacket.CreatedAt,
		Receives:       receiveInfos,
	}

	return &balance.GetRedPacketDetailResponse{
		Detail: detail,
	}, 0, nil
}

