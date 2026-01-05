package models

import (
	"gorm.io/gorm"
)

// ConversationSetting 会话设置模型（用于管理会话的置顶、静音、隐藏等）
type ConversationSetting struct {
	ID          int64          `gorm:"primaryKey;comment:会话设置ID"`
	UserID      int64          `gorm:"column:user_id;not null;index;comment:用户ID"`
	OtherUserID int64          `gorm:"column:other_user_id;not null;default:0;index;comment:私聊对方ID，群聊时为0"`
	GroupID     int64          `gorm:"column:group_id;not null;default:0;index;comment:群聊ID，私聊时为0"`
	IsPinned    bool           `gorm:"column:is_pinned;type:tinyint(1);not null;default:0;comment:是否置顶"`
	IsMuted     bool           `gorm:"column:is_muted;type:tinyint(1);not null;default:0;comment:是否静音"`
	IsHidden    bool           `gorm:"column:is_hidden;type:tinyint(1);not null;default:0;comment:是否隐藏"`
	UnreadCount int            `gorm:"column:unread_count;type:int;not null;default:0;comment:未读数量"`
	CreatedAt   int64          `gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)"`
	UpdatedAt   int64          `gorm:"column:updated_at;type:bigint;not null;default:0;comment:更新时间(时间戳)"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index;comment:删除时间"`
}

// TableName 指定表名
func (ConversationSetting) TableName() string {
	return "conversation_settings"
}




