package dao

import (
	"gochat/internal/infrastructure/models"
	"time"

	"gorm.io/gorm"
)

type NotificationDao struct {
	db *gorm.DB
}

func NewNotificationDao(db *gorm.DB) *NotificationDao {
	return &NotificationDao{db: db}
}

// CreateNotification 创建通知
func (d *NotificationDao) CreateNotification(notification *models.NotificationMessage) error {
	now := time.Now().Unix()
	notification.CreatedAt = now
	return d.db.Create(notification).Error
}

// GetNotifications 获取通知列表（分页）
func (d *NotificationDao) GetNotifications(userID int64, page, pageSize int32, notificationType string, isRead *bool) ([]*models.NotificationMessage, int64, error) {
	query := d.db.Model(&models.NotificationMessage{}).Where("user_id = ?", userID)

	// 按类型筛选
	if notificationType != "" {
		query = query.Where("type = ?", notificationType)
	}

	// 按已读状态筛选
	if isRead != nil {
		query = query.Where("is_read = ?", *isRead)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	var notifications []*models.NotificationMessage
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&notifications).Error

	return notifications, total, err
}

// GetUnreadCount 获取未读通知数量
func (d *NotificationDao) GetUnreadCount(userID int64, notificationType string) (int64, error) {
	query := d.db.Model(&models.NotificationMessage{}).
		Where("user_id = ? AND is_read = ?", userID, false)

	if notificationType != "" {
		query = query.Where("type = ?", notificationType)
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}

// MarkAsRead 标记通知为已读
func (d *NotificationDao) MarkAsRead(userID int64, notificationIDs []int64) (int64, error) {
	result := d.db.Model(&models.NotificationMessage{}).
		Where("user_id = ? AND id IN ?", userID, notificationIDs).
		Update("is_read", true)

	return result.RowsAffected, result.Error
}

// MarkAllAsRead 标记所有通知为已读
func (d *NotificationDao) MarkAllAsRead(userID int64, notificationType string) (int64, error) {
	query := d.db.Model(&models.NotificationMessage{}).
		Where("user_id = ? AND is_read = ?", userID, false)

	if notificationType != "" {
		query = query.Where("type = ?", notificationType)
	}

	result := query.Update("is_read", true)
	return result.RowsAffected, result.Error
}

// DeleteNotifications 删除通知（软删除）
func (d *NotificationDao) DeleteNotifications(userID int64, notificationIDs []int64) (int64, error) {
	result := d.db.Where("user_id = ? AND id IN ?", userID, notificationIDs).
		Delete(&models.NotificationMessage{})

	return result.RowsAffected, result.Error
}

// GetNotificationByID 根据ID获取通知
func (d *NotificationDao) GetNotificationByID(notificationID, userID int64) (*models.NotificationMessage, error) {
	var notification models.NotificationMessage
	err := d.db.Where("id = ? AND user_id = ?", notificationID, userID).First(&notification).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &notification, nil
}

