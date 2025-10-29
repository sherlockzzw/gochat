package dao

import (
	"gochat/models"

	"gorm.io/gorm"
)

type FriendDao struct {
	db *gorm.DB
}

func NewFriendDao(db *gorm.DB) *FriendDao {
	return &FriendDao{db: db}
}

// CreateFriendRequest 创建好友申请
func (d *FriendDao) CreateFriendRequest(req *models.FriendRequest) error {
	return d.db.Create(req).Error
}

// GetFriendRequestByID 根据ID获取好友申请
func (d *FriendDao) GetFriendRequestByID(requestID uint) (*models.FriendRequest, error) {
	var request models.FriendRequest
	err := d.db.Where("id = ?", requestID).First(&request).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &request, nil
}

// GetFriendRequest 获取好友申请
func (d *FriendDao) GetFriendRequest(fromUserID, toUserID uint) (*models.FriendRequest, error) {
	var req models.FriendRequest
	err := d.db.Where("from_user_id = ? AND to_user_id = ? AND status = ?", fromUserID, toUserID, "pending").First(&req).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &req, nil
}

// UpdateFriendRequestStatus 更新好友申请状态
func (d *FriendDao) UpdateFriendRequestStatus(requestID uint, status string) error {
	return d.db.Model(&models.FriendRequest{}).Where("id = ?", requestID).Update("status", status).Error
}

// GetFriendRequests 获取好友申请列表
func (d *FriendDao) GetFriendRequests(userID uint, page, pageSize int) ([]*models.FriendRequest, error) {
	var requests []*models.FriendRequest
	offset := (page - 1) * pageSize

	err := d.db.Where("to_user_id = ? AND status = ?", userID, "pending").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&requests).Error

	return requests, err
}

// GetFriendRequestsWithUserInfo 获取好友申请列表（包含申请人信息）
func (d *FriendDao) GetFriendRequestsWithUserInfo(userID uint, page, pageSize int) ([]*models.FriendRequestWithUser, error) {
	var requests []*models.FriendRequestWithUser
	offset := (page - 1) * pageSize

	err := d.db.Table("friend_requests fr").
		Select("fr.*, ub.name as from_user_name, ub.avatar as from_user_avatar").
		Joins("LEFT JOIN user_basic ub ON fr.from_user_id = ub.id").
		Where("fr.to_user_id = ? AND fr.status = ?", userID, "pending").
		Order("fr.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Scan(&requests).Error

	return requests, err
}

// CreateFriend 创建好友关系
func (d *FriendDao) CreateFriend(friend *models.Friend) error {
	return d.db.Create(friend).Error
}

// CreateFriendsInTransaction 在事务中创建双向好友关系
func (d *FriendDao) CreateFriendsInTransaction(friend1, friend2 *models.Friend) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		// 创建第一条好友关系
		if err := tx.Create(friend1).Error; err != nil {
			return err
		}

		// 创建第二条好友关系
		if err := tx.Create(friend2).Error; err != nil {
			return err
		}

		return nil
	})
}

// CheckIsFriend 检查是否为好友
func (d *FriendDao) CheckIsFriend(userID, friendID uint) (*models.Friend, error) {
	var friend models.Friend
	err := d.db.Where("user_id = ? AND friend_id = ?", userID, friendID).First(&friend).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &friend, nil
}

// GetFriendList 获取好友列表
func (d *FriendDao) GetFriendList(userID uint, page, pageSize int) ([]*models.Friend, error) {
	var friends []*models.Friend
	offset := (page - 1) * pageSize

	err := d.db.Where("user_id = ? AND is_blocked = ?", userID, false).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&friends).Error

	return friends, err
}

// GetFriendListWithUserInfo 获取好友列表（包含用户信息）
func (d *FriendDao) GetFriendListWithUserInfo(userID uint, page, pageSize int) ([]*models.FriendWithUser, error) {
	var friends []*models.FriendWithUser
	offset := (page - 1) * pageSize

	err := d.db.Table("friends f").
		Select("f.*, ub.name as friend_name, ub.phone as friend_phone, ub.email as friend_email, ub.avatar as friend_avatar").
		Joins("LEFT JOIN user_basic ub ON f.friend_id = ub.id").
		Where("f.user_id = ? AND f.is_blocked = ?", userID, false).
		Order("f.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Scan(&friends).Error

	return friends, err
}

// GetFriendListWithUserInfoAndOnlineStatus 获取好友列表（包含用户信息和在线状态）
func (d *FriendDao) GetFriendListWithUserInfoAndOnlineStatus(userID uint, page, pageSize int, onlineUserIDs []uint) ([]*models.FriendWithUser, error) {
	var friends []*models.FriendWithUser
	offset := (page - 1) * pageSize

	onlineMap := make(map[uint]bool)
	for _, id := range onlineUserIDs {
		onlineMap[id] = true
	}

	err := d.db.Table("friends f").
		Select("f.*, ub.name as friend_name, ub.phone as friend_phone, ub.email as friend_email, ub.avatar as friend_avatar").
		Joins("LEFT JOIN user_basic ub ON f.friend_id = ub.id").
		Where("f.user_id = ? AND f.is_blocked = ?", userID, false).
		Order("f.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Scan(&friends).Error

	if err != nil {
		return nil, err
	}

	// 设置在线状态
	for _, friend := range friends {
		friend.IsOnline = onlineMap[friend.FriendID]
	}

	return friends, err
}

// DeleteFriend 删除好友
func (d *FriendDao) DeleteFriend(userID, friendID uint) error {
	// 删除双向好友关系
	err := d.db.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
		userID, friendID, friendID, userID).Delete(&models.Friend{}).Error
	return err
}

// GetFriendDetail 获取好友详情
func (d *FriendDao) GetFriendDetail(userID, friendID uint) (*models.Friend, error) {
	var friend models.Friend
	err := d.db.Where("user_id = ? AND friend_id = ?", userID, friendID).First(&friend).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &friend, nil
}

// SetFriendRemark 设置好友备注
func (d *FriendDao) SetFriendRemark(userID, friendID uint, remark string) error {
	return d.db.Model(&models.Friend{}).
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		Update("remark", remark).Error
}

// BlockFriend 屏蔽好友
func (d *FriendDao) BlockFriend(userID, friendID uint, isBlocked bool) error {
	return d.db.Model(&models.Friend{}).
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		Update("is_blocked", isBlocked).Error
}
