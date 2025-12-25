package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

// ChatMessage 聊天消息模型
type ChatMessage struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id" gorm:"-"`                                    // MongoDB ID，MySQL中忽略
	MessageID   string             `bson:"message_id" json:"message_id" gorm:"primaryKey;uniqueIndex;not null"` // 业务消息ID作为主键
	FromUserID  int64              `bson:"from_user_id" json:"from_user_id" gorm:"not null;index"`
	ToUserID    int64              `bson:"to_user_id" json:"to_user_id" gorm:"index"`                          // 私聊接收者ID，群聊时为0
	GroupID     int64              `bson:"group_id" json:"group_id" gorm:"index"`                               // 群聊群组ID，私聊时为0
	MessageType int                `bson:"message_type" json:"message_type" gorm:"not null"`                    // 消息类型
	Content     string             `bson:"content" json:"content" gorm:"type:text"`                            // 消息内容
	FileURL     string             `bson:"file_url,omitempty" json:"file_url,omitempty"`                        // 文件URL（图片/文件消息）
	FileName    string             `bson:"file_name,omitempty" json:"file_name,omitempty"`                       // 文件名
	FileSize    int64              `bson:"file_size,omitempty" json:"file_size,omitempty"`                      // 文件大小
	// 扩展字段
	VideoURL      string `bson:"video_url,omitempty" json:"video_url,omitempty"`                    // 视频URL
	VideoThumb    string `bson:"video_thumb,omitempty" json:"video_thumb,omitempty"`                // 视频缩略图
	VoiceURL      string `bson:"voice_url,omitempty" json:"voice_url,omitempty"`                    // 语音URL
	VoiceDuration int    `bson:"voice_duration,omitempty" json:"voice_duration,omitempty"`          // 语音时长（秒）
	EmojiURL      string `bson:"emoji_url,omitempty" json:"emoji_url,omitempty"`                      // 表情包URL
	MergeMessages string `bson:"merge_messages,omitempty" json:"merge_messages,omitempty"`          // JSON数组，合并消息的message_id列表
	QuoteMessageID string `bson:"quote_message_id,omitempty" json:"quote_message_id,omitempty"`     // 引用的消息ID
	ContactUserID int64  `bson:"contact_user_id,omitempty" json:"contact_user_id,omitempty"`        // 分享的联系人ID
	RedPacketID   int64  `bson:"red_packet_id,omitempty" json:"red_packet_id,omitempty"`              // 关联的红包ID
	TransferID    int64  `bson:"transfer_id,omitempty" json:"transfer_id,omitempty"`                // 关联的转账ID
	// 状态字段
	IsRecalled bool `bson:"is_recalled,omitempty" json:"is_recalled,omitempty"` // 是否已撤回
	IsDeleted  bool `bson:"is_deleted,omitempty" json:"is_deleted,omitempty"`    // 是否已删除
	ReadStatus int  `bson:"read_status,omitempty" json:"read_status,omitempty"` // 已读状态：1=已发送，2=已读
	Status     int  `bson:"status" json:"status" gorm:"not null;default:1"`     // 0:发送中 1:已发送 2:已送达 3:已读 4:失败
	CreatedAt  int64              `bson:"created_at" json:"created_at" gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)"`
	UpdatedAt  int64              `bson:"updated_at" json:"updated_at" gorm:"column:updated_at;type:bigint;not null;default:0;comment:更新时间(时间戳)"`
	DeletedAt  gorm.DeletedAt     `bson:"deleted_at,omitempty" json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (ChatMessage) TableName() string {
	return "chat_messages"
}

// Conversation 会话模型（用于快速查询会话列表）
type Conversation struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id" gorm:"-"` // MongoDB ID，MySQL中忽略
	UserID          int64             `bson:"user_id" json:"user_id" gorm:"primaryKey;not null;index"`
	OtherUserID     int64             `bson:"other_user_id" json:"other_user_id" gorm:"primaryKey;not null;index"`
	LastMessage     string             `bson:"last_message" json:"last_message" gorm:"type:text"`
	LastMessageType int                `bson:"last_message_type" json:"last_message_type"`
	UnreadCount     int                `bson:"unread_count" json:"unread_count" gorm:"default:0"`
	LastMessageAt   int64              `bson:"last_message_at" json:"last_message_at" gorm:"column:last_message_at;type:bigint;not null;default:0;comment:最后消息时间(时间戳)"`
	CreatedAt       int64              `bson:"created_at" json:"created_at" gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)"`
	UpdatedAt       int64              `bson:"updated_at" json:"updated_at" gorm:"column:updated_at;type:bigint;not null;default:0;comment:更新时间(时间戳)"`
	DeletedAt       gorm.DeletedAt     `bson:"deleted_at,omitempty" json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (Conversation) TableName() string {
	return "conversations"
}

// MessageReadStatus 消息已读状态
type MessageReadStatus struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id" gorm:"-"` // MongoDB ID，MySQL中忽略
	MessageID string             `bson:"message_id" json:"message_id" gorm:"primaryKey;not null;index"`
	UserID    int64              `bson:"user_id" json:"user_id" gorm:"primaryKey;not null;index"`
	ReadAt    int64              `bson:"read_at" json:"read_at" gorm:"column:read_at;type:bigint;not null;default:0;comment:已读时间(时间戳)"`
	CreatedAt int64              `bson:"created_at" json:"created_at" gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)"`
}

// TableName 指定表名
func (MessageReadStatus) TableName() string {
	return "message_read_status"
}

// 消息类型常量
const (
	MessageTypeText     = 0  // 文字消息
	MessageTypeImage    = 1  // 图片消息
	MessageTypeFile     = 2  // 文件消息
	MessageTypeSystem   = 3  // 系统消息
	MessageTypeVideo    = 4  // 视频消息
	MessageTypeVoice    = 5  // 语音条消息
	MessageTypeEmoji    = 6  // 表情包消息
	MessageTypeMerge    = 7  // 合并消息
	MessageTypeQuote    = 8  // 引用消息
	MessageTypeContact  = 9  // 联系人分享消息
	MessageTypeRedPacket = 10 // 红包消息
	MessageTypeTransfer = 11 // 转账消息
)

// 消息状态常量
const (
	MessageStatusSending   = 0 // 发送中
	MessageStatusSent      = 1 // 已发送
	MessageStatusDelivered = 2 // 已送达
	MessageStatusRead      = 3 // 已读
	MessageStatusFailed    = 4 // 发送失败
)
