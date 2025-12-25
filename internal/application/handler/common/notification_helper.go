package common

import (
	"gochat/internal/infrastructure/dao"
	"gochat/internal/infrastructure/models"
	"gochat/utils"
	"time"
)

// CreateNotification 创建通知的辅助函数
func CreateNotification(userID int64, notificationType, title, content string, amount int64, relatedID int64) error {
	notificationDao := dao.NewNotificationDao(utils.DB)
	
	notification := &models.NotificationMessage{
		UserID:    userID,
		Type:      notificationType,
		Title:     title,
		Content:   content,
		Amount:    amount,
		RelatedID: relatedID,
		IsRead:    false,
		CreatedAt: time.Now().Unix(),
	}

	return notificationDao.CreateNotification(notification)
}

// CreateRedPacketNotification 创建红包通知
func CreateRedPacketNotification(userID int64, title, content string, amount, redPacketID int64) error {
	return CreateNotification(
		userID,
		models.NotificationTypeRedPacket,
		title,
		content,
		amount,
		redPacketID,
	)
}

// CreateTransferNotification 创建转账通知
func CreateTransferNotification(userID int64, title, content string, amount, transferID int64) error {
	return CreateNotification(
		userID,
		models.NotificationTypeTransfer,
		title,
		content,
		amount,
		transferID,
	)
}

// CreateSystemNotification 创建系统通知
func CreateSystemNotification(userID int64, title, content string, relatedID int64) error {
	return CreateNotification(
		userID,
		models.NotificationTypeSystem,
		title,
		content,
		0,
		relatedID,
	)
}

// CreateDeductionNotification 创建扣款通知（充值/提现等）
func CreateDeductionNotification(userID int64, title, content string, amount, relatedID int64) error {
	return CreateNotification(
		userID,
		models.NotificationTypeDeduction,
		title,
		content,
		amount,
		relatedID,
	)
}

