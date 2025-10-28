package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

// ChatMessage 聊天消息模型
type ChatMessage struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id" gorm:"-"`                                    // MongoDB ID，MySQL中忽略
	MessageID   string             `bson:"message_id" json:"message_id" gorm:"primaryKey;uniqueIndex;not null"` // 业务消息ID作为主键
	FromUserID  uint               `bson:"from_user_id" json:"from_user_id" gorm:"not null;index"`
	ToUserID    uint               `bson:"to_user_id" json:"to_user_id" gorm:"not null;index"`
	MessageType int                `bson:"message_type" json:"message_type" gorm:"not null"` // 0:文字 1:图片 2:文件 3:系统
	Content     string             `bson:"content" json:"content" gorm:"type:text"`
	FileURL     string             `bson:"file_url,omitempty" json:"file_url,omitempty"`
	FileName    string             `bson:"file_name,omitempty" json:"file_name,omitempty"`
	FileSize    int64              `bson:"file_size,omitempty" json:"file_size,omitempty"`
	Status      int                `bson:"status" json:"status" gorm:"not null;default:1"` // 0:发送中 1:已发送 2:已送达 3:已读 4:失败
	CreatedAt   time.Time          `bson:"created_at" json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt     `bson:"deleted_at,omitempty" json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (ChatMessage) TableName() string {
	return "chat_messages"
}

// Conversation 会话模型（用于快速查询会话列表）
type Conversation struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id" gorm:"-"` // MongoDB ID，MySQL中忽略
	UserID          uint               `bson:"user_id" json:"user_id" gorm:"primaryKey;not null;index"`
	OtherUserID     uint               `bson:"other_user_id" json:"other_user_id" gorm:"primaryKey;not null;index"`
	LastMessage     string             `bson:"last_message" json:"last_message" gorm:"type:text"`
	LastMessageType int                `bson:"last_message_type" json:"last_message_type"`
	UnreadCount     int                `bson:"unread_count" json:"unread_count" gorm:"default:0"`
	LastMessageAt   time.Time          `bson:"last_message_at" json:"last_message_at"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt     `bson:"deleted_at,omitempty" json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (Conversation) TableName() string {
	return "conversations"
}

// MessageReadStatus 消息已读状态
type MessageReadStatus struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id" gorm:"-"` // MongoDB ID，MySQL中忽略
	MessageID string             `bson:"message_id" json:"message_id" gorm:"primaryKey;uniqueIndex;not null"`
	UserID    uint               `bson:"user_id" json:"user_id" gorm:"not null;index"`
	ReadAt    time.Time          `bson:"read_at" json:"read_at" gorm:"autoCreateTime"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at" gorm:"autoCreateTime"`
}

// TableName 指定表名
func (MessageReadStatus) TableName() string {
	return "message_read_status"
}

// 消息类型常量
const (
	MessageTypeText   = 0 // 文字消息
	MessageTypeImage  = 1 // 图片消息
	MessageTypeFile   = 2 // 文件消息
	MessageTypeSystem = 3 // 系统消息
)

// 消息状态常量
const (
	MessageStatusSending   = 0 // 发送中
	MessageStatusSent      = 1 // 已发送
	MessageStatusDelivered = 2 // 已送达
	MessageStatusRead      = 3 // 已读
	MessageStatusFailed    = 4 // 发送失败
)
