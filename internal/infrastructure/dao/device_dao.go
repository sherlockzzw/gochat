package dao

import (
	"time"

	"gochat/internal/infrastructure/models"

	"gorm.io/gorm"
)

type DeviceDao struct {
	db *gorm.DB
}

func NewDeviceDao(db *gorm.DB) *DeviceDao {
	return &DeviceDao{db: db}
}

// CreateDevice 创建设备记录
func (d *DeviceDao) CreateDevice(device *models.UserDevice) error {
	return d.db.Create(device).Error
}

// GetDevice 获取设备信息
func (d *DeviceDao) GetDevice(deviceID string) (*models.UserDevice, error) {
	var device models.UserDevice
	err := d.db.Where("device_id = ?", deviceID).First(&device).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &device, nil
}

// GetDeviceByID 根据ID获取设备
func (d *DeviceDao) GetDeviceByID(id int64) (*models.UserDevice, error) {
	var device models.UserDevice
	err := d.db.Where("id = ?", id).First(&device).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &device, nil
}

// GetUserDevices 获取用户的所有设备
func (d *DeviceDao) GetUserDevices(userID int64) ([]*models.UserDevice, error) {
	var devices []*models.UserDevice
	err := d.db.Where("user_id = ?", userID).
		Order("last_active_at DESC").
		Find(&devices).Error
	return devices, err
}

// UpdateDeviceToken 更新设备Token
func (d *DeviceDao) UpdateDeviceToken(deviceID string, token string) error {
	return d.db.Model(&models.UserDevice{}).
		Where("device_id = ?", deviceID).
		Update("token", token).Error
}

// UpdateDeviceActiveTime 更新设备最后活跃时间
func (d *DeviceDao) UpdateDeviceActiveTime(deviceID string) error {
	return d.db.Model(&models.UserDevice{}).
		Where("device_id = ?", deviceID).
		Update("last_active_at", time.Now().Unix()).Error
}

// DeleteDevice 删除设备(下线)
func (d *DeviceDao) DeleteDevice(id int64) error {
	return d.db.Delete(&models.UserDevice{}, id).Error
}

// DeleteUserDevices 删除用户的所有设备(除指定设备外)
func (d *DeviceDao) DeleteUserDevices(userID int64, excludeDeviceID string) error {
	query := d.db.Where("user_id = ?", userID)
	if excludeDeviceID != "" {
		query = query.Where("device_id != ?", excludeDeviceID)
	}
	return query.Delete(&models.UserDevice{}).Error
}

// GetDeviceByToken 根据Token获取设备
func (d *DeviceDao) GetDeviceByToken(token string) (*models.UserDevice, error) {
	var device models.UserDevice
	err := d.db.Where("token = ?", token).First(&device).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &device, nil
}

// GetTerminalStatistics 获取终端统计
// 参数:
//   - startTime: 开始时间戳（秒），0表示不限制
//   - endTime: 结束时间戳（秒），0表示不限制
// 返回:
//   - terminalData: 终端统计数据，key为设备类型，value为活跃用户数
//   - error: 错误信息
func (d *DeviceDao) GetTerminalStatistics(startTime, endTime int64) (map[string]int64, error) {
	terminalData := make(map[string]int64)

	// 构建基础查询条件
	buildTimeQuery := func(query *gorm.DB) *gorm.DB {
		if startTime > 0 {
			query = query.Where("last_active_at >= ?", startTime)
		}
		if endTime > 0 {
			query = query.Where("last_active_at <= ?", endTime)
		}
		return query
	}

	// 定义设备类型列表
	deviceTypes := []string{
		models.DeviceTypeIOS,
		models.DeviceTypeAndroid,
		models.DeviceTypeWebMobile,
		models.DeviceTypePCDesktop,
		models.DeviceTypePCWeb,
	}

	// 统计每种设备类型的活跃用户数（去重）
	for _, deviceType := range deviceTypes {
		var count int64
		query := d.db.Model(&models.UserDevice{}).
			Where("device_type = ?", deviceType)
		query = buildTimeQuery(query)
		
		// 使用 DISTINCT 去重统计用户数
		if err := query.Select("COUNT(DISTINCT user_id)").Scan(&count).Error; err != nil {
			return nil, err
		}

		// 将设备类型转换为中文显示名称
		var displayName string
		switch deviceType {
		case models.DeviceTypeIOS:
			displayName = "iOS"
		case models.DeviceTypeAndroid:
			displayName = "Android"
		case models.DeviceTypeWebMobile:
			displayName = "手机网页"
		case models.DeviceTypePCDesktop:
			displayName = "PC桌面"
		case models.DeviceTypePCWeb:
			displayName = "PC网页"
		default:
			displayName = deviceType
		}

		terminalData[displayName] = count
	}

	return terminalData, nil
}
