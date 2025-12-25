package dao

import (
	"fmt"
	"gochat/internal/infrastructure/models"
	"time"

	"gorm.io/gorm"
)

type CallDao struct {
	db    *gorm.DB
	cache *CallCache
}

func NewCallDao(db *gorm.DB) *CallDao {
	return &CallDao{
		db:    db,
		cache: NewCallCache(),
	}
}

// GetDB 获取数据库连接（用于事务）
func (d *CallDao) GetDB() *gorm.DB {
	return d.db
}

// CreateRoom 创建通话房间
func (d *CallDao) CreateRoom(room *models.CallRoom) error {
	now := time.Now().Unix()
	room.CreatedAt = now
	room.UpdatedAt = now
	return d.db.Create(room).Error
}

// GetRoomByID 根据ID获取房间（带缓存）
func (d *CallDao) GetRoomByID(roomID int64) (*models.CallRoom, error) {
	// 先尝试从缓存获取
	if room, err := d.cache.GetRoomFromCache(roomID); err == nil {
		return room, nil
	}

	// 缓存未命中，从数据库查询
	var room models.CallRoom
	err := d.db.Where("id = ?", roomID).First(&room).Error
	if err != nil {
		return nil, err
	}

	// 异步更新缓存，不阻塞
	go d.cache.SetRoomToCache(&room)

	return &room, nil
}

// GetRoomByToken 根据Token获取房间
func (d *CallDao) GetRoomByToken(token string) (*models.CallRoom, error) {
	var room models.CallRoom
	err := d.db.Where("room_token = ?", token).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// UpdateRoomStatus 更新房间状态（同时更新缓存）
func (d *CallDao) UpdateRoomStatus(roomID int64, status string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now().Unix(),
	}

	if status == models.RoomStatusConnected {
		updates["started_at"] = time.Now().Unix()
	} else if status == models.RoomStatusEnded {
		updates["ended_at"] = time.Now().Unix()
		// 房间结束时，使缓存失效
		go d.cache.InvalidateRoomCache(roomID)
		// 计算通话时长
		var room models.CallRoom
		if err := d.db.Where("id = ?", roomID).First(&room).Error; err == nil && room.StartedAt > 0 {
			updates["duration"] = time.Now().Unix() - room.StartedAt
		}
	}

	return d.db.Model(&models.CallRoom{}).Where("id = ?", roomID).Updates(updates).Error
}

// CreateParticipant 创建参与者
func (d *CallDao) CreateParticipant(participant *models.CallParticipant) error {
	now := time.Now().Unix()
	participant.CreatedAt = now
	participant.UpdatedAt = now
	return d.db.Create(participant).Error
}

// GetParticipant 获取参与者
func (d *CallDao) GetParticipant(roomID, userID int64) (*models.CallParticipant, error) {
	var participant models.CallParticipant
	err := d.db.Where("room_id = ? AND user_id = ?", roomID, userID).First(&participant).Error
	if err != nil {
		return nil, err
	}
	return &participant, nil
}

// GetRoomParticipants 获取房间所有参与者（带缓存）
func (d *CallDao) GetRoomParticipants(roomID int64) ([]*models.CallParticipant, error) {
	// 先尝试从缓存获取
	if participants, err := d.cache.GetParticipantsFromCache(roomID); err == nil {
		return participants, nil
	}

	// 缓存未命中，从数据库查询
	var participants []*models.CallParticipant
	err := d.db.Where("room_id = ?", roomID).Find(&participants).Error
	if err != nil {
		return nil, err
	}

	// 异步更新缓存，不阻塞
	go d.cache.SetParticipantsToCache(roomID, participants)

	return participants, nil
}

// UpdateParticipantStatus 更新参与者状态
func (d *CallDao) UpdateParticipantStatus(roomID, userID int64, status string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now().Unix(),
	}

	if status == models.ParticipantStatusJoined {
		updates["joined_at"] = time.Now().Unix()
	} else if status == models.ParticipantStatusLeft {
		updates["left_at"] = time.Now().Unix()
		// 计算参与时长
		var participant models.CallParticipant
		if err := d.db.Where("room_id = ? AND user_id = ?", roomID, userID).First(&participant).Error; err == nil && participant.JoinedAt > 0 {
			updates["duration"] = time.Now().Unix() - participant.JoinedAt
		}
	}

	return d.db.Model(&models.CallParticipant{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Updates(updates).Error
}

// CreateCallRecord 创建通话记录
func (d *CallDao) CreateCallRecord(record *models.CallRecord) error {
	record.CreatedAt = time.Now().Unix()
	return d.db.Create(record).Error
}

// GetUserCallRecords 获取用户通话记录
func (d *CallDao) GetUserCallRecords(userID int64, page, pageSize int, callType string) ([]*models.CallRecord, int64, error) {
	var records []*models.CallRecord
	var total int64

	query := d.db.Model(&models.CallRecord{}).Where("user_id = ?", userID)

	if callType != "" && callType != "all" {
		query = query.Where("type = ?", callType)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error

	return records, total, err
}

// GetRoomCallRecords 获取房间通话记录
func (d *CallDao) GetRoomCallRecords(roomID int64) ([]*models.CallRecord, error) {
	var records []*models.CallRecord
	err := d.db.Where("room_id = ?", roomID).Find(&records).Error
	return records, err
}

// EndRoom 结束房间（更新所有参与者状态）
func (d *CallDao) EndRoom(roomID int64) error {
	now := time.Now().Unix()

	// 更新房间状态
	err := d.db.Model(&models.CallRoom{}).Where("id = ?", roomID).Updates(map[string]interface{}{
		"status":     models.RoomStatusEnded,
		"ended_at":   now,
		"updated_at": now,
	}).Error
	if err != nil {
		return err
	}

	// 更新所有参与者的离开时间和时长
	participants, err := d.GetRoomParticipants(roomID)
	if err != nil {
		return err
	}

	for _, p := range participants {
		if p.Status == models.ParticipantStatusJoined && p.JoinedAt > 0 {
			duration := now - p.JoinedAt
			d.db.Model(p).Updates(map[string]interface{}{
				"status":     models.ParticipantStatusLeft,
				"left_at":    now,
				"duration":   duration,
				"updated_at": now,
			})
		}
	}

	return nil
}

// GenerateRoomToken 生成房间令牌
func (d *CallDao) GenerateRoomToken() string {
	// 使用时间戳+随机数生成唯一token
	return fmt.Sprintf("room_%d_%d", time.Now().UnixNano(), time.Now().Unix())
}
