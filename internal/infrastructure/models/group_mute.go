package models

import (
	"gorm.io/gorm"
)

// GroupMute 群组禁言模型
type GroupMute struct {
	ID        int64          `gorm:"primaryKey;comment:禁言ID"`
	GroupID   int64          `gorm:"column:group_id;not null;index;comment:群组ID"`
	UserID    int64          `gorm:"column:user_id;not null;index;comment:被禁言用户ID"`
	MutedBy   int64          `gorm:"column:muted_by;not null;index;comment:禁言操作者ID"`
	MutedUntil int64         `gorm:"column:muted_until;type:bigint;not null;default:0;comment:禁言到期时间(时间戳,0表示永久禁言)"`
	Reason    string         `gorm:"column:reason;type:varchar(500);comment:禁言原因"`
	CreatedAt int64          `gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)"`
	UpdatedAt int64          `gorm:"column:updated_at;type:bigint;not null;default:0;comment:更新时间(时间戳)"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index;comment:删除时间"`
}

// TableName 指定表名
func (GroupMute) TableName() string {
	return "group_mutes"
}


