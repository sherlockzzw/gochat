package call

import (
	"encoding/json"
	"gochat/api/api/call"
	"gochat/internal/infrastructure/models"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
	globalUtils "gochat/utils"

	"github.com/gin-gonic/gin"
)

// AcceptCall 接受通话
func (h *CallHandler) AcceptCall(ctx *gin.Context) {
	req, err := analysis.BindParameter[call.AcceptCallRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.acceptCallLogic(ctx, &req)
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

func (h *CallHandler) acceptCallLogic(ctx *gin.Context, req *call.AcceptCallRequest) (resp *call.AcceptCallResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	roomID := int64(req.GetRoomId())
	if roomID <= 0 {
		return nil, code_msg.BadRequest, nil
	}

	// 获取房间信息
	room, err := h.dao.GetRoomByID(roomID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if room == nil {
		return nil, code_msg.NotFound, nil
	}

	// 验证用户是否为参与者
	participant, err := h.dao.GetParticipant(roomID, userID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if participant == nil {
		return nil, code_msg.Unauthorized, nil
	}

	// 验证房间状态
	if room.Status != models.RoomStatusCalling && room.Status != models.RoomStatusRinging {
		return nil, code_msg.BadRequest, nil
	}

	// 更新参与者状态为已加入
	err = h.dao.UpdateParticipantStatus(roomID, userID, models.ParticipantStatusJoined)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 更新房间状态为已连接
	err = h.dao.UpdateRoomStatus(roomID, models.RoomStatusConnected)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 取消通话超时定时器（已接听，不需要超时）
	if timeoutMgr := globalUtils.GetCallTimeoutManager(); timeoutMgr != nil {
		if manager, ok := timeoutMgr.(*websocket.CallTimeoutManager); ok {
			manager.CancelTimeout(roomID)
		}
	}

	// 获取所有参与者，通知他们通话已连接
	participants, err := h.dao.GetRoomParticipants(roomID)
	if err == nil {
		wsHub := getWebSocketHub()
		if wsHub != nil {
			acceptMessage := map[string]interface{}{
				"type":       "call_accept",
				"room_id":    roomID,
				"room_token": room.RoomToken,
				"user_id":    userID,
			}

			acceptBytes, _ := json.Marshal(acceptMessage)

			// 通知所有参与者
			for _, p := range participants {
				if p.UserID != userID {
					wsHub.SendToUser(p.UserID, acceptBytes)
				}
			}
		}
	}

	resp = &call.AcceptCallResponse{
		Success: true,
	}

	return resp, 0, nil
}

