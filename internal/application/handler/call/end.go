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

// EndCall 结束通话
func (h *CallHandler) EndCall(ctx *gin.Context) {
	req, err := analysis.BindParameter[call.EndCallRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.endCallLogic(ctx, &req)
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

func (h *CallHandler) endCallLogic(ctx *gin.Context, req *call.EndCallRequest) (resp *call.EndCallResponse, errCode code_msg.BusinessCode, err error) {
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

	// 结束房间（更新所有参与者状态）
	err = h.dao.EndRoom(roomID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 取消通话超时定时器（已结束，不需要超时）
	if timeoutMgr := globalUtils.GetCallTimeoutManager(); timeoutMgr != nil {
		if manager, ok := timeoutMgr.(*websocket.CallTimeoutManager); ok {
			manager.CancelTimeout(roomID)
		}
	}

	// 重新获取房间信息（获取更新后的时长）
	room, _ = h.dao.GetRoomByID(roomID)

	// 创建通话记录
	now := time.Now().Unix()
	participants, _ := h.dao.GetRoomParticipants(roomID)

	for _, p := range participants {
		var otherUserID int64
		var direction string

		// 确定对方用户ID和方向
		if room.CallType == models.CallScenePrivate {
			for _, other := range participants {
				if other.UserID != p.UserID {
					otherUserID = other.UserID
					break
				}
			}
			if p.UserID == room.CreatorID {
				direction = models.CallDirectionOutgoing
			} else {
				direction = models.CallDirectionIncoming
			}
		} else {
			direction = models.CallDirectionOutgoing
		}

		// 确定状态
		var status string
		if p.Status == models.ParticipantStatusJoined {
			status = models.RecordStatusAnswered
		} else if p.Status == models.ParticipantStatusRejected {
			status = models.RecordStatusRejected
		} else {
			status = models.RecordStatusMissed
		}

		record := &models.CallRecord{
			RoomID:      roomID,
			UserID:      p.UserID,
			Type:        room.Type,
			CallType:    room.CallType,
			OtherUserID: otherUserID,
			GroupID:     room.GroupID,
			Direction:   direction,
			Status:      status,
			Duration:    room.Duration,
			StartedAt:   room.StartedAt,
			EndedAt:     room.EndedAt,
			CreatedAt:   now,
		}

		h.dao.CreateCallRecord(record)
	}

	// 通知所有参与者通话已结束
	wsHub := getWebSocketHub()
	if wsHub != nil {
		endMessage := map[string]interface{}{
			"type":       "call_end",
			"room_id":    roomID,
			"room_token": room.RoomToken,
			"user_id":    userID,
			"media_type": room.Type, // 传递媒体类型（voice/video）
			"call_type":  room.CallType, // 传递通话类型（private/group）
			"duration":   room.Duration,
		}

		endBytes, _ := json.Marshal(endMessage)

		// 通知所有参与者
		for _, p := range participants {
			if p.UserID != userID {
				wsHub.SendToUser(p.UserID, endBytes)
			}
		}
	}

	resp = &call.EndCallResponse{
		Success: true,
	}

	return resp, 0, nil
}
