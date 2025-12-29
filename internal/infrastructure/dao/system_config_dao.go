package dao

import (
	"encoding/json"
	"gochat/internal/infrastructure/models"
	"gorm.io/gorm"
)

type SystemConfigDao struct {
	db *gorm.DB
}

func NewSystemConfigDao(db *gorm.DB) *SystemConfigDao {
	return &SystemConfigDao{db: db}
}

// GetConfig 获取配置
func (d *SystemConfigDao) GetConfig(key string) (*models.SystemConfig, error) {
	var config models.SystemConfig
	err := d.db.Where("config_key = ?", key).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// SetConfig 设置配置
func (d *SystemConfigDao) SetConfig(key, value, description string) error {
	config := &models.SystemConfig{
		ConfigKey:   key,
		ConfigValue: value,
		Description: description,
	}
	return d.db.Where("config_key = ?", key).
		Assign(models.SystemConfig{ConfigValue: value, Description: description}).
		FirstOrCreate(config).Error
}

// GetConfigValue 获取配置值（JSON反序列化）
func (d *SystemConfigDao) GetConfigValue(key string, v interface{}) error {
	config, err := d.GetConfig(key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(config.ConfigValue), v)
}

// SetConfigValue 设置配置值（JSON序列化）
func (d *SystemConfigDao) SetConfigValue(key string, v interface{}, description string) error {
	valueBytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return d.SetConfig(key, string(valueBytes), description)
}

