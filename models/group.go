package models

import (
	"time"

	"gorm.io/gorm"
)

// Group 群组模型
type Group struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`  // 群名称
	Avatar      string         `gorm:"size:255" json:"avatar"`         // 群头像
	OwnerID     uint           `gorm:"not null;index" json:"owner_id"` // 群主ID
	Notice      string         `gorm:"type:text" json:"notice"`        // 群公告
	MemberCount int            `gorm:"default:0" json:"member_count"`  // 成员数量
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (Group) TableName() string {
	return "groups"
}

// GroupMember 群成员模型
type GroupMember struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	GroupID   uint           `gorm:"not null;index" json:"group_id"`       // 群ID
	UserID    uint           `gorm:"not null;index" json:"user_id"`        // 用户ID
	Role      string         `gorm:"size:20;default:'member'" json:"role"` // 角色: owner(群主), admin(管理员), member(普通成员)
	Nickname  string         `gorm:"size:100" json:"nickname"`             // 群内昵称
	JoinedAt  time.Time      `json:"joined_at"`                            // 加入时间
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
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
