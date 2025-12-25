package models

import (
	"gorm.io/gorm"
)

// MessageFavorite 消息收藏模型
type MessageFavorite struct {
	ID        int64          `gorm:"primaryKey;comment:收藏ID"`
	UserID    int64          `gorm:"column:user_id;not null;index;comment:用户ID"`
	MessageID string         `gorm:"column:message_id;type:varchar(100);not null;index:idx_user_message,unique;comment:消息ID"`
	CreatedAt int64          `gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)"`
	UpdatedAt int64          `gorm:"column:updated_at;type:bigint;not null;default:0;comment:更新时间(时间戳)"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index;comment:删除时间"`
}

// TableName 指定表名
func (MessageFavorite) TableName() string {
	return "message_favorites"
}

