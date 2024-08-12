package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Name     string `gorm:"column:name;type:varchar(32);not null;default:'';comment:用户名"`
	Password string `gorm:"column:password;type:varchar(128);not null;default:'';comment:密码"`
}

func (table *Admin) TableName() string {
	return "admin"
}
