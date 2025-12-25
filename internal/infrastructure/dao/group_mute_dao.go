package dao

import (
	"gochat/internal/infrastructure/models"
	"time"

	"gorm.io/gorm"
)

type GroupMuteDao struct {
	db *gorm.DB
}

func NewGroupMuteDao(db *gorm.DB) *GroupMuteDao {
	return &GroupMuteDao{db: db}
}

// CreateMute 创建禁言记录
func (d *GroupMuteDao) CreateMute(mute *models.GroupMute) error {
	now := time.Now().Unix()
	mute.CreatedAt = now
	mute.UpdatedAt = now
	return d.db.Create(mute).Error
}

// GetMute 获取禁言记录
func (d *GroupMuteDao) GetMute(groupID, userID int64) (*models.GroupMute, error) {
	var mute models.GroupMute
	err := d.db.Where("group_id = ? AND user_id = ? AND deleted_at IS NULL", groupID, userID).
		First(&mute).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &mute, nil
}

// IsMuted 检查用户是否被禁言
func (d *GroupMuteDao) IsMuted(groupID, userID int64) (bool, error) {
	mute, err := d.GetMute(groupID, userID)
	if err != nil {
		return false, err
	}
	if mute == nil {
		return false, nil
	}

	// 检查禁言是否已过期
	now := time.Now().Unix()
	if mute.MutedUntil > 0 && mute.MutedUntil <= now {
		// 禁言已过期，删除记录
		_ = d.RemoveMute(groupID, userID)
		return false, nil
	}

	// 永久禁言或未过期
	return mute.MutedUntil == 0 || mute.MutedUntil > now, nil
}

// RemoveMute 解除禁言
func (d *GroupMuteDao) RemoveMute(groupID, userID int64) error {
	return d.db.Model(&models.GroupMute{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Update("deleted_at", time.Now()).Error
}

// UpdateMute 更新禁言记录
func (d *GroupMuteDao) UpdateMute(groupID, userID int64, mutedUntil int64, reason string) error {
	now := time.Now().Unix()
	return d.db.Model(&models.GroupMute{}).
		Where("group_id = ? AND user_id = ? AND deleted_at IS NULL", groupID, userID).
		Updates(map[string]interface{}{
			"muted_until": mutedUntil,
			"reason":      reason,
			"updated_at":  now,
		}).Error
}

