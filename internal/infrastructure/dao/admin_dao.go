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

// ListWithPagination 分页获取管理员列表
func (d *AdminDao) ListWithPagination(page, pageSize int) (list []*models.Admin, total int64, err error) {
    if page <= 0 {
        page = 1
    }
    if pageSize <= 0 {
        pageSize = 20
    }
    query := d.db.Model(&models.Admin{})
    if err = query.Count(&total).Error; err != nil {
        return
    }
    offset := (page - 1) * pageSize
    err = query.Offset(offset).Limit(pageSize).Order("id desc").Find(&list).Error
    return
}

// CreateAdmin 创建管理员
func (d *AdminDao) CreateAdmin(admin *models.Admin) error {
    return d.db.Create(admin).Error
}

// UpdateAdmin 局部更新
func (d *AdminDao) UpdateAdmin(id int64, updates map[string]interface{}) error {
    return d.db.Model(&models.Admin{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteAdmin 软删除
func (d *AdminDao) DeleteAdmin(id int64) error {
    return d.db.Delete(&models.Admin{}, id).Error
}