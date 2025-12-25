package models

import "gorm.io/gorm"

// CallRoom 通话房间
type CallRoom struct {
	ID          int64          `gorm:"primaryKey;comment:房间ID"`
	Type        string         `gorm:"column:type;type:varchar(20);not null;index;comment:类型(voice语音,video视频)"`
	CallType    string         `gorm:"column:call_type;type:varchar(20);not null;index;comment:通话类型(private私聊,group群聊)"`
	CreatorID   int64          `gorm:"column:creator_id;not null;index;comment:创建者ID"`
	GroupID     int64          `gorm:"column:group_id;index;comment:群组ID(群聊时使用,私聊时为0)"`
	RoomToken   string         `gorm:"column:room_token;type:varchar(100);uniqueIndex;not null;comment:房间令牌"`
	Status      string         `gorm:"column:status;type:varchar(20);not null;default:'calling';index;comment:状态(calling呼叫中,ringing响铃中,connected已连接,ended已结束)"`
	StartedAt   int64          `gorm:"column:started_at;type:bigint;not null;default:0;comment:开始时间(时间戳)"`
	EndedAt     int64          `gorm:"column:ended_at;type:bigint;not null;default:0;comment:结束时间(时间戳)"`
	Duration    int64          `gorm:"column:duration;type:bigint;not null;default:0;comment:通话时长(秒)"`
	CreatedAt   int64          `gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)"`
	UpdatedAt   int64          `gorm:"column:updated_at;type:bigint;not null;default:0;comment:更新时间(时间戳)"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index;comment:删除时间"`
}

func (CallRoom) TableName() string {
	return "call_rooms"
}

// CallParticipant 通话参与者
type CallParticipant struct {
	ID        int64  `gorm:"primaryKey;comment:参与者ID"`
	RoomID    int64  `gorm:"column:room_id;not null;index;comment:房间ID"`
	UserID    int64  `gorm:"column:user_id;not null;index;comment:用户ID"`
	Status    string `gorm:"column:status;type:varchar(20);not null;default:'invited';index;comment:状态(invited已邀请,ringing响铃中,joined已加入,rejected已拒绝,left已离开)"`
	JoinedAt  int64  `gorm:"column:joined_at;type:bigint;not null;default:0;comment:加入时间(时间戳)"`
	LeftAt    int64  `gorm:"column:left_at;type:bigint;not null;default:0;comment:离开时间(时间戳)"`
	Duration  int64  `gorm:"column:duration;type:bigint;not null;default:0;comment:参与时长(秒)"`
	CreatedAt int64  `gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)"`
	UpdatedAt int64  `gorm:"column:updated_at;type:bigint;not null;default:0;comment:更新时间(时间戳)"`
}

func (CallParticipant) TableName() string {
	return "call_participants"
}

// CallRecord 通话记录
type CallRecord struct {
	ID          int64  `gorm:"primaryKey;comment:记录ID"`
	RoomID      int64  `gorm:"column:room_id;not null;index;comment:房间ID"`
	UserID      int64  `gorm:"column:user_id;not null;index;comment:用户ID"`
	Type        string `gorm:"column:type;type:varchar(20);not null;index;comment:类型(voice语音,video视频)"`
	CallType    string `gorm:"column:call_type;type:varchar(20);not null;index;comment:通话类型(private私聊,group群聊)"`
	OtherUserID int64  `gorm:"column:other_user_id;index;comment:对方用户ID(私聊时使用)"`
	GroupID     int64  `gorm:"column:group_id;index;comment:群组ID(群聊时使用)"`
	Direction   string `gorm:"column:direction;type:varchar(20);not null;index;comment:方向(incoming来电,outgoing去电)"`
	Status      string `gorm:"column:status;type:varchar(20);not null;index;comment:状态(missed未接,answered已接,rejected已拒绝)"`
	Duration    int64  `gorm:"column:duration;type:bigint;not null;default:0;comment:通话时长(秒)"`
	StartedAt   int64  `gorm:"column:started_at;type:bigint;not null;default:0;comment:开始时间(时间戳)"`
	EndedAt     int64  `gorm:"column:ended_at;type:bigint;not null;default:0;comment:结束时间(时间戳)"`
	CreatedAt   int64  `gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)"`
}

func (CallRecord) TableName() string {
	return "call_records"
}

// 常量定义
const (
	// 通话类型
	CallTypeVoice = "voice"
	CallTypeVideo = "video"

	// 通话场景
	CallScenePrivate = "private"
	CallSceneGroup   = "group"

	// 房间状态
	RoomStatusCalling   = "calling"   // 呼叫中
	RoomStatusRinging   = "ringing"   // 响铃中
	RoomStatusConnected = "connected" // 已连接
	RoomStatusEnded     = "ended"     // 已结束
	RoomStatusCancelled = "cancelled" // 已取消

	// 参与者状态
	ParticipantStatusInvited  = "invited"  // 已邀请
	ParticipantStatusRinging  = "ringing"  // 响铃中
	ParticipantStatusJoined   = "joined"   // 已加入
	ParticipantStatusRejected = "rejected" // 已拒绝
	ParticipantStatusLeft     = "left"     // 已离开

	// 通话方向
	CallDirectionIncoming = "incoming" // 来电
	CallDirectionOutgoing = "outgoing" // 去电

	// 通话记录状态
	RecordStatusMissed   = "missed"   // 未接
	RecordStatusAnswered = "answered" // 已接
	RecordStatusRejected = "rejected" // 已拒绝
	RecordStatusCancelled = "cancelled" // 已取消

	// 参与者状态
	ParticipantStatusMissed = "missed" // 已错过
)

