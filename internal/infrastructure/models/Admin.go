package models

import "gorm.io/gorm"

type Admin struct {
	ID        int64          `gorm:"primaryKey" json:"id"`
	CreatedAt int64          `gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)" json:"created_at"`
	UpdatedAt int64          `gorm:"column:updated_at;type:bigint;not null;default:0;comment:更新时间(时间戳)" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index;comment:删除时间" json:"deleted_at,omitempty"`
	Name      string         `gorm:"column:name;type:varchar(32);not null;default:'';comment:用户名"`
	Password  string         `gorm:"column:password;type:varchar(128);not null;default:'';comment:密码"`
}

func (table *Admin) TableName() string {
	return "admin"
}
