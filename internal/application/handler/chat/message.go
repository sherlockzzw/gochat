package chat

import (
	"encoding/json"
	"fmt"
	"gochat/api/api/chat"
	"gochat/internal/application/handler/common"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"time"

	"github.com/gin-gonic/gin"
)

// SendMessage 发送消息
func (h *ChatHandler) SendMessage(ctx *gin.Context) {
	req, err := analysis.BindParameter[chat.SendMessageRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.sendMessageLogic(ctx, req)
	if code != 0 {
		h.response.JsonErrorFixation(ctx, code)
		return
	}
	if err != nil {
		h.response.JsonError(ctx, err, err.Error())
		return
	}

	h.response.JsonSuccess(ctx, resp)
}

func (h *ChatHandler) sendMessageLogic(ctx *gin.Context, req chat.SendMessageRequest) (resp *chat.SendMessageResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	fromUserID, errCode, err := common.GetUserIDFromContext(ctx)
	if errCode != 0 {
		return nil, errCode, err
	}

	// 验证消息类型：私聊或群聊二选一
	groupID := int64(req.GetGroupId())
	toUserID := int64(req.GetToUserId())

	if groupID == 0 && toUserID == 0 {
		return nil, code_msg.BadRequest, nil
	}

	if groupID > 0 && toUserID > 0 {
		return nil, code_msg.BadRequest, nil
	}

	// 如果是私聊，验证是否为好友关系（删除好友后不能发送新消息）
	if toUserID > 0 {
		friendRelation, err := h.friendDao.CheckIsFriend(fromUserID, toUserID)
		if err != nil {
			return nil, code_msg.ServerError, err
		}
		if friendRelation == nil {
			return nil, code_msg.NotFriend, nil
		}
		// 检查是否被屏蔽
		if friendRelation.IsBlocked {
			return nil, code_msg.BadRequest, nil
		}
	}

	// 如果是群聊，验证用户是否在群组中
	if groupID > 0 {
		isMember, err := h.groupDao.IsMemberInGroup(groupID, fromUserID)
		if err != nil {
			return nil, code_msg.ServerError, err
		}
		if !isMember {
			return nil, code_msg.NotGroupMember, nil
		}

		// 检查是否被禁言
		isMuted, err := h.muteDao.IsMuted(groupID, fromUserID)
		if err != nil {
			return nil, code_msg.ServerError, err
		}
		if isMuted {
			return nil, code_msg.BadRequest, nil
		}
	}

	// 验证消息类型相关的业务逻辑
	validateCode, validateErr := h.validateMessageType(&req, fromUserID)
	if validateCode != 0 {
		return nil, validateCode, validateErr
	}
	if validateErr != nil {
		return nil, code_msg.ServerError, validateErr
	}

	// 处理合并消息ID列表
	var mergeMessagesJSON string
	if len(req.GetMergeMessageIds()) > 0 {
		mergeMessagesBytes, _ := json.Marshal(req.GetMergeMessageIds())
		mergeMessagesJSON = string(mergeMessagesBytes)
	}

	// 创建消息模型
	now := time.Now().Unix()
	message := &models.ChatMessage{
		FromUserID:     fromUserID,
		ToUserID:       toUserID,
		GroupID:        groupID,
		MessageType:    int(req.GetMessageType()),
		Content:        req.GetContent(),
		FileURL:        req.GetFileUrl(),
		FileName:       req.GetFileName(),
		FileSize:       req.GetFileSize(),
		VideoURL:       req.GetVideoUrl(),
		VideoThumb:     req.GetVideoThumb(),
		VoiceURL:       req.GetVoiceUrl(),
		VoiceDuration:  int(req.GetVoiceDuration()),
		EmojiURL:       req.GetEmojiUrl(),
		MergeMessages:  mergeMessagesJSON,
		QuoteMessageID: req.GetQuoteMessageId(),
		ContactUserID:  int64(req.GetContactUserId()),
		RedPacketID:    int64(req.GetRedPacketId()),
		TransferID:     int64(req.GetTransferId()),
		ReadStatus:     1, // 默认已发送
		Status:         models.MessageStatusSent,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// 保存消息
	err = h.dao.CreateMessage(message)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 通过WebSocket实时推送消息
	h.sendMessageViaWebSocket(message)

	// 转换为响应格式（包含详细信息）
	chatMessage := h.convertToChatMessageWithDetails(message)

	resp = &chat.SendMessageResponse{
		Message: chatMessage,
		Success: true,
	}

	return resp, 0, nil
}

// GetMessageHistory 获取消息历史
func (h *ChatHandler) GetMessageHistory(ctx *gin.Context) {
	req, err := analysis.BindQuery[chat.GetMessageHistoryRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.getMessageHistoryLogic(ctx, req)
	if code != 0 {
		h.response.JsonErrorFixation(ctx, code)
		return
	}
	if err != nil {
		h.response.JsonError(ctx, err, err.Error())
		return
	}

	// 手动构建响应结构体，确保messages字段始终存在（即使为空数组）
	// 将protobuf消息转换为map，确保JSON序列化正确
	messagesList := make([]map[string]interface{}, 0)
	for _, msg := range resp.GetMessages() {
		msgMap := map[string]interface{}{
			"id":               msg.GetId(),
			"from_user_id":     msg.GetFromUserId(),
			"to_user_id":       msg.GetToUserId(),
			"group_id":         msg.GetGroupId(),
			"message_type":     msg.GetMessageType(),
			"content":          msg.GetContent(),
			"file_url":         msg.GetFileUrl(),
			"file_name":        msg.GetFileName(),
			"file_size":        msg.GetFileSize(),
			"video_url":        msg.GetVideoUrl(),
			"video_thumb":      msg.GetVideoThumb(),
			"voice_url":        msg.GetVoiceUrl(),
			"voice_duration":   msg.GetVoiceDuration(),
			"emoji_url":        msg.GetEmojiUrl(),
			"merge_messages":   msg.GetMergeMessages(),
			"quote_message_id": msg.GetQuoteMessageId(),
			"contact_user_id":  msg.GetContactUserId(),
			"red_packet_id":    msg.GetRedPacketId(),
			"transfer_id":      msg.GetTransferId(),
			"is_recalled":      msg.GetIsRecalled(),
			"is_deleted":       msg.GetIsDeleted(),
			"read_status":      msg.GetReadStatus(),
			"status":           msg.GetStatus(),
		}

		// 处理时间戳
		if msg.GetCreatedAt() != nil {
			msgMap["created_at"] = msg.GetCreatedAt().AsTime().Unix()
		}
		if msg.GetUpdatedAt() != nil {
			msgMap["updated_at"] = msg.GetUpdatedAt().AsTime().Unix()
		}

		messagesList = append(messagesList, msgMap)
	}

	responseData := map[string]interface{}{
		"messages":     messagesList,
		"total_count":  resp.GetTotalCount(),
		"current_page": resp.GetCurrentPage(),
		"page_size":    resp.GetPageSize(),
		"has_more":     resp.GetHasMore(),
	}

	h.response.JsonSuccess(ctx, responseData)
}

func (h *ChatHandler) getMessageHistoryLogic(ctx *gin.Context, req chat.GetMessageHistoryRequest) (resp *chat.GetMessageHistoryResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	fromUserID, errCode, err := common.GetUserIDFromContext(ctx)
	if errCode != 0 {
		return nil, errCode, err
	}

	// 解析时间参数
	var beforeTime *time.Time
	if req.GetBeforeTime() != nil {
		t := req.GetBeforeTime().AsTime()
		beforeTime = &t
	}

	// 判断是私聊还是群聊
	groupID := int64(req.GetGroupId())
	otherUserID := int64(req.GetOtherUserId())

	if groupID == 0 && otherUserID == 0 {
		return &chat.GetMessageHistoryResponse{
			Messages:    []*chat.ChatMessage{},
			TotalCount:  0,
			CurrentPage: req.GetPage(),
			PageSize:    req.GetPageSize(),
			HasMore:     false,
		}, 0, nil
	}

	// 如果是群聊，验证用户是否在群组中
	if groupID > 0 {
		isMember, err := h.groupDao.IsMemberInGroup(groupID, fromUserID)
		if err != nil {
			return nil, code_msg.ServerError, err
		}
		if !isMember {
			return &chat.GetMessageHistoryResponse{
				Messages:    []*chat.ChatMessage{},
				TotalCount:  0,
				CurrentPage: req.GetPage(),
				PageSize:    req.GetPageSize(),
				HasMore:     false,
			}, 0, nil
		}
	}

	// 获取消息历史
	fmt.Printf("准备查询消息历史: fromUserID=%d, otherUserID=%d, groupID=%d, page=%d, pageSize=%d\n",
		fromUserID, otherUserID, groupID, req.GetPage(), req.GetPageSize())

	messages, totalCount, err := h.dao.GetMessageHistory(
		fromUserID,
		otherUserID,
		groupID,
		int(req.GetPage()),
		int(req.GetPageSize()),
		beforeTime,
	)
	if err != nil {
		fmt.Printf("查询消息历史失败: %v\n", err)
		return nil, code_msg.ServerError, err
	}

	fmt.Printf("从数据库查询到 %d 条消息，总数: %d\n", len(messages), totalCount)

	// 转换为响应格式（包含详细信息）
	var chatMessages []*chat.ChatMessage
	for _, msg := range messages {
		chatMsg := h.convertToChatMessageWithDetails(msg)
		chatMessages = append(chatMessages, chatMsg)
	}

	fmt.Printf("转换后的消息数量: %d\n", len(chatMessages))

	hasMore := int64(req.GetPage()*req.GetPageSize()) < totalCount

	// 确保Messages不为nil
	if chatMessages == nil {
		chatMessages = []*chat.ChatMessage{}
	}

	// 构建响应（使用protobuf结构体）
	resp = &chat.GetMessageHistoryResponse{
		Messages:    chatMessages,
		TotalCount:  int32(totalCount),
		CurrentPage: req.GetPage(),
		PageSize:    req.GetPageSize(),
		HasMore:     hasMore,
	}

	fmt.Printf("构建响应: Messages数量=%d, TotalCount=%d, CurrentPage=%d, PageSize=%d\n",
		len(resp.Messages), resp.TotalCount, resp.CurrentPage, resp.PageSize)

	// 测试JSON序列化
	respJSON, _ := json.Marshal(resp)
	fmt.Printf("响应JSON序列化结果: %s\n", string(respJSON))

	return resp, 0, nil
}

// MarkMessageRead 标记消息为已读
func (h *ChatHandler) MarkMessageRead(ctx *gin.Context) {
	req, err := analysis.BindParameter[chat.MarkMessageReadRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.markMessageReadLogic(ctx, req)
	if code != 0 {
		h.response.JsonErrorFixation(ctx, code)
		return
	}
	if err != nil {
		h.response.JsonError(ctx, err, err.Error())
		return
	}

	h.response.JsonSuccess(ctx, resp)
}

func (h *ChatHandler) markMessageReadLogic(ctx *gin.Context, req chat.MarkMessageReadRequest) (resp *chat.MarkMessageReadResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	toUserID, errCode, err := common.GetUserIDFromContext(ctx)
	if errCode != 0 {
		return nil, errCode, err
	}

	// 标记消息为已读
	markedCount, err := h.dao.MarkMessagesAsRead(
		int64(req.GetOtherUserId()),
		toUserID,
		req.GetMessageId(),
	)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	resp = &chat.MarkMessageReadResponse{
		Success:      true,
		MarkedCount:  int32(markedCount),
		ErrorMessage: "",
	}

	return resp, 0, nil
}

// GetUnreadCount 获取未读消息数量
func (h *ChatHandler) GetUnreadCount(ctx *gin.Context) {
	req, err := analysis.BindQuery[chat.GetUnreadCountRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.getUnreadCountLogic(ctx, req)
	if code != 0 {
		h.response.JsonErrorFixation(ctx, code)
		return
	}
	if err != nil {
		h.response.JsonError(ctx, err, err.Error())
		return
	}

	h.response.JsonSuccess(ctx, resp)
}

func (h *ChatHandler) getUnreadCountLogic(ctx *gin.Context, req chat.GetUnreadCountRequest) (resp *chat.GetUnreadCountResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, errCode, err := common.GetUserIDFromContext(ctx)
	if errCode != 0 {
		return nil, errCode, err
	}

	// 获取未读消息数量
	unreadByUser, err := h.dao.GetUnreadCount(userID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 计算总未读数量
	totalUnread := 0
	unreadByUserMap := make(map[uint32]int32)
	for userID, count := range unreadByUser {
		totalUnread += count
		unreadByUserMap[uint32(userID)] = int32(count)
	}

	resp = &chat.GetUnreadCountResponse{
		TotalUnread:  int32(totalUnread),
		UnreadByUser: unreadByUserMap,
	}

	return resp, 0, nil
}

// sendMessageViaWebSocket 通过WebSocket发送消息
func (h *ChatHandler) sendMessageViaWebSocket(message *models.ChatMessage) {
	// 构建WebSocket消息
	wsMessage := map[string]interface{}{
		"type": "message",
		"data": map[string]interface{}{
			"id":           message.MessageID,
			"from_user_id": message.FromUserID,
			"to_user_id":   message.ToUserID,
			"group_id":     message.GroupID,
			"message_type": message.MessageType,
			"content":      message.Content,
			"file_url":     message.FileURL,
			"file_name":    message.FileName,
			"file_size":    message.FileSize,
			"status":       message.Status,
			"created_at":   message.CreatedAt,
		},
	}

	// 将消息转换为JSON
	messageBytes, err := json.Marshal(wsMessage)
	if err != nil {
		fmt.Printf("Failed to marshal WebSocket message: %v\n", err)
		return
	}

	// 判断是私聊还是群聊
	if message.GroupID > 0 {
		// 群聊：推送给所有群成员（包括发送者自己）
		members, err := h.groupDao.GetGroupMembers(message.GroupID)
		if err != nil {
			fmt.Printf("获取群成员失败: %v\n", err)
			return
		}

		if len(members) == 0 {
			fmt.Printf("警告: 群组 %d 没有成员，无法推送消息\n", message.GroupID)
			return
		}

		fmt.Printf("准备通过WebSocket推送群聊消息给群组 %d 的 %d 个成员，消息长度: %d\n",
			message.GroupID, len(members), len(messageBytes))

		// 统计推送结果
		successCount := 0
		onlineCount := 0
		offlineCount := 0

		// 获取在线用户列表
		onlineUserIDs := make(map[int64]bool)
		onlineIDs := common.GetOnlineUserIDs()
		for _, uid := range onlineIDs {
			onlineUserIDs[uid] = true
		}

		// 遍历所有成员，推送给每个人（包括发送者自己）
		for _, member := range members {
			userID := member.UserID

			// 检查用户是否在线
			isOnline := onlineUserIDs[userID]
			if isOnline {
				onlineCount++
			} else {
				offlineCount++
			}

			// 发送消息给每个成员（无论在线与否都尝试发送）
			// SendToUser会检查用户是否在线，如果不在线会记录日志但不阻塞
			// 消息已经保存到数据库，离线用户上线后可以通过历史消息获取
			common.SendToUser(userID, messageBytes)

			// 记录日志（只记录前10个成员，避免日志过多）
			if successCount < 10 {
				status := "离线"
				if isOnline {
					status = "在线"
				}
				fmt.Printf("  - 推送消息给成员 %d (用户ID: %d, 状态: %s)\n",
					member.ID, userID, status)
			}
			successCount++
		}

		fmt.Printf("群聊消息推送完成: 群组ID=%d, 总成员数=%d, 在线=%d, 离线=%d, 消息ID=%s\n",
			message.GroupID, len(members), onlineCount, offlineCount, message.MessageID)

		// 如果所有成员都离线，记录警告
		if onlineCount == 0 && len(members) > 0 {
			fmt.Printf("警告: 群组 %d 的所有 %d 个成员都不在线，消息已保存但无法实时推送\n",
				message.GroupID, len(members))
		}
	} else {
		// 私聊：同时推送给接收者和发送者
		fmt.Printf("准备通过WebSocket推送私聊消息给接收者 %d 和发送者 %d，消息长度: %d\n",
			message.ToUserID, message.FromUserID, len(messageBytes))

		// 推送给接收者
		common.SendToUser(message.ToUserID, messageBytes)
		fmt.Printf("WebSocket消息已发送给接收者 %d\n", message.ToUserID)

		// 推送给发送者（确保发送者也能实时看到自己发送的消息）
		common.SendToUser(message.FromUserID, messageBytes)
		fmt.Printf("WebSocket消息已发送给发送者 %d\n", message.FromUserID)
	}
}
