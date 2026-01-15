package dao

import (
	"gochat/internal/infrastructure/models"

	"gorm.io/gorm"
)

// PermissionGroupDao  权限组 DAO
type PermissionGroupDao struct {
	db *gorm.DB
}

func NewPermissionGroupDao(db *gorm.DB) *PermissionGroupDao {
	return &PermissionGroupDao{db: db}
}

func (d *PermissionGroupDao) GetByID(id int64) (*models.PermissionGroup, error) {
	var pg models.PermissionGroup
	if err := d.db.First(&pg, id).Error; err != nil {
		return nil, err
	}
	return &pg, nil
}

func (d *PermissionGroupDao) List(ids []int64) ([]*models.PermissionGroup, error) {
	if len(ids) == 0 {
		return []*models.PermissionGroup{}, nil
	}
	var list []*models.PermissionGroup
	if err := d.db.Where("id IN ?", ids).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
