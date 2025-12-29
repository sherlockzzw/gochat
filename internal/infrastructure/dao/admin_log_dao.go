package dao

import (
	"gochat/internal/infrastructure/models"
	"gorm.io/gorm"
)

type AdminLogDao struct {
	db *gorm.DB
}

func NewAdminLogDao(db *gorm.DB) *AdminLogDao {
	return &AdminLogDao{db: db}
}

// CreateLog 创建日志
func (d *AdminLogDao) CreateLog(log *models.AdminLog) error {
	return d.db.Create(log).Error
}

// GetLogs 获取日志列表
func (d *AdminLogDao) GetLogs(page, pageSize int, actionType, targetType string, startTime, endTime int64) ([]*models.AdminLog, int64, error) {
	var logs []*models.AdminLog
	var total int64

	query := d.db.Model(&models.AdminLog{})

	if actionType != "" {
		query = query.Where("action_type = ?", actionType)
	}
	if targetType != "" {
		query = query.Where("target_type = ?", targetType)
	}
	if startTime > 0 {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime > 0 {
		query = query.Where("created_at <= ?", endTime)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

