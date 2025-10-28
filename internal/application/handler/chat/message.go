package chat

import (
	"encoding/json"
	"fmt"
	"gochat/api/api/chat"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/models"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	// 从JWT中获取当前用户ID（这里需要从中间件中获取）
	// TODO: 从JWT token中获取用户ID
	fromUserID := uint(1) // 临时硬编码，实际应该从JWT中获取

	// 创建消息模型
	message := &models.ChatMessage{
		FromUserID:  fromUserID,
		ToUserID:    uint(req.GetToUserId()),
		MessageType: int(req.GetMessageType()),
		Content:     req.GetContent(),
		FileURL:     req.GetFileUrl(),
		FileName:    req.GetFileName(),
		FileSize:    req.GetFileSize(),
		Status:      models.MessageStatusSent,
	}

	// 保存消息
	err = h.dao.CreateMessage(message)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 通过WebSocket实时推送消息
	h.sendMessageViaWebSocket(message)

	// 转换为响应格式
	chatMessage := &chat.ChatMessage{
		Id:          message.MessageID,
		FromUserId:  uint32(message.FromUserID),
		ToUserId:    uint32(message.ToUserID),
		MessageType: chat.MessageType(message.MessageType),
		Content:     message.Content,
		FileUrl:     message.FileURL,
		FileName:    message.FileName,
		FileSize:    message.FileSize,
		Status:      chat.MessageStatus(message.Status),
		CreatedAt:   timestamppb.New(message.CreatedAt),
		UpdatedAt:   timestamppb.New(message.UpdatedAt),
	}

	resp = &chat.SendMessageResponse{
		Message:      chatMessage,
		Success:      true,
		ErrorMessage: "",
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

	h.response.JsonSuccess(ctx, resp)
}

func (h *ChatHandler) getMessageHistoryLogic(ctx *gin.Context, req chat.GetMessageHistoryRequest) (resp *chat.GetMessageHistoryResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	// TODO: 从JWT token中获取用户ID
	fromUserID := uint(1) // 临时硬编码

	// 解析时间参数
	var beforeTime *time.Time
	if req.GetBeforeTime() != nil {
		t := req.GetBeforeTime().AsTime()
		beforeTime = &t
	}

	// 获取消息历史
	messages, totalCount, err := h.dao.GetMessageHistory(
		fromUserID,
		uint(req.GetOtherUserId()),
		int(req.GetPage()),
		int(req.GetPageSize()),
		beforeTime,
	)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 转换为响应格式
	var chatMessages []*chat.ChatMessage
	for _, msg := range messages {
		chatMsg := &chat.ChatMessage{
			Id:          msg.MessageID,
			FromUserId:  uint32(msg.FromUserID),
			ToUserId:    uint32(msg.ToUserID),
			MessageType: chat.MessageType(msg.MessageType),
			Content:     msg.Content,
			FileUrl:     msg.FileURL,
			FileName:    msg.FileName,
			FileSize:    msg.FileSize,
			Status:      chat.MessageStatus(msg.Status),
			CreatedAt:   timestamppb.New(msg.CreatedAt),
			UpdatedAt:   timestamppb.New(msg.UpdatedAt),
		}
		chatMessages = append(chatMessages, chatMsg)
	}

	hasMore := int64(req.GetPage()*req.GetPageSize()) < totalCount

	resp = &chat.GetMessageHistoryResponse{
		Messages:    chatMessages,
		TotalCount:  int32(totalCount),
		CurrentPage: req.GetPage(),
		PageSize:    req.GetPageSize(),
		HasMore:     hasMore,
	}

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
	// TODO: 从JWT token中获取用户ID
	toUserID := uint(1) // 临时硬编码

	// 标记消息为已读
	markedCount, err := h.dao.MarkMessagesAsRead(
		uint(req.GetOtherUserId()),
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
	// TODO: 从JWT token中获取用户ID
	userID := uint(1) // 临时硬编码

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
			"message_type": message.MessageType,
			"content":      message.Content,
			"file_url":     message.FileURL,
			"file_name":    message.FileName,
			"file_size":    message.FileSize,
			"status":       message.Status,
			"created_at":   message.CreatedAt.Unix(),
		},
	}

	// 将消息转换为JSON
	messageBytes, err := json.Marshal(wsMessage)
	if err != nil {
		fmt.Printf("Failed to marshal WebSocket message: %v\n", err)
		return
	}

	// 发送给接收者
	// TODO: 这里需要获取WebSocket Hub实例
	// 暂时通过Redis Pub/Sub实现
	h.sendMessageViaRedis(message.ToUserID, messageBytes)
}

// sendMessageViaRedis 通过Redis Pub/Sub发送消息
func (h *ChatHandler) sendMessageViaRedis(userID uint, message []byte) {
	// TODO: 实现Redis Pub/Sub消息推送
	fmt.Printf("Sending message to user %d via Redis: %s\n", userID, string(message))
}
