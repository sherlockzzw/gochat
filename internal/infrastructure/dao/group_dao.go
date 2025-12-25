package dao

import (
	"fmt"
	"time"

	"gochat/internal/infrastructure/models"

	"gorm.io/gorm"
)

type GroupDao struct {
	db *gorm.DB
}

func NewGroupDao(db *gorm.DB) *GroupDao {
	return &GroupDao{db: db}
}

// CreateGroup 创建群组
func (d *GroupDao) CreateGroup(group *models.Group) error {
	return d.db.Create(group).Error
}

// GetGroupByID 根据ID获取群组
func (d *GroupDao) GetGroupByID(groupID int64) (*models.Group, error) {
	var group models.Group
	err := d.db.Where("id = ?", groupID).First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &group, nil
}

// GetUserGroups 获取用户加入的群组列表
func (d *GroupDao) GetUserGroups(userID int64) ([]*models.Group, error) {
	var groups []*models.Group
	err := d.db.Table("`groups` AS g").
		Select("g.*").
		Joins("INNER JOIN group_members gm ON g.id = gm.group_id").
		Where("gm.user_id = ? AND gm.deleted_at IS NULL", userID).
		Order("g.updated_at DESC").
		Find(&groups).Error
	return groups, err
}

// UpdateGroup 更新群组信息
func (d *GroupDao) UpdateGroup(group *models.Group) error {
	return d.db.Model(group).Updates(group).Error
}

// AddMember 添加群成员（如果已存在则返回错误）
func (d *GroupDao) AddMember(groupID, userID int64, role string, nickname string) error {
	// 检查成员是否已存在
	var existingMember models.GroupMember
	err := d.db.Where("group_id = ? AND user_id = ? AND deleted_at IS NULL", groupID, userID).
		First(&existingMember).Error
	
	if err == nil {
		// 成员已存在
		return fmt.Errorf("用户已在群组中")
	}
	
	if err != gorm.ErrRecordNotFound {
		// 其他错误
		return err
	}

	// 创建新成员
	member := &models.GroupMember{
		GroupID:  groupID,
		UserID:   userID,
		Role:     role,
		Nickname: nickname,
		JoinedAt: time.Now().Unix(),
	}
	
	if err := d.db.Create(member).Error; err != nil {
		return err
	}

	// 更新群组成员数量
	return d.db.Model(&models.Group{}).
		Where("id = ?", groupID).
		UpdateColumn("member_count", gorm.Expr("member_count + 1")).Error
}

// AddMembers 批量添加群成员
func (d *GroupDao) AddMembers(groupID int64, userIDs []int64, inviterID int64) error {
	// 开始事务
	tx := d.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	addedCount := 0
	for _, userID := range userIDs {
		// 检查成员是否已存在
		var existingMember models.GroupMember
		err := tx.Where("group_id = ? AND user_id = ? AND deleted_at IS NULL", groupID, userID).
			First(&existingMember).Error
		
		if err == nil {
			// 成员已存在，跳过
			continue
		}
		
		if err != gorm.ErrRecordNotFound {
			tx.Rollback()
			return err
		}

		// 创建新成员
		member := &models.GroupMember{
			GroupID:  groupID,
			UserID:   userID,
			Role:     models.GroupRoleMember,
			Nickname: "",
			JoinedAt: time.Now().Unix(),
		}
		
		if err := tx.Create(member).Error; err != nil {
			tx.Rollback()
			return err
		}
		addedCount++
	}

	// 更新群组成员数量
	if addedCount > 0 {
		if err := tx.Model(&models.Group{}).
			Where("id = ?", groupID).
			UpdateColumn("member_count", gorm.Expr("member_count + ?", addedCount)).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// RemoveMember 移除群成员
func (d *GroupDao) RemoveMember(groupID, userID int64) error {
	// 软删除成员
	result := d.db.Model(&models.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Update("deleted_at", time.Now())
	
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("成员不存在")
	}

	// 更新群组成员数量
	return d.db.Model(&models.Group{}).
		Where("id = ?", groupID).
		UpdateColumn("member_count", gorm.Expr("member_count - 1")).Error
}

// GetGroupMembers 获取群成员列表
func (d *GroupDao) GetGroupMembers(groupID int64) ([]*models.GroupMemberInfo, error) {
	var members []*models.GroupMemberInfo
	err := d.db.Table("group_members gm").
		Select("gm.*, ub.name as user_name, ub.avatar as user_avatar, ub.phone as user_phone").
		Joins("LEFT JOIN user_basic ub ON gm.user_id = ub.id").
		Where("gm.group_id = ? AND gm.deleted_at IS NULL", groupID).
		Order("gm.role DESC, gm.joined_at ASC").
		Scan(&members).Error
	return members, err
}

// IsMemberInGroup 检查用户是否在群组中
func (d *GroupDao) IsMemberInGroup(groupID, userID int64) (bool, error) {
	var count int64
	err := d.db.Model(&models.GroupMember{}).
		Where("group_id = ? AND user_id = ? AND deleted_at IS NULL", groupID, userID).
		Count(&count).Error
	return count > 0, err
}

// GetMemberRole 获取成员角色
func (d *GroupDao) GetMemberRole(groupID, userID int64) (string, error) {
	var member models.GroupMember
	err := d.db.Where("group_id = ? AND user_id = ? AND deleted_at IS NULL", groupID, userID).
		First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("用户不在群组中")
		}
		return "", err
	}
	return member.Role, nil
}

// UpdateMemberNickname 更新成员群内昵称
func (d *GroupDao) UpdateMemberNickname(groupID, userID int64, nickname string) error {
	return d.db.Model(&models.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Update("nickname", nickname).Error
}

// GetGroupMemberCount 获取群成员数量
func (d *GroupDao) GetGroupMemberCount(groupID int64) (int64, error) {
	var count int64
	err := d.db.Model(&models.GroupMember{}).
		Where("group_id = ? AND deleted_at IS NULL", groupID).
		Count(&count).Error
	return count, err
}

// UpdateMemberRole 更新成员角色
func (d *GroupDao) UpdateMemberRole(groupID, userID int64, role string) error {
	return d.db.Model(&models.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Update("role", role).Error
}

// GetGroupMember 获取群成员信息
func (d *GroupDao) GetGroupMember(groupID, userID int64) (*models.GroupMember, error) {
	var member models.GroupMember
	err := d.db.Where("group_id = ? AND user_id = ? AND deleted_at IS NULL", groupID, userID).
		First(&member).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// DissolveGroup 解散群组（软删除）
func (d *GroupDao) DissolveGroup(groupID int64) error {
	now := time.Now().Unix()
	// 软删除群组
	err := d.db.Model(&models.Group{}).
		Where("id = ?", groupID).
		Update("deleted_at", now).Error
	if err != nil {
		return err
	}
	// 软删除所有成员
	return d.db.Model(&models.GroupMember{}).
		Where("group_id = ?", groupID).
		Update("deleted_at", now).Error
}

// TransferGroup 转让群组
func (d *GroupDao) TransferGroup(groupID, oldOwnerID, newOwnerID int64) error {
	// 使用事务确保原子性
	return d.db.Transaction(func(tx *gorm.DB) error {
		// 更新群组owner_id
		if err := tx.Model(&models.Group{}).
			Where("id = ?", groupID).
			Update("owner_id", newOwnerID).Error; err != nil {
			return err
		}

		// 更新原群主角色为member
		if err := tx.Model(&models.GroupMember{}).
			Where("group_id = ? AND user_id = ?", groupID, oldOwnerID).
			Update("role", models.GroupRoleMember).Error; err != nil {
			return err
		}

		// 更新新群主角色为owner
		if err := tx.Model(&models.GroupMember{}).
			Where("group_id = ? AND user_id = ?", groupID, newOwnerID).
			Update("role", models.GroupRoleOwner).Error; err != nil {
			return err
		}

		return nil
	})
}

