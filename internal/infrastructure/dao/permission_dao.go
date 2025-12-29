package dao

import (
	"gochat/internal/infrastructure/models"
	"gorm.io/gorm"
)

type PermissionDao struct {
	db *gorm.DB
}

func NewPermissionDao(db *gorm.DB) *PermissionDao {
	return &PermissionDao{db: db}
}

// CreatePermissionGroup 创建权限组
func (d *PermissionDao) CreatePermissionGroup(group *models.PermissionGroup) error {
	return d.db.Create(group).Error
}

// GetPermissionGroupByID 根据ID获取权限组
func (d *PermissionDao) GetPermissionGroupByID(id int64) (*models.PermissionGroup, error) {
	var group models.PermissionGroup
	err := d.db.Where("id = ?", id).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// GetPermissionGroups 获取权限组列表
func (d *PermissionDao) GetPermissionGroups(page, pageSize int) ([]*models.PermissionGroup, int64, error) {
	var groups []*models.PermissionGroup
	var total int64

	query := d.db.Model(&models.PermissionGroup{})

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&groups).Error; err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}

// UpdatePermissionGroup 更新权限组
func (d *PermissionDao) UpdatePermissionGroup(group *models.PermissionGroup) error {
	return d.db.Save(group).Error
}

// DeletePermissionGroup 删除权限组
func (d *PermissionDao) DeletePermissionGroup(id int64) error {
	return d.db.Delete(&models.PermissionGroup{}, id).Error
}

// AssignUserToGroup 分配用户到权限组
func (d *PermissionDao) AssignUserToGroup(userID, groupID int64) error {
	upg := &models.UserPermissionGroup{
		UserID:  userID,
		GroupID: groupID,
	}
	return d.db.Create(upg).Error
}

// GetUserGroups 获取用户的权限组
func (d *PermissionDao) GetUserGroups(userID int64) ([]*models.PermissionGroup, error) {
	var groups []*models.PermissionGroup
	err := d.db.Table("permission_groups").
		Joins("JOIN user_permission_groups ON permission_groups.id = user_permission_groups.group_id").
		Where("user_permission_groups.user_id = ?", userID).
		Find(&groups).Error
	return groups, err
}

// GetGroupUserCount 获取权限组的用户数量
func (d *PermissionDao) GetGroupUserCount(groupID int64) (int64, error) {
	var count int64
	err := d.db.Model(&models.UserPermissionGroup{}).
		Where("group_id = ?", groupID).
		Count(&count).Error
	return count, err
}

