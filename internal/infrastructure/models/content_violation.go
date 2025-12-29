package models

import "gorm.io/gorm"

// ContentViolation 违规内容
type ContentViolation struct {
	ID          int64          `gorm:"primaryKey;comment:违规ID"`
	Type        string         `gorm:"column:type;type:varchar(50);not null;index;comment:类型(emoji违规表情包,finance异常资金操作)"`
	UserID      int64          `gorm:"column:user_id;not null;index;comment:用户ID"`
	Content     string         `gorm:"column:content;type:text;comment:违规内容"`
	ContentID   int64          `gorm:"column:content_id;index;comment:内容ID(表情包ID/交易ID等)"`
	Reason      string         `gorm:"column:reason;type:varchar(200);comment:违规原因"`
	Status      string         `gorm:"column:status;type:varchar(20);not null;default:'pending';index;comment:状态(pending待处理,processed已处理,ignored已忽略)"`
	ProcessedBy int64          `gorm:"column:processed_by;index;comment:处理人ID"`
	ProcessedAt int64          `gorm:"column:processed_at;type:bigint;default:0;comment:处理时间(时间戳)"`
	CreatedAt   int64          `gorm:"column:created_at;type:bigint;not null;default:0;index;comment:创建时间(时间戳)"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index;comment:删除时间"`
}

func (ContentViolation) TableName() string {
	return "content_violations"
}

// 违规类型常量
const (
	ViolationTypeEmoji  = "emoji"   // 违规表情包
	ViolationTypeFinance = "finance" // 异常资金操作
)

// 违规状态常量
const (
	ViolationStatusPending   = "pending"   // 待处理
	ViolationStatusProcessed = "processed" // 已处理
	ViolationStatusIgnored   = "ignored"   // 已忽略
)

