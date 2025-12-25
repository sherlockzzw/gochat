package models

import (
	"gochat/utils"

	"gorm.io/gorm"
)

type UserBasic struct {
	ID              int64          `gorm:"primaryKey" json:"id"`
	CreatedAt       int64          `gorm:"column:created_at;type:bigint;not null;default:0;comment:创建时间(时间戳)" json:"created_at"`
	UpdatedAt       int64          `gorm:"column:updated_at;type:bigint;not null;default:0;comment:更新时间(时间戳)" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index;comment:删除时间" json:"deleted_at,omitempty"`
	Name            string         `gorm:"column:name;type:varchar(32);not null;default:'';comment:用户名"`
	LoginAccount    string         `gorm:"column:login_account;type:varchar(32);uniqueIndex;not null;default:'';comment:登录账号(5-32位,字母开头,不可修改)"`
	Password        string         `gorm:"column:password;type:varchar(128);not null;default:'';comment:密码"`
	Phone           string         `gorm:"column:phone;type:varchar(32);not null;default:'';comment:手机号"`
	PhoneVerified   bool           `gorm:"column:phone_verified;type:tinyint(1);not null;default:0;comment:手机号是否已验证"`
	Email           string         `gorm:"column:email;type:varchar(64);not null;default:'';comment:邮箱1"`
	Email1          string         `gorm:"column:email1;type:varchar(64);not null;default:'';comment:邮箱2"`
	Email2          string         `gorm:"column:email2;type:varchar(64);not null;default:'';comment:邮箱3"`
	Avatar          string         `gorm:"column:avatar;type:varchar(255);not null;default:'';comment:头像"`
	ClientIp        string         `gorm:"column:client_ip;type:varchar(32);not null;default:'';comment:客户端ip"`
	ClientPort      string         `gorm:"column:client_port;type:varchar(32);not null;default:'';comment:客户端端口"`
	LoginTime       int64          `gorm:"column:login_time;type:bigint;not null;default:0;comment:登录时间"`
	Identity        string         `gorm:"column:identity;type:varchar(36);not null;default:'';comment:唯一标识"`
	HeartTime       int64          `gorm:"column:heart_time;type:bigint;not null;default:0;comment:心跳时间"`
	LogoutTime      int64          `gorm:"column:logout_time;type:bigint;not null;default:0;comment:登出时间"`
	DeviceInfo      string         `gorm:"column:device_info;type:varchar(32);not null;default:'';comment:用户设备"`
	Signature       string         `gorm:"column:signature;type:varchar(200);not null;default:'';comment:个性签名"`
	PrivacySettings string         `gorm:"column:privacy_settings;type:json;comment:隐私设置(JSON格式)"`
	LoginFailCount  int            `gorm:"column:login_fail_count;type:int;not null;default:0;comment:登录失败次数"`
	LockedUntil     int64          `gorm:"column:locked_until;type:bigint;not null;default:0;comment:锁定到期时间(时间戳,0表示未锁定)"`
	LastActiveTime  int64          `gorm:"column:last_active_time;type:bigint;not null;default:0;comment:最后活跃时间(时间戳)"`
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
