package dao

import (
    "gochat/internal/infrastructure/models"

    "gorm.io/gorm"
)

type UserPermissionGroupDao struct {
    db *gorm.DB
}

func NewUserPermissionGroupDao(db *gorm.DB) *UserPermissionGroupDao {
    return &UserPermissionGroupDao{db: db}
}

// GetGroupIDByAdminID 获取管理员绑定的权限组ID（一个管理员只属于一个组的简单场景）
func (d *UserPermissionGroupDao) GetGroupIDByAdminID(adminID int64) (int64, error) {
    var upg models.UserPermissionGroup
    if err := d.db.Where("user_id = ?", adminID).First(&upg).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return 0, nil
        }
        return 0, err
    }
    return upg.GroupID, nil
}

// GetGroupIDsByAdminIDs 批量获取
func (d *UserPermissionGroupDao) GetGroupIDsByAdminIDs(adminIDs []int64) (map[int64]int64, error) {
    if len(adminIDs) == 0 {
        return map[int64]int64{}, nil
    }
    var list []models.UserPermissionGroup
    if err := d.db.Where("user_id IN ?", adminIDs).Find(&list).Error; err != nil {
        return nil, err
    }
    res := make(map[int64]int64, len(list))
    for _, v := range list {
        res[v.UserID] = v.GroupID
    }
    return res, nil
}

// Bind 绑定管理员至权限组（若已存在则更新）
func (d *UserPermissionGroupDao) Bind(adminID, groupID int64) error {
    var upg models.UserPermissionGroup
    if err := d.db.Where("user_id = ?", adminID).First(&upg).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return d.db.Create(&models.UserPermissionGroup{UserID: adminID, GroupID: groupID}).Error
        }
        return err
    }
    return d.db.Model(&models.UserPermissionGroup{}).Where("user_id = ?", adminID).Update("group_id", groupID).Error
}

// DeleteByAdminID 删除绑定
func (d *UserPermissionGroupDao) DeleteByAdminID(adminID int64) error {
    return d.db.Where("user_id = ?", adminID).Delete(&models.UserPermissionGroup{}).Error
}
