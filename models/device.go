package models

import (
	"gorm.io/gorm"
)

// UserDevice 用户设备表(多端登录管理)
type UserDevice struct {
	ID          int64          `gorm:"primaryKey;comment:设备ID"`
	UserID      int64          `gorm:"column:user_id;not null;index;comment:用户ID"`
	DeviceType  string         `gorm:"column:device_type;type:varchar(20);not null;index;comment:设备类型(ios,android,web_mobile,pc_desktop,pc_web)"`
	DeviceID    string         `gorm:"column:device_id;type:varchar(100);not null;index;comment:设备唯一标识"`
	DeviceName  string         `gorm:"column:device_name;type:varchar(100);comment:设备名称"`
	Token       string         `gorm:"column:token;type:varchar(500);index;comment:登录Token"`
	LastActiveAt int64          `gorm:"column:last_active_at;type:bigint;not null;default:0;index;comment:最后活跃时间(时间戳)"`
	CreatedAt    int64          `gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)"`
	UpdatedAt    int64          `gorm:"column:updated_at;type:bigint;not null;default:0;comment:更新时间(时间戳)"`
	DeletedAt   gorm.DeletedAt  `gorm:"column:deleted_at;index;comment:删除时间"`
}

// TableName 指定表名
func (UserDevice) TableName() string {
	return "user_devices"
}

// 设备类型常量
const (
	DeviceTypeIOS        = "ios"         // iOS APP
	DeviceTypeAndroid    = "android"     // Android APP
	DeviceTypeWebMobile   = "web_mobile"  // 手机网页端
	DeviceTypePCDesktop  = "pc_desktop"  // PC桌面端
	DeviceTypePCWeb      = "pc_web"      // PC网页端
)

