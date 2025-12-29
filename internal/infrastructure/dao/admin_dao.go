package dao

import (
	"gochat/internal/infrastructure/models"
	"gorm.io/gorm"
)

type AdminDao struct {
	db *gorm.DB
}

func NewAdminDao(db *gorm.DB) *AdminDao {
	return &AdminDao{db: db}
}

// GetAdminByName 根据用户名获取管理员
func (d *AdminDao) GetAdminByName(name string) (*models.Admin, error) {
	var admin models.Admin
	err := d.db.Where("name = ?", name).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// GetAdminByID 根据ID获取管理员
func (d *AdminDao) GetAdminByID(id int64) (*models.Admin, error) {
	var admin models.Admin
	err := d.db.Where("id = ?", id).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

