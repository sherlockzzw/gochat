package models

import (
	"gorm.io/gorm"
)

// Group 群组模型
type Group struct {
	ID          int64          `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`  // 群名称
	Avatar      string         `gorm:"size:255" json:"avatar"`         // 群头像
	OwnerID     int64          `gorm:"not null;index" json:"owner_id"` // 群主ID
	Notice      string         `gorm:"type:text" json:"notice"`        // 群公告
	MemberCount int            `gorm:"default:0" json:"member_count"`  // 成员数量
	CreatedAt   int64          `gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)" json:"created_at"`
	UpdatedAt   int64          `gorm:"column:updated_at;type:bigint;not null;default:0;comment:更新时间(时间戳)" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (Group) TableName() string {
	return "groups"
}

// GroupMember 群成员模型
type GroupMember struct {
	ID        int64          `gorm:"primaryKey" json:"id"`
	GroupID   int64          `gorm:"not null;index" json:"group_id"`       // 群ID
	UserID    int64          `gorm:"not null;index" json:"user_id"`        // 用户ID
	Role      string         `gorm:"size:20;default:'member'" json:"role"` // 角色: owner(群主), admin(管理员), member(普通成员)
	Nickname  string         `gorm:"size:100" json:"nickname"`             // 群内昵称
	JoinedAt  int64          `gorm:"column:joined_at;type:bigint;not null;default:0;comment:加入时间(时间戳)" json:"joined_at"`
	CreatedAt int64          `gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)" json:"created_at"`
	UpdatedAt int64          `gorm:"column:updated_at;type:bigint;not null;default:0;comment:更新时间(时间戳)" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (GroupMember) TableName() string {
	return "group_members"
}

// 群成员角色常量
const (
	GroupRoleOwner  = "owner"  // 群主
	GroupRoleAdmin  = "admin"  // 管理员
	GroupRoleMember = "member" // 普通成员
)

// GroupWithMembers 群组信息（包含成员列表）
type GroupWithMembers struct {
	Group
	Members []GroupMemberInfo `json:"members"`
}

// GroupMemberInfo 群成员信息（包含用户信息）
type GroupMemberInfo struct {
	GroupMember
	UserName   string `json:"user_name"`   // 用户名
	UserAvatar string `json:"user_avatar"` // 用户头像
	UserPhone  string `json:"user_phone"`  // 用户手机号
	IsOnline   bool   `json:"is_online"`   // 是否在线
}
