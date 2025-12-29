package dao

import (
	"gochat/internal/infrastructure/models"
	"gorm.io/gorm"
)

type ContentViolationDao struct {
	db *gorm.DB
}

func NewContentViolationDao(db *gorm.DB) *ContentViolationDao {
	return &ContentViolationDao{db: db}
}

// CreateViolation 创建违规记录
func (d *ContentViolationDao) CreateViolation(violation *models.ContentViolation) error {
	return d.db.Create(violation).Error
}

// GetViolations 获取违规内容列表
func (d *ContentViolationDao) GetViolations(page, pageSize int, violationType, status string) ([]*models.ContentViolation, int64, error) {
	var violations []*models.ContentViolation
	var total int64

	query := d.db.Model(&models.ContentViolation{})

	if violationType != "" {
		query = query.Where("type = ?", violationType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&violations).Error; err != nil {
		return nil, 0, err
	}

	return violations, total, nil
}

// DeleteViolation 删除违规内容
func (d *ContentViolationDao) DeleteViolation(id int64) error {
	return d.db.Delete(&models.ContentViolation{}, id).Error
}

// UpdateViolationStatus 更新违规状态
func (d *ContentViolationDao) UpdateViolationStatus(id int64, status string, processedBy int64) error {
	return d.db.Model(&models.ContentViolation{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       status,
			"processed_by": processedBy,
			"processed_at": 0, // TODO: 使用当前时间戳
		}).Error
}

