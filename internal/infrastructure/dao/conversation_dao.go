package dao

import (
	"gochat/internal/infrastructure/models"
	"time"

	"gorm.io/gorm"
)

type ConversationDao struct {
	db *gorm.DB
}

func NewConversationDao(db *gorm.DB) *ConversationDao {
	return &ConversationDao{db: db}
}

// GetOrCreateConversationSetting 获取或创建会话设置
func (d *ConversationDao) GetOrCreateConversationSetting(userID, otherUserID, groupID int64) (*models.ConversationSetting, error) {
	var setting models.ConversationSetting

	// 构建查询条件
	query := d.db.Where("user_id = ?", userID)
	if groupID > 0 {
		query = query.Where("group_id = ? AND other_user_id = 0", groupID)
	} else {
		query = query.Where("other_user_id = ? AND group_id = 0", otherUserID)
	}

	err := query.First(&setting).Error
	if err == gorm.ErrRecordNotFound {
		// 创建新设置
		now := time.Now().Unix()
		setting = models.ConversationSetting{
			UserID:      userID,
			OtherUserID: otherUserID,
			GroupID:     groupID,
			IsPinned:    false,
			IsMuted:     false,
			IsHidden:    false,
			UnreadCount: 0,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		err = d.db.Create(&setting).Error
		if err != nil {
			return nil, err
		}
		return &setting, nil
	} else if err != nil {
		return nil, err
	}

	return &setting, nil
}

// GetConversationSetting 获取会话设置
func (d *ConversationDao) GetConversationSetting(userID, otherUserID, groupID int64) (*models.ConversationSetting, error) {
	var setting models.ConversationSetting

	query := d.db.Where("user_id = ?", userID)
	if groupID > 0 {
		query = query.Where("group_id = ? AND other_user_id = 0", groupID)
	} else {
		query = query.Where("other_user_id = ? AND group_id = 0", otherUserID)
	}

	err := query.First(&setting).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &setting, nil
}

// UpdateConversationSetting 更新会话设置
func (d *ConversationDao) UpdateConversationSetting(userID, otherUserID, groupID int64, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now().Unix()

	query := d.db.Model(&models.ConversationSetting{}).Where("user_id = ?", userID)
	if groupID > 0 {
		query = query.Where("group_id = ? AND other_user_id = 0", groupID)
	} else {
		query = query.Where("other_user_id = ? AND group_id = 0", otherUserID)
	}

	return query.Updates(updates).Error
}

// GetConversationsWithSettings 获取会话列表（包含设置信息）
func (d *ConversationDao) GetConversationsWithSettings(userID int64, page, pageSize int32, includeHidden bool) ([]*ConversationWithSetting, int64, error) {
	// 先获取会话列表
	var conversations []*models.Conversation
	query := d.db.Model(&models.Conversation{}).Where("user_id = ?", userID)

	// 如果不包含隐藏的会话，需要关联查询排除隐藏的会话
	if !includeHidden {
		// 通过子查询排除隐藏的会话（Conversation模型只支持私聊）
		query = query.Where("NOT EXISTS (SELECT 1 FROM conversation_settings cs WHERE cs.user_id = ? AND cs.other_user_id = conversations.other_user_id AND cs.group_id = 0 AND cs.is_hidden = 1)", userID)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("last_message_at DESC").
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&conversations).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取每个会话的设置
	result := make([]*ConversationWithSetting, 0, len(conversations))
	for _, conv := range conversations {
		setting, _ := d.GetConversationSetting(userID, conv.OtherUserID, 0) // 私聊，groupID为0
		if setting == nil {
			// 如果没有设置，创建默认设置
			setting, _ = d.GetOrCreateConversationSetting(userID, conv.OtherUserID, 0)
		}

		result = append(result, &ConversationWithSetting{
			Conversation: conv,
			Setting:      setting,
		})
	}

	return result, total, nil
}

// ConversationWithSetting 会话和设置的组合
type ConversationWithSetting struct {
	Conversation *models.Conversation
	Setting      *models.ConversationSetting
}

// UpdateUnreadCount 更新未读数量
func (d *ConversationDao) UpdateUnreadCount(userID, otherUserID, groupID int64, count int) error {
	// 获取或创建设置
	setting, err := d.GetOrCreateConversationSetting(userID, otherUserID, groupID)
	if err != nil {
		return err
	}

	// 更新未读数量
	return d.db.Model(setting).Update("unread_count", count).Error
}

// ClearUnreadCount 清除未读数量
func (d *ConversationDao) ClearUnreadCount(userID, otherUserID, groupID int64) (int, error) {
	setting, err := d.GetOrCreateConversationSetting(userID, otherUserID, groupID)
	if err != nil {
		return 0, err
	}

	clearedCount := setting.UnreadCount
	err = d.db.Model(setting).Update("unread_count", 0).Error
	if err != nil {
		return 0, err
	}

	// 同时更新Conversation表的未读数量
	convQuery := d.db.Model(&models.Conversation{}).Where("user_id = ?", userID)
	if groupID > 0 {
		// 群聊（暂时不支持，因为Conversation模型只支持私聊）
	} else {
		convQuery = convQuery.Where("other_user_id = ?", otherUserID)
		convQuery.Update("unread_count", 0)
	}

	return clearedCount, nil
}

// DeleteConversation 删除会话（软删除）
func (d *ConversationDao) DeleteConversation(userID, otherUserID, groupID int64) error {
	// 删除会话记录
	convQuery := d.db.Model(&models.Conversation{}).Where("user_id = ?", userID)
	if groupID > 0 {
		// 群聊（暂时不支持，因为Conversation模型只支持私聊）
		return nil
	} else {
		convQuery = convQuery.Where("other_user_id = ?", otherUserID)
	}

	return convQuery.Delete(&models.Conversation{}).Error
}

