package models

import (
	"gochat/utils"
	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name       string `gorm:"column:name;type:varchar(32);not null;default:'';comment:用户名"`
	Password   string `gorm:"column:password;type:varchar(128);not null;default:'';comment:密码"`
	Phone      string `gorm:"column:phone;type:varchar(32);not null;default:'';comment:手机号"`
	Email      string `gorm:"column:email;type:varchar(32);not null;default:'';comment:邮箱"`
	ClientIp   string `gorm:"column:client_ip;type:varchar(32);not null;default:'';comment:客户端ip"`
	ClientPort string `gorm:"column:client_port;type:varchar(32);not null;default:'';comment:客户端端口"`
	LoginTime  int64  `gorm:"column:login_time;type:bigint not null;default:0;comment:登录时间"`
	Identity   string `gorm:"column:identity;type:varchar(36);not null;default:'';comment:唯一标识"`
	HeartTime  int64  `gorm:"column:heart_time;type:bigint not null;default:0;comment:心跳时间"`
	LogoutTime int64  `gorm:"column:logout_time;type:bigint not null;default:0;comment:登出时间"`
	DeviceInfo string `gorm:"column:device_info;type:varchar(32);not null;default:'';comment:用户设备"`
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}
func GetUserList() []*UserBasic {
	var data []*UserBasic
	utils.DB.Find(&data)
	return data
}

func CreateUser(user *UserBasic) error {
	result := utils.DB.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func GetUserByName(name string) (*UserBasic, error) {
	var user UserBasic
	result := utils.DB.Where("name = ?", name).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}
