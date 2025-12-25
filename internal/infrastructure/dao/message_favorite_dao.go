package dao

import (
	"gochat/internal/infrastructure/models"
	"time"

	"gorm.io/gorm"
)

type MessageFavoriteDao struct {
	db *gorm.DB
}

func NewMessageFavoriteDao(db *gorm.DB) *MessageFavoriteDao {
	return &MessageFavoriteDao{db: db}
}

// CreateFavorite 创建收藏
func (d *MessageFavoriteDao) CreateFavorite(userID int64, messageID string) error {
	now := time.Now().Unix()
	favorite := &models.MessageFavorite{
		UserID:    userID,
		MessageID: messageID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return d.db.Create(favorite).Error
}

// DeleteFavorite 取消收藏
func (d *MessageFavoriteDao) DeleteFavorite(userID int64, messageID string) error {
	return d.db.Where("user_id = ? AND message_id = ?", userID, messageID).
		Delete(&models.MessageFavorite{}).Error
}

// GetFavorite 获取收藏记录
func (d *MessageFavoriteDao) GetFavorite(userID int64, messageID string) (*models.MessageFavorite, error) {
	var favorite models.MessageFavorite
	err := d.db.Where("user_id = ? AND message_id = ?", userID, messageID).
		First(&favorite).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &favorite, nil
}

// GetFavoriteMessageIDs 获取用户收藏的消息ID列表
func (d *MessageFavoriteDao) GetFavoriteMessageIDs(userID int64, page, pageSize int32) ([]string, int64, error) {
	var favorites []models.MessageFavorite
	var total int64

	// 获取总数
	err := d.db.Model(&models.MessageFavorite{}).
		Where("user_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = d.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&favorites).Error
	if err != nil {
		return nil, 0, err
	}

	// 提取消息ID列表
	messageIDs := make([]string, 0, len(favorites))
	for _, fav := range favorites {
		messageIDs = append(messageIDs, fav.MessageID)
	}

	return messageIDs, total, nil
}

