package call

import (
	"encoding/json"
	"fmt"
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

// CancelCall 取消通话（只有发起者可以取消）
func (h *CallHandler) CancelCall(ctx *gin.Context) {
	req, err := analysis.BindParameter[call.CancelCallRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.cancelCallLogic(ctx, &req)
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

func (h *CallHandler) cancelCallLogic(ctx *gin.Context, req *call.CancelCallRequest) (resp *call.CancelCallResponse, errCode code_msg.BusinessCode, err error) {
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

	// 验证是否为房间创建者（只有创建者可以取消）
	if room.CreatorID != userID {
		return nil, code_msg.NoPermission, nil
	}

	// 验证房间状态（只有呼叫中或响铃中的通话可以取消）
	if room.Status != models.RoomStatusCalling && room.Status != models.RoomStatusRinging {
		return nil, code_msg.BadRequest, nil
	}

	// 更新房间状态为已取消
	err = h.dao.UpdateRoomStatus(roomID, models.RoomStatusCancelled)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 取消通话超时定时器（已取消，不需要超时）
	if timeoutMgr := globalUtils.GetCallTimeoutManager(); timeoutMgr != nil {
		if manager, ok := timeoutMgr.(*websocket.CallTimeoutManager); ok {
			manager.CancelTimeout(roomID)
		}
	}

	// 获取所有参与者
	participants, err := h.dao.GetRoomParticipants(roomID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 更新所有参与者状态
	now := time.Now().Unix()
	for _, p := range participants {
		if p.UserID != userID {
			// 其他参与者标记为已错过
			if err := h.dao.UpdateParticipantStatus(roomID, p.UserID, models.ParticipantStatusMissed); err != nil {
				// 记录错误但不影响整体流程
				fmt.Printf("Failed to update participant %d status: %v\n", p.UserID, err)
			}
		}
	}

	// 创建通话记录
	for _, p := range participants {
		var direction string
		if p.UserID == room.CreatorID {
			direction = models.CallDirectionOutgoing
		} else {
			direction = models.CallDirectionIncoming
		}

		var status string
		if p.UserID == userID {
			status = models.RecordStatusCancelled
		} else {
			status = models.RecordStatusMissed
		}

		record := &models.CallRecord{
			RoomID:    roomID,
			UserID:    p.UserID,
			Type:      room.Type,
			CallType:  room.CallType,
			Direction: direction,
			Status:    status,
			CreatedAt: now,
		}

		// 如果是私聊，设置对方用户ID
		if room.CallType == models.CallScenePrivate {
			for _, other := range participants {
				if other.UserID != p.UserID {
					record.OtherUserID = other.UserID
					break
				}
			}
		}

		h.dao.CreateCallRecord(record)
	}

	// 通过WebSocket通知所有参与者
	wsHub := getWebSocketHub()
	if wsHub != nil {
		cancelMessage := map[string]interface{}{
			"type":       "call_cancel",
			"room_id":    roomID,
			"room_token": room.RoomToken,
			"user_id":    userID,
			"media_type": room.Type, // 传递媒体类型（voice/video）
			"call_type":  room.CallType, // 传递通话类型（private/group）
			"reason":     req.GetReason(),
		}

		cancelBytes, _ := json.Marshal(cancelMessage)

		// 通知所有参与者
		for _, p := range participants {
			if p.UserID != userID {
				wsHub.SendToUser(p.UserID, cancelBytes)
			}
		}
	}

	resp = &call.CancelCallResponse{
		Success: true,
	}

	return resp, 0, nil
}

