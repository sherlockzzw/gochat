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
	"time"

	"github.com/gin-gonic/gin"
)

// RejectCall 拒绝通话
func (h *CallHandler) RejectCall(ctx *gin.Context) {
	req, err := analysis.BindParameter[call.RejectCallRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.rejectCallLogic(ctx, &req)
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

func (h *CallHandler) rejectCallLogic(ctx *gin.Context, req *call.RejectCallRequest) (resp *call.RejectCallResponse, errCode code_msg.BusinessCode, err error) {
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

	// 更新参与者状态为已拒绝
	err = h.dao.UpdateParticipantStatus(roomID, userID, models.ParticipantStatusRejected)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 更新房间状态为已结束
	err = h.dao.UpdateRoomStatus(roomID, models.RoomStatusEnded)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 取消通话超时定时器（已拒绝，不需要超时）
	if timeoutMgr := globalUtils.GetCallTimeoutManager(); timeoutMgr != nil {
		if manager, ok := timeoutMgr.(*websocket.CallTimeoutManager); ok {
			manager.CancelTimeout(roomID)
		}
	}

	// 创建通话记录
	now := time.Now().Unix()
	record := &models.CallRecord{
		RoomID:    roomID,
		UserID:    userID,
		Type:      room.Type,
		CallType:  room.CallType,
		Direction: models.CallDirectionIncoming,
		Status:    models.RecordStatusRejected,
		CreatedAt: now,
	}

	// 如果是私聊，设置对方用户ID
	if room.CallType == models.CallScenePrivate {
		participants, _ := h.dao.GetRoomParticipants(roomID)
		for _, p := range participants {
			if p.UserID != userID {
				record.OtherUserID = p.UserID
				break
			}
		}
	}

	h.dao.CreateCallRecord(record)

	// 获取所有参与者，通知他们通话已拒绝
	participants, err := h.dao.GetRoomParticipants(roomID)
	if err == nil {
		wsHub := getWebSocketHub()
		if wsHub != nil {
			rejectMessage := map[string]interface{}{
				"type":       "call_reject",
				"room_id":    roomID,
				"room_token": room.RoomToken,
				"user_id":    userID,
				"media_type": room.Type, // 传递媒体类型（voice/video）
				"call_type":  room.CallType, // 传递通话类型（private/group）
				"reason":     req.GetReason(),
			}

			rejectBytes, _ := json.Marshal(rejectMessage)

			// 通知所有参与者
			for _, p := range participants {
				if p.UserID != userID {
					wsHub.SendToUser(p.UserID, rejectBytes)
				}
			}
		}
	}

	resp = &call.RejectCallResponse{
		Success: true,
	}

	return resp, 0, nil
}

