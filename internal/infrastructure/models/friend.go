package models

// Friend 好友关系表
type Friend struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    int64     `gorm:"not null;index" json:"user_id"`   // 用户ID
	FriendID  int64     `gorm:"not null;index" json:"friend_id"` // 好友ID
	Remark    string    `gorm:"size:100" json:"remark"`          // 备注
	GroupName string    `gorm:"size:50;default:''" json:"group_name"` // 自定义分组名称
	IsBlocked bool      `gorm:"default:false" json:"is_blocked"` // 是否屏蔽
	CreatedAt int64 `gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)" json:"created_at"`
	UpdatedAt int64 `gorm:"column:updated_at;type:bigint;not null;default:0;comment:更新时间(时间戳)" json:"updated_at"`
}

// FriendRequest 好友申请表
type FriendRequest struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	FromUserID int64     `gorm:"not null;index" json:"from_user_id"`      // 申请人ID
	ToUserID   int64     `gorm:"not null;index" json:"to_user_id"`        // 被申请人ID
	Message    string    `gorm:"size:500" json:"message"`                 // 申请消息
	Status     string    `gorm:"size:20;default:'pending'" json:"status"` // 状态: pending, accepted, rejected
	CreatedAt  int64 `gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)" json:"created_at"`
	UpdatedAt  int64 `gorm:"column:updated_at;type:bigint;not null;default:0;comment:更新时间(时间戳)" json:"updated_at"`
}

// TableName 指定表名
func (Friend) TableName() string {
	return "friends"
}

// FriendRequestWithUser 好友申请（包含申请人信息）
type FriendRequestWithUser struct {
	FriendRequest
	FromUserName   string `json:"from_user_name"`   // 申请人姓名
	FromUserAvatar string `json:"from_user_avatar"` // 申请人头像
}

// FriendWithUser 好友关系（包含好友信息）
type FriendWithUser struct {
	Friend
	FriendName   string `json:"friend_name"`   // 好友姓名
	FriendPhone  string `json:"friend_phone"`  // 好友手机号
	FriendEmail  string `json:"friend_email"`  // 好友邮箱
	FriendAvatar string `json:"friend_avatar"` // 好友头像
	IsOnline     bool   `json:"is_online"`     // 是否在线
}

// TableName 指定表名
func (FriendRequest) TableName() string {
	return "friend_requests"
}
