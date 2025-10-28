package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MessageType 消息类型
type MessageType string

const (
	MessageTypeText  MessageType = "text"  // 文本消息
	MessageTypeImage MessageType = "image" // 图片消息
	MessageTypeFile  MessageType = "file"  // 文件消息
	MessageTypeAudio MessageType = "audio" // 语音消息
	MessageTypeVideo MessageType = "video" // 视频消息
)

// MessageStatus 消息状态
type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "sent"      // 已发送
	MessageStatusDelivered MessageStatus = "delivered" // 已送达
	MessageStatusRead      MessageStatus = "read"      // 已读
	MessageStatusFailed    MessageStatus = "failed"    // 发送失败
)

// Message 消息模型
type Message struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FromUserID  uint               `bson:"from_user_id" json:"from_user_id"`
	ToUserID    uint               `bson:"to_user_id,omitempty" json:"to_user_id,omitempty"` // 私聊消息
	RoomID      string             `bson:"room_id,omitempty" json:"room_id,omitempty"`       // 群聊消息
	Type        MessageType        `bson:"type" json:"type"`
	Content     string             `bson:"content" json:"content"`
	FileURL     string             `bson:"file_url,omitempty" json:"file_url,omitempty"`
	FileName    string             `bson:"file_name,omitempty" json:"file_name,omitempty"`
	FileSize    int64              `bson:"file_size,omitempty" json:"file_size,omitempty"`
	Status      MessageStatus      `bson:"status" json:"status"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	ReadAt      *time.Time         `bson:"read_at,omitempty" json:"read_at,omitempty"`
	DeliveredAt *time.Time         `bson:"delivered_at,omitempty" json:"delivered_at,omitempty"`
}

// ChatRoom 聊天室模型
type ChatRoom struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Avatar      string             `bson:"avatar,omitempty" json:"avatar,omitempty"`
	Type        string             `bson:"type" json:"type"` // "private", "group"
	OwnerID     uint               `bson:"owner_id" json:"owner_id"`
	Members     []uint             `bson:"members" json:"members"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// Friend 好友关系模型
type Friend struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    uint               `bson:"user_id" json:"user_id"`
	FriendID  uint               `bson:"friend_id" json:"friend_id"`
	Status    string             `bson:"status" json:"status"` // "pending", "accepted", "blocked"
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// MessageRead 消息已读记录
type MessageRead struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MessageID primitive.ObjectID `bson:"message_id" json:"message_id"`
	UserID    uint               `bson:"user_id" json:"user_id"`
	ReadAt    time.Time          `bson:"read_at" json:"read_at"`
}

// UnreadCount 未读消息计数
type UnreadCount struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     uint               `bson:"user_id" json:"user_id"`
	FromUserID uint               `bson:"from_user_id,omitempty" json:"from_user_id,omitempty"` // 私聊未读数
	RoomID     string             `bson:"room_id,omitempty" json:"room_id,omitempty"`           // 群聊未读数
	Count      int                `bson:"count" json:"count"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}
