package dao

import (
	"gochat/internal/infrastructure/models"

	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

// GetUserList 获取用户列表
func (d *UserDao) GetUserList() ([]*models.UserBasic, error) {
	var users []*models.UserBasic
	err := d.db.Find(&users).Error
	return users, err
}

// CreateUser 创建用户
func (d *UserDao) CreateUser(user *models.UserBasic) error {
	return d.db.Create(user).Error
}

// GetUserByName 根据用户名获取用户
func (d *UserDao) GetUserByName(name string) (*models.UserBasic, error) {
	var user models.UserBasic
	err := d.db.Where("name = ?", name).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 用户不存在，返回 nil, nil
		}
		return nil, err // 其他数据库错误
	}
	return &user, nil
}

// GetUserByLoginAccount 根据登录账号获取用户
func (d *UserDao) GetUserByLoginAccount(loginAccount string) (*models.UserBasic, error) {
	var user models.UserBasic
	err := d.db.Where("login_account = ?", loginAccount).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByPhone 根据手机号获取用户
func (d *UserDao) GetUserByPhone(phone string) (*models.UserBasic, error) {
	var user models.UserBasic
	err := d.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail 根据邮箱获取用户
func (d *UserDao) GetUserByEmail(email string) (*models.UserBasic, error) {
	var user models.UserBasic
	err := d.db.Where("email = ? OR email1 = ? OR email2 = ?", email, email, email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据ID获取用户
func (d *UserDao) GetUserByID(id int64) (*models.UserBasic, error) {
	var user models.UserBasic
	err := d.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 用户不存在，返回 nil, nil
		}
		return nil, err // 其他数据库错误
	}
	return &user, nil
}

// UpdateUser 更新用户（支持部分更新）
func (d *UserDao) UpdateUser(id int64, updates map[string]interface{}) error {
	return d.db.Model(&models.UserBasic{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteUser 删除用户
func (d *UserDao) DeleteUser(id int64) error {
	return d.db.Delete(&models.UserBasic{}, id).Error
}
