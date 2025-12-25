package dao

import (
	"context"
	"fmt"
	"time"

	"gochat/internal/infrastructure/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type ChatDao struct {
	mysqlDB *gorm.DB
	mongoDB *mongo.Database
}

func NewChatDao(mysqlDB *gorm.DB, mongoDB *mongo.Database) *ChatDao {
	return &ChatDao{
		mysqlDB: mysqlDB,
		mongoDB: mongoDB,
	}
}

// CreateMessage 创建消息（同时存储到MySQL和MongoDB）
func (d *ChatDao) CreateMessage(message *models.ChatMessage) error {
	// 生成消息ID
	message.MessageID = generateMessageID()
	message.CreatedAt = time.Now().Unix()
	message.UpdatedAt = time.Now().Unix()

	fmt.Printf("准备存储消息: MessageID=%s, FromUserID=%d, ToUserID=%d, Content=%s\n",
		message.MessageID, message.FromUserID, message.ToUserID, message.Content)

	// 存储到MongoDB（主要存储）
	collection := d.mongoDB.Collection("chat_messages")
	result, err := collection.InsertOne(context.Background(), message)
	if err != nil {
		fmt.Printf("MongoDB存储失败: %v\n", err)
		return fmt.Errorf("failed to insert message to MongoDB: %w", err)
	}
	fmt.Printf("MongoDB存储成功: InsertedID=%v\n", result.InsertedID)

	// 存储到MySQL（用于快速查询和统计）
	err = d.mysqlDB.Create(message).Error
	if err != nil {
		fmt.Printf("MySQL存储失败: %v\n", err)
		// MySQL存储失败不影响整体流程，但记录错误
		fmt.Printf("Warning: failed to insert message to MySQL: %v\n", err)
	} else {
		fmt.Printf("MySQL存储成功\n")
	}

	// 更新会话信息
	err = d.updateConversation(message)
	if err != nil {
		// 记录错误但不影响消息发送
		fmt.Printf("Failed to update conversation: %v\n", err)
	}

	return nil
}

// GetMessageHistory 获取消息历史（仅从MongoDB查询）
func (d *ChatDao) GetMessageHistory(fromUserID, toUserID, groupID int64, page, pageSize int, beforeTime *time.Time) ([]*models.ChatMessage, int64, error) {
	// 检查MongoDB连接
	if d.mongoDB == nil {
		return nil, 0, fmt.Errorf("MongoDB连接未初始化")
	}

	collection := d.mongoDB.Collection("chat_messages")

	fmt.Printf("MongoDB查询: fromUserID=%d, toUserID=%d, groupID=%d, page=%d, pageSize=%d\n", fromUserID, toUserID, groupID, page, pageSize)

	// 构建查询条件
	var filter bson.M

	if groupID > 0 {
		// 群聊：查询该群组的所有消息
		filter = bson.M{
			"group_id": groupID,
		}
	} else {
		// 私聊：查询双方的消息
		// 兼容旧数据：group_id 可能不存在
		groupFilter := []bson.M{
			{"group_id": 0},
			{"group_id": bson.M{"$exists": false}},
		}

		filter = bson.M{
			"$or": []bson.M{
				{
					"from_user_id": fromUserID,
					"to_user_id":   toUserID,
					"$or":          groupFilter,
				},
				{
					"from_user_id": toUserID,
					"to_user_id":   fromUserID,
					"$or":          groupFilter,
				},
			},
		}
	}

	// 如果指定了时间，只获取该时间之前的消息
	if beforeTime != nil {
		filter["created_at"] = bson.M{"$lt": beforeTime}
	}

	fmt.Printf("MongoDB查询条件: %+v\n", filter)

	// 计算总数
	totalCount, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		fmt.Printf("MongoDB计数失败: %v\n", err)
		return nil, 0, fmt.Errorf("failed to count messages: %w", err)
	}

	fmt.Printf("MongoDB查询到总数: %d\n", totalCount)

	// 构建查询选项（按时间正序）
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: 1}}). // 按时间正序（最早的在前）
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))

	// 执行查询
	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		fmt.Printf("MongoDB查询失败: %v\n", err)
		return nil, 0, fmt.Errorf("failed to find messages: %w", err)
	}
	defer cursor.Close(context.Background())

	var messages []*models.ChatMessage
	for cursor.Next(context.Background()) {
		var message models.ChatMessage
		if err := cursor.Decode(&message); err != nil {
			fmt.Printf("MongoDB解码消息失败: %v\n", err)
			return nil, 0, fmt.Errorf("failed to decode message: %w", err)
		}
		messages = append(messages, &message)
		fmt.Printf("解码消息: ID=%s, From=%d, To=%d, Content=%s, CreatedAt=%v\n",
			message.MessageID, message.FromUserID, message.ToUserID, message.Content, message.CreatedAt)
	}

	if err := cursor.Err(); err != nil {
		fmt.Printf("MongoDB游标错误: %v\n", err)
		return nil, 0, fmt.Errorf("cursor error: %w", err)
	}

	fmt.Printf("MongoDB返回消息数量: %d\n", len(messages))

	return messages, totalCount, nil
}

// RecallMessage 撤回消息
func (d *ChatDao) RecallMessage(messageID string, userID int64) error {
	collection := d.mongoDB.Collection("chat_messages")
	
	// 查询消息是否存在且是发送者
	filter := bson.M{
		"message_id":  messageID,
		"from_user_id": userID,
	}
	
	update := bson.M{
		"$set": bson.M{
			"is_recalled": true,
			"updated_at":  time.Now().Unix(),
		},
	}
	
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to recall message: %w", err)
	}
	
	if result.MatchedCount == 0 {
		return fmt.Errorf("message not found or no permission")
	}
	
	// 同时更新MySQL
	err = d.mysqlDB.Model(&models.ChatMessage{}).
		Where("message_id = ? AND from_user_id = ?", messageID, userID).
		Updates(map[string]interface{}{
			"is_recalled": true,
			"updated_at":  time.Now().Unix(),
		}).Error
	
	return err
}

// DeleteMessage 删除消息
func (d *ChatDao) DeleteMessage(messageID string, userID int64, deleteForBoth bool) error {
	collection := d.mongoDB.Collection("chat_messages")
	
	// 查询消息
	var message models.ChatMessage
	filter := bson.M{"message_id": messageID}
	err := collection.FindOne(context.Background(), filter).Decode(&message)
	if err != nil {
		return fmt.Errorf("message not found: %w", err)
	}
	
	// 检查权限：只有发送者或接收者可以删除
	if message.FromUserID != userID && message.ToUserID != userID {
		return fmt.Errorf("no permission to delete message")
	}
	
	// 如果是群聊，不支持双向删除
	if message.GroupID > 0 {
		deleteForBoth = false
	}
	
	// 更新消息
	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"updated_at": time.Now().Unix(),
		},
	}
	
	// 如果是双向删除，需要更新双方的消息
	if deleteForBoth && message.GroupID == 0 {
		// 私聊双向删除：更新发送者和接收者的消息
		filter = bson.M{
			"message_id": messageID,
		}
	} else {
		// 单方删除：只删除当前用户的消息
		filter = bson.M{
			"message_id": messageID,
			"$or": []bson.M{
				{"from_user_id": userID},
				{"to_user_id": userID},
			},
		}
	}
	
	_, err = collection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	
	// 同时更新MySQL
	mysqlFilter := d.mysqlDB.Model(&models.ChatMessage{}).
		Where("message_id = ?", messageID)
	
	if !deleteForBoth || message.GroupID > 0 {
		mysqlFilter = mysqlFilter.Where("from_user_id = ? OR to_user_id = ?", userID, userID)
	}
	
	err = mysqlFilter.Updates(map[string]interface{}{
		"is_deleted": true,
		"updated_at": time.Now().Unix(),
	}).Error
	
	return err
}

// GetUnreadCount 获取未读消息数量
func (d *ChatDao) GetUnreadCount(userID int64) (map[int64]int, error) {
	collection := d.mongoDB.Collection("chat_messages")

	// 查询发送给当前用户且未读的消息
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{
			"to_user_id": userID,
			"status":     bson.M{"$lt": models.MessageStatusRead},
		}}},
		{{"$group", bson.M{
			"_id":   "$from_user_id",
			"count": bson.M{"$sum": 1},
		}}},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate unread count: %w", err)
	}
	defer cursor.Close(context.Background())

	unreadCount := make(map[int64]int)
	for cursor.Next(context.Background()) {
		var result struct {
			FromUserID int64 `bson:"_id"`
			Count      int   `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode unread count: %v", err)
		}
		unreadCount[result.FromUserID] = result.Count
	}

	return unreadCount, nil
}

// MarkMessagesAsRead 标记消息为已读
func (d *ChatDao) MarkMessagesAsRead(fromUserID, toUserID int64, messageID string) (int64, error) {
	collection := d.mongoDB.Collection("chat_messages")

	// 构建查询条件
	filter := bson.M{
		"from_user_id": fromUserID,
		"to_user_id":   toUserID,
		"status":       bson.M{"$lt": models.MessageStatusRead},
	}

	// 如果指定了消息ID，只标记该消息
	if messageID != "" {
		filter["message_id"] = messageID
	}

	// 更新消息状态
	update := bson.M{
		"$set": bson.M{
			"status":     models.MessageStatusRead,
			"updated_at": time.Now().Unix(),
		},
	}

	result, err := collection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return 0, fmt.Errorf("failed to mark messages as read: %w", err)
	}

	// 更新会话的未读数量
	err = d.updateConversationUnreadCount(fromUserID, toUserID)
	if err != nil {
		fmt.Printf("Failed to update conversation unread count: %v\n", err)
	}

	// 更新会话设置的未读数量
	convDao := NewConversationDao(d.mysqlDB)
	unreadCount, err := d.GetUnreadCount(toUserID)
	if err == nil {
		if count, ok := unreadCount[fromUserID]; ok {
			convDao.UpdateUnreadCount(toUserID, fromUserID, 0, count)
		} else {
			convDao.UpdateUnreadCount(toUserID, fromUserID, 0, 0)
		}
	}

	return result.ModifiedCount, nil
}

// SearchUsers 搜索用户
func (d *ChatDao) SearchUsers(keyword string, loginAccount string, onlineOnly bool, onlineUserIDs []int64, limit int) ([]*models.UserBasic, error) {
	query := d.mysqlDB.Model(&models.UserBasic{})

	// 精准搜索登录账号（优先级最高）
	if loginAccount != "" {
		query = query.Where("login_account = ?", loginAccount)
	} else if keyword != "" {
		// 关键词搜索：昵称/手机号/邮箱
		query = query.Where("name LIKE ? OR phone LIKE ? OR email LIKE ? OR email1 LIKE ? OR email2 LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 在线状态筛选
	if onlineOnly && len(onlineUserIDs) > 0 {
		query = query.Where("id IN ?", onlineUserIDs)
	}

	var users []*models.UserBasic
	err := query.Limit(limit).Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	return users, nil
}

// GetConversations 获取会话列表
func (d *ChatDao) GetConversations(userID int64, page, pageSize int) ([]*models.Conversation, int64, error) {
	var conversations []*models.Conversation
	var totalCount int64

	// 获取总数
	err := d.mysqlDB.Model(&models.Conversation{}).Where("user_id = ?", userID).Count(&totalCount).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count conversations: %w", err)
	}

	// 获取会话列表
	err = d.mysqlDB.Where("user_id = ?", userID).
		Order("last_message_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&conversations).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get conversations: %w", err)
	}

	return conversations, totalCount, nil
}

// updateConversation 更新会话信息
func (d *ChatDao) updateConversation(message *models.ChatMessage) error {
	// 更新发送者的会话
	err := d.updateUserConversation(message.FromUserID, message.ToUserID, message)
	if err != nil {
		return err
	}

	// 更新接收者的会话
	err = d.updateUserConversation(message.ToUserID, message.FromUserID, message)
	if err != nil {
		return err
	}

	return nil
}

// updateUserConversation 更新单个用户的会话
func (d *ChatDao) updateUserConversation(userID, otherUserID int64, message *models.ChatMessage) error {
	var conversation models.Conversation

	// 查找或创建会话
	err := d.mysqlDB.Where("user_id = ? AND other_user_id = ?", userID, otherUserID).
		First(&conversation).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新会话
		conversation = models.Conversation{
			UserID:          userID,
			OtherUserID:     otherUserID,
			LastMessage:     message.Content,
			LastMessageType: message.MessageType,
			UnreadCount:     0,
			LastMessageAt:   message.CreatedAt,
		}
		err = d.mysqlDB.Create(&conversation).Error
	} else if err == nil {
		// 更新现有会话
		updates := map[string]interface{}{
			"last_message":      message.Content,
			"last_message_type": message.MessageType,
			"last_message_at":   message.CreatedAt,
		}

		// 如果是接收者，增加未读数量
		if userID == message.ToUserID {
			updates["unread_count"] = gorm.Expr("unread_count + 1")
			// 同时更新会话设置的未读数量
			convDao := NewConversationDao(d.mysqlDB)
			setting, _ := convDao.GetOrCreateConversationSetting(userID, otherUserID, 0)
			if setting != nil {
				convDao.UpdateUnreadCount(userID, otherUserID, 0, setting.UnreadCount+1)
			}
		}

		err = d.mysqlDB.Model(&conversation).Updates(updates).Error
	}

	return err
}

// updateConversationUnreadCount 更新会话未读数量
func (d *ChatDao) updateConversationUnreadCount(fromUserID, toUserID int64) error {
	// 重新计算未读数量
	unreadCount, err := d.GetUnreadCount(toUserID)
	if err != nil {
		return err
	}

	// 更新会话的未读数量
	return d.mysqlDB.Model(&models.Conversation{}).
		Where("user_id = ? AND other_user_id = ?", toUserID, fromUserID).
		Update("unread_count", unreadCount[fromUserID]).Error
}

// GetMessagesByIDs 根据消息ID列表获取消息
func (d *ChatDao) GetMessagesByIDs(messageIDs []string) ([]*models.ChatMessage, error) {
	if len(messageIDs) == 0 {
		return nil, fmt.Errorf("message IDs list is empty")
	}

	collection := d.mongoDB.Collection("chat_messages")

	// 构建查询条件
	filter := bson.M{
		"message_id": bson.M{"$in": messageIDs},
	}

	// 执行查询
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find messages: %w", err)
	}
	defer cursor.Close(context.Background())

	var messages []*models.ChatMessage
	for cursor.Next(context.Background()) {
		var message models.ChatMessage
		if err := cursor.Decode(&message); err != nil {
			return nil, fmt.Errorf("failed to decode message: %w", err)
		}
		messages = append(messages, &message)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return messages, nil
}

// generateMessageID 生成消息ID
func generateMessageID() string {
	return primitive.NewObjectID().Hex()
}
