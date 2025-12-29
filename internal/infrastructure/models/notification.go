package models

import (
	"gorm.io/gorm"
)

// NotificationMessage 通知消息模型（独立列表，不可删除、不可静音）
type NotificationMessage struct {
	ID        int64          `gorm:"primaryKey;comment:通知消息ID"`
	UserID    int64          `gorm:"column:user_id;not null;index;comment:用户ID"`
	Type      string         `gorm:"column:type;type:varchar(30);not null;index;comment:通知类型(deduction扣款,redpacket红包,transfer转账,system系统)"`
	Title     string         `gorm:"column:title;type:varchar(200);not null;comment:通知标题"`
	Content   string         `gorm:"column:content;type:text;comment:通知内容"`
	Amount    int64          `gorm:"column:amount;type:bigint;not null;default:0;comment:金额(单位:分，如果有)"`
	RelatedID int64          `gorm:"column:related_id;index;comment:关联记录ID(充值/提现/红包/转账ID)"`
	IsRead    bool           `gorm:"column:is_read;type:tinyint(1);not null;default:0;index;comment:是否已读"`
	CreatedAt int64          `gorm:"column:created_at;type:bigint;not null;default:0;index;comment:创建时间(时间戳)"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index;comment:删除时间"`
}

// TableName 指定表名
func (NotificationMessage) TableName() string {
	return "notification_messages"
}

// 通知类型常量
const (
	NotificationTypeDeduction = "deduction" // 扣款通知
	NotificationTypeRedPacket = "redpacket"  // 红包通知
	NotificationTypeTransfer  = "transfer"   // 转账通知
	NotificationTypeSystem   = "system"      // 系统通知
)



