package chat

import (
	"encoding/json"
	"fmt"
	"gochat/api/api/chat"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
	globalUtils "gochat/utils"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ForwardMessage 转发消息
func (h *ChatHandler) ForwardMessage(ctx *gin.Context) {
	req, err := analysis.BindParameter[chat.ForwardMessageRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.forwardMessageLogic(ctx, &req)
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

func (h *ChatHandler) forwardMessageLogic(ctx *gin.Context, req *chat.ForwardMessageRequest) (resp *chat.ForwardMessageResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	fromUserID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 验证参数：私聊或群聊二选一
	groupID := int64(req.GetGroupId())
	toUserID := int64(req.GetToUserId())

	if groupID == 0 && toUserID == 0 {
		return nil, code_msg.BadRequest, nil
	}

	if groupID > 0 && toUserID > 0 {
		return nil, code_msg.BadRequest, nil
	}

	// 验证消息数量限制（1-100条）
	messageIDs := req.GetMessageIds()
	if len(messageIDs) == 0 || len(messageIDs) > 100 {
		return nil, code_msg.BadRequest, nil
	}

	// 如果是群聊，验证用户是否在群组中
	if groupID > 0 {
		groupDao := dao.NewGroupDao(globalUtils.DB)
		isMember, err := groupDao.IsMemberInGroup(groupID, fromUserID)
		if err != nil {
			return nil, code_msg.ServerError, err
		}
		if !isMember {
			return nil, code_msg.NotGroupMember, nil
		}
	}

	// 如果是私聊，验证是否为好友关系
	if toUserID > 0 {
		friendDao := dao.NewFriendDao(globalUtils.DB)
		friendRelation, err := friendDao.CheckIsFriend(fromUserID, toUserID)
		if err != nil {
			return nil, code_msg.ServerError, err
		}
		if friendRelation == nil {
			return nil, code_msg.NotFriend, nil
		}
	}

	// 获取要转发的原始消息
	originalMessages, err := h.dao.GetMessagesByIDs(messageIDs)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	if len(originalMessages) == 0 {
		return nil, code_msg.BadRequest, nil
	}

	// 验证消息是否已撤回或已删除
	for _, msg := range originalMessages {
		if msg.IsRecalled || msg.IsDeleted {
			return nil, code_msg.BadRequest, nil
		}
	}

	// 判断转发方式：单条转发或多条合并转发
	var forwardedMessages []*models.ChatMessage
	now := time.Now().Unix()

	if len(originalMessages) == 1 {
		// 单条转发：直接复制消息内容
		originalMsg := originalMessages[0]
		forwardedMsg := &models.ChatMessage{
			FromUserID:     fromUserID,
			ToUserID:       toUserID,
			GroupID:        groupID,
			MessageType:    originalMsg.MessageType,
			Content:        req.GetContent(), // 如果有附加内容，使用附加内容；否则使用原消息内容
			FileURL:        originalMsg.FileURL,
			FileName:       originalMsg.FileName,
			FileSize:       originalMsg.FileSize,
			VideoURL:       originalMsg.VideoURL,
			VideoThumb:     originalMsg.VideoThumb,
			VoiceURL:       originalMsg.VoiceURL,
			VoiceDuration:  originalMsg.VoiceDuration,
			EmojiURL:       originalMsg.EmojiURL,
			QuoteMessageID: originalMsg.MessageID, // 引用原消息ID
			ContactUserID:   originalMsg.ContactUserID,
			// 注意：红包和转账不能转发，需要清空
			RedPacketID:    0,
			TransferID:     0,
			ReadStatus:     1,
			Status:         models.MessageStatusSent,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		// 如果没有附加内容，使用原消息内容
		if forwardedMsg.Content == "" {
			forwardedMsg.Content = originalMsg.Content
		}

		forwardedMessages = append(forwardedMessages, forwardedMsg)
	} else {
		// 多条转发：使用合并消息类型
		// 收集所有要转发的消息ID
		mergeMessageIDs := make([]string, 0, len(originalMessages))
		for _, msg := range originalMessages {
			mergeMessageIDs = append(mergeMessageIDs, msg.MessageID)
		}

		mergeMessagesJSON, _ := json.Marshal(mergeMessageIDs)

		// 构建合并消息内容（显示转发了几条消息）
		mergeContent := fmt.Sprintf("转发了 %d 条消息", len(originalMessages))
		if req.GetContent() != "" {
			mergeContent = req.GetContent() + "\n" + mergeContent
		}

		forwardedMsg := &models.ChatMessage{
			FromUserID:    fromUserID,
			ToUserID:      toUserID,
			GroupID:       groupID,
			MessageType:   models.MessageTypeMerge,
			Content:       mergeContent,
			MergeMessages: string(mergeMessagesJSON),
			ReadStatus:    1,
			Status:        models.MessageStatusSent,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		forwardedMessages = append(forwardedMessages, forwardedMsg)
	}

	// 保存转发的消息
	var savedMessages []*models.ChatMessage
	for _, msg := range forwardedMessages {
		err = h.dao.CreateMessage(msg)
		if err != nil {
			// 如果部分消息保存失败，记录错误但继续处理其他消息
			fmt.Printf("Failed to forward message: %v\n", err)
			continue
		}
		savedMessages = append(savedMessages, msg)
	}

	if len(savedMessages) == 0 {
		return nil, code_msg.ServerError, fmt.Errorf("failed to forward any message")
	}

	// 通过WebSocket实时推送消息
	for _, msg := range savedMessages {
		h.sendMessageViaWebSocket(msg)
	}

	// 转换为响应格式
	chatMessages := make([]*chat.ChatMessage, 0, len(savedMessages))
	for _, msg := range savedMessages {
		chatMsg := &chat.ChatMessage{
			Id:             msg.MessageID,
			FromUserId:     uint32(msg.FromUserID),
			ToUserId:       uint32(msg.ToUserID),
			GroupId:        uint32(msg.GroupID),
			MessageType:    chat.MessageType(msg.MessageType),
			Content:        msg.Content,
			FileUrl:        msg.FileURL,
			FileName:       msg.FileName,
			FileSize:       msg.FileSize,
			VideoUrl:      msg.VideoURL,
			VideoThumb:    msg.VideoThumb,
			VoiceUrl:      msg.VoiceURL,
			VoiceDuration: int32(msg.VoiceDuration),
			EmojiUrl:      msg.EmojiURL,
			MergeMessages: msg.MergeMessages,
			QuoteMessageId: msg.QuoteMessageID,
			ContactUserId:  uint32(msg.ContactUserID),
			RedPacketId:    uint32(msg.RedPacketID),
			TransferId:     uint32(msg.TransferID),
			IsRecalled:    msg.IsRecalled,
			IsDeleted:     msg.IsDeleted,
			ReadStatus:    int32(msg.ReadStatus),
			Status:        chat.MessageStatus(msg.Status),
			CreatedAt:     timestamppb.New(time.Unix(msg.CreatedAt, 0)),
			UpdatedAt:     timestamppb.New(time.Unix(msg.UpdatedAt, 0)),
		}
		chatMessages = append(chatMessages, chatMsg)
	}

	resp = &chat.ForwardMessageResponse{
		Messages:        chatMessages,
		ForwardedCount:  int32(len(savedMessages)),
		Success:         true,
	}

	return resp, 0, nil
}

