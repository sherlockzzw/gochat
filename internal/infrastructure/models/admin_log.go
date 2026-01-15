package models

// AdminLog 管理员操作日志
type AdminLog struct {
	ID          int64  `gorm:"primaryKey;comment:日志ID"`
	AdminID     int64  `gorm:"column:admin_id;not null;index;comment:管理员ID"`
	AdminName   string `gorm:"column:admin_name;type:varchar(50);comment:管理员名称"`
	ActionType  string `gorm:"column:action_type;type:varchar(50);not null;index;comment:操作类型"`
	Description string `gorm:"column:description;type:text;comment:操作描述"`
	TargetType  string `gorm:"column:target_type;type:varchar(50);index;comment:目标类型(user/permission/config/finance/content)"`
	TargetID    int64  `gorm:"column:target_id;index;comment:目标ID"`
	IPAddress   string `gorm:"column:ip_address;type:varchar(50);comment:IP地址"`
	UserAgent   string `gorm:"column:user_agent;type:varchar(500);comment:用户代理"`
	CreatedAt   int64  `gorm:"column:created_at;type:bigint;not null;default:0;index;comment:创建时间(时间戳)"`
}

func (AdminLog) TableName() string {
	return "admin_logs"
}

// 操作类型常量
const (
	ActionTypePermissionCreate   = "permission_create"   // 创建权限组
	ActionTypePermissionUpdate   = "permission_update"  // 更新权限组
	ActionTypePermissionDelete   = "permission_delete"  // 删除权限组
	ActionTypePermissionAssign   = "permission_assign"  // 分配用户权限
	ActionTypeConfigUpdate       = "config_update"      // 更新配置
	ActionTypeFinanceAudit       = "finance_audit"      // 审核资金
	ActionTypeFinanceAdjust      = "finance_adjust"      // 调整余额
	ActionTypeContentDelete      = "content_delete"     // 删除违规内容
	ActionTypeContentBlock       = "content_block"       // 拦截异常操作
	ActionTypeUserCreate         = "user_create"        // 创建用户
	ActionTypeUserUpdate         = "user_update"        // 更新用户
	ActionTypeUserDisable        = "user_disable"      // 禁用用户
	ActionTypeFeatureToggle      = "feature_toggle"     // 功能开关
	ActionTypeAdminLogin         = "admin_login"        // 管理员登录
	ActionTypeAdminLogout        = "admin_logout"       // 管理员登出
	ActionTypeAdminCreate        = "admin_create"       // 创建管理员
	ActionTypeAdminUpdate        = "admin_update"       // 更新管理员
	ActionTypeAdminDelete        = "admin_delete"       // 删除管理员
)

