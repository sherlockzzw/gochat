package websocket

import (
	"context"
	"encoding/json"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/infrastructure/models"
	"log"
	"sync"
	"time"
)

// CallTimeoutManager 通话超时管理器
type CallTimeoutManager struct {
	timeouts    map[int64]*time.Timer // roomID -> timer
	mutex       sync.RWMutex
	timeoutDuration time.Duration     // 超时时间（默认60秒）
	dao         *dao.CallDao
	wsHub       *Hub
	ctx         context.Context
	cancel      context.CancelFunc
}

var globalCallTimeoutManager *CallTimeoutManager

// NewCallTimeoutManager 创建通话超时管理器
func NewCallTimeoutManager(timeoutSeconds int, callDao *dao.CallDao, wsHub *Hub) *CallTimeoutManager {
	if timeoutSeconds <= 0 {
		timeoutSeconds = 60 // 默认60秒
	}

	ctx, cancel := context.WithCancel(context.Background())

	manager := &CallTimeoutManager{
		timeouts:        make(map[int64]*time.Timer),
		timeoutDuration: time.Duration(timeoutSeconds) * time.Second,
		dao:             callDao,
		wsHub:           wsHub,
		ctx:             ctx,
		cancel:          cancel,
	}

	return manager
}

// SetCallTimeoutManager 设置全局超时管理器
func SetCallTimeoutManager(manager *CallTimeoutManager) {
	globalCallTimeoutManager = manager
}

// GetCallTimeoutManager 获取全局超时管理器
func GetCallTimeoutManager() *CallTimeoutManager {
	return globalCallTimeoutManager
}

// StartTimeout 启动通话超时定时器
func (m *CallTimeoutManager) StartTimeout(roomID int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 如果已有定时器，先取消
	if timer, exists := m.timeouts[roomID]; exists {
		timer.Stop()
	}

	// 创建新的定时器
	timer := time.AfterFunc(m.timeoutDuration, func() {
		m.handleTimeout(roomID)
	})

	m.timeouts[roomID] = timer
	log.Printf("Call timeout started for room %d, duration: %v", roomID, m.timeoutDuration)
}

// CancelTimeout 取消通话超时定时器
func (m *CallTimeoutManager) CancelTimeout(roomID int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if timer, exists := m.timeouts[roomID]; exists {
		timer.Stop()
		delete(m.timeouts, roomID)
		log.Printf("Call timeout cancelled for room %d", roomID)
	}
}

// handleTimeout 处理超时
func (m *CallTimeoutManager) handleTimeout(roomID int64) {
	log.Printf("Call timeout triggered for room %d", roomID)

	// 获取房间信息
	room, err := m.dao.GetRoomByID(roomID)
	if err != nil || room == nil {
		log.Printf("Failed to get room %d: %v", roomID, err)
		return
	}

	// 检查房间状态（只有呼叫中或响铃中的通话才需要超时处理）
	if room.Status != models.RoomStatusCalling && room.Status != models.RoomStatusRinging {
		log.Printf("Room %d is not in calling/ringing state, skip timeout", roomID)
		m.CancelTimeout(roomID)
		return
	}

	// 更新房间状态为已结束（超时）
	err = m.dao.UpdateRoomStatus(roomID, models.RoomStatusEnded)
	if err != nil {
		log.Printf("Failed to update room status: %v", err)
		return
	}

	// 获取所有参与者
	participants, err := m.dao.GetRoomParticipants(roomID)
	if err != nil {
		log.Printf("Failed to get participants: %v", err)
		return
	}

	// 更新参与者状态
	now := time.Now().Unix()
	for _, p := range participants {
		if p.Status == models.ParticipantStatusInvited || p.Status == models.ParticipantStatusRinging {
			// 未接听的参与者标记为已错过
			m.dao.UpdateParticipantStatus(roomID, p.UserID, models.ParticipantStatusMissed)
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
		if p.Status == models.ParticipantStatusJoined {
			// 已加入的参与者（理论上不应该发生，但为了安全）
			status = models.RecordStatusAnswered
		} else {
			// 未接听的参与者
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

		m.dao.CreateCallRecord(record)
	}

	// 通过WebSocket通知所有参与者
	timeoutMessage := map[string]interface{}{
		"type":       "call_timeout",
		"room_id":    roomID,
		"room_token": room.RoomToken,
		"reason":     "通话超时，对方未接听",
	}

	timeoutBytes, _ := json.Marshal(timeoutMessage)

	// 通知所有参与者
	for _, p := range participants {
		m.wsHub.SendToUser(p.UserID, timeoutBytes)
	}

	// 清理定时器
	m.CancelTimeout(roomID)

	log.Printf("Call timeout handled for room %d", roomID)
}

// Stop 停止超时管理器
func (m *CallTimeoutManager) Stop() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 取消所有定时器
	for roomID, timer := range m.timeouts {
		timer.Stop()
		delete(m.timeouts, roomID)
	}

	// 取消上下文
	m.cancel()

	log.Println("Call timeout manager stopped")
}

// Stats 获取统计信息
func (m *CallTimeoutManager) Stats() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return map[string]interface{}{
		"active_timeouts": len(m.timeouts),
		"timeout_duration": m.timeoutDuration.String(),
	}
}

