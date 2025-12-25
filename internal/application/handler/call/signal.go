package call

import (
	"encoding/json"
	"fmt"
	"gochat/internal/application/handler/common"
	"gochat/internal/infrastructure/models"
	"log"
)

// ProcessCallSignal 处理来自WebSocket的通话信令
// 这个方法由WebSocket层调用，用于验证和转发信令
func (h *CallHandler) ProcessCallSignal(signalData []byte, fromUserID int64) error {
	var signal map[string]interface{}
	if err := json.Unmarshal(signalData, &signal); err != nil {
		return err
	}

	signalType, ok := signal["type"].(string)
	if !ok {
		return nil
	}

	roomIDFloat, ok := signal["room_id"].(float64)
	if !ok {
		return nil
	}
	roomID := int64(roomIDFloat)

	// 验证房间和用户权限
	room, participants, err := h.validateCallRoom(roomID, fromUserID)
	if err != nil {
		log.Printf("Failed to validate call room: %v", err)
		return err
	}

	// 设置发送者ID
	signal["from_user_id"] = float64(fromUserID)

	// 根据信令类型处理
	switch signalType {
	case "call_accept", "call_reject", "call_cancel", "call_end",
		"call_joined", "call_left",
		"call_mute", "call_unmute",
		"call_video_on", "call_video_off":
		// 广播给房间内其他参与者
		h.broadcastToRoomParticipants(room, participants, signal, fromUserID)

	case "offer", "answer", "ice_candidate":
		// WebRTC信令
		h.forwardWebRTCSignal(room, participants, signal, fromUserID)
	}

	return nil
}

// validateCallRoom 验证通话房间和用户权限
func (h *CallHandler) validateCallRoom(roomID, userID int64) (*models.CallRoom, []*models.CallParticipant, error) {
	// 获取房间信息
	room, err := h.dao.GetRoomByID(roomID)
	if err != nil {
		return nil, nil, err
	}

	// 获取参与者列表
	participants, err := h.dao.GetRoomParticipants(roomID)
	if err != nil {
		return nil, nil, err
	}

	// 验证用户是否为参与者
	isParticipant := false
	for _, p := range participants {
		if p.UserID == userID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		return nil, nil, fmt.Errorf("user is not a participant")
	}

	return room, participants, nil
}

// broadcastToRoomParticipants 向房间内所有其他参与者广播信令
func (h *CallHandler) broadcastToRoomParticipants(room *models.CallRoom, participants []*models.CallParticipant, signal map[string]interface{}, excludeUserID int64) {
	// 序列化信令消息
	signalBytes, err := json.Marshal(signal)
	if err != nil {
		log.Printf("Failed to marshal signal: %v", err)
		return
	}

	// 发送给所有其他参与者
	for _, p := range participants {
		if p.UserID != excludeUserID {
			common.SendToUser(p.UserID, signalBytes)
		}
	}
}

// forwardWebRTCSignal 转发WebRTC信令
func (h *CallHandler) forwardWebRTCSignal(room *models.CallRoom, participants []*models.CallParticipant, signal map[string]interface{}, fromUserID int64) {
	// 序列化信令消息
	signalBytes, err := json.Marshal(signal)
	if err != nil {
		log.Printf("Failed to marshal WebRTC signal: %v", err)
		return
	}

	// 如果有指定接收者，发送给该用户
	if toUserIDFloat, ok := signal["to_user_id"].(float64); ok && toUserIDFloat > 0 {
		common.SendToUser(int64(toUserIDFloat), signalBytes)
		return
	}

	// 如果是私聊，转发给对方
	if room.CallType == models.CallScenePrivate {
		for _, p := range participants {
			if p.UserID != fromUserID {
				common.SendToUser(p.UserID, signalBytes)
				return
			}
		}
	} else {
		// 群聊：转发给所有其他参与者
		for _, p := range participants {
			if p.UserID != fromUserID {
				common.SendToUser(p.UserID, signalBytes)
			}
		}
	}
}

