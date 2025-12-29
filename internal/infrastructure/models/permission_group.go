package models

import "gorm.io/gorm"

// PermissionGroup 权限组
type PermissionGroup struct {
	ID          int64          `gorm:"primaryKey;comment:权限组ID"`
	Name        string         `gorm:"column:name;type:varchar(50);not null;uniqueIndex;comment:权限组名称"`
	Description string         `gorm:"column:description;type:text;comment:描述"`
	Permissions string         `gorm:"column:permissions;type:json;comment:权限配置JSON"`
	IsDefault   bool           `gorm:"column:is_default;type:tinyint(1);not null;default:0;comment:是否默认组"`
	CreatedAt   int64          `gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)"`
	UpdatedAt   int64          `gorm:"column:updated_at;type:bigint;not null;default:0;comment:更新时间(时间戳)"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index;comment:删除时间"`
}

func (PermissionGroup) TableName() string {
	return "permission_groups"
}

// UserPermissionGroup 用户权限组关联表
type UserPermissionGroup struct {
	ID        int64 `gorm:"primaryKey;comment:关联ID"`
	UserID    int64 `gorm:"column:user_id;not null;index;comment:用户ID"`
	GroupID   int64 `gorm:"column:group_id;not null;index;comment:权限组ID"`
	CreatedAt int64 `gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)"`
}

func (UserPermissionGroup) TableName() string {
	return "user_permission_groups"
}

// 权限常量定义
const (
	// 消息权限
	PermissionDeleteMessage     = "delete_message"      // 双向删除消息
	PermissionRecallMessage     = "recall_message"      // 撤回消息
	// 群组权限
	PermissionCreateGroup       = "create_group"        // 创建群聊
	PermissionCreateChannel     = "create_channel"      // 创建频道
	// 资金权限
	PermissionRedPacket         = "red_packet"          // 红包使用
	PermissionTransfer           = "transfer"            // 转账使用
	PermissionRecharge          = "recharge"             // 充值
	PermissionWithdraw          = "withdraw"             // 提现
	// 功能权限
	PermissionTranslate         = "translate"           // 易翻译调用
	PermissionCustomEmoji       = "custom_emoji"        // 自定义表情包
)

