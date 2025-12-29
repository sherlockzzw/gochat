package chat

import (
	"fmt"
	"gochat/api/api/chat"
	"gochat/internal/infrastructure/models"
	globalUtils "gochat/utils"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// convertToChatMessage 将模型消息转换为API响应消息
func (h *ChatHandler) convertToChatMessage(message *models.ChatMessage) *chat.ChatMessage {
	return &chat.ChatMessage{
		Id:             message.MessageID,
		FromUserId:     uint32(message.FromUserID),
		ToUserId:       uint32(message.ToUserID),
		GroupId:        uint32(message.GroupID),
		MessageType:    chat.MessageType(message.MessageType),
		Content:        message.Content,
		FileUrl:        message.FileURL,
		FileName:       message.FileName,
		FileSize:       message.FileSize,
		VideoUrl:       message.VideoURL,
		VideoThumb:     message.VideoThumb,
		VoiceUrl:       message.VoiceURL,
		VoiceDuration:  int32(message.VoiceDuration),
		EmojiUrl:       message.EmojiURL,
		MergeMessages:  message.MergeMessages,
		QuoteMessageId: message.QuoteMessageID,
		ContactUserId:  uint32(message.ContactUserID),
		RedPacketId:    uint32(message.RedPacketID),
		TransferId:     uint32(message.TransferID),
		IsRecalled:     message.IsRecalled,
		IsDeleted:      message.IsDeleted,
		ReadStatus:     int32(message.ReadStatus),
		Status:         chat.MessageStatus(message.Status),
		CreatedAt:      timestamppb.New(time.Unix(message.CreatedAt, 0)),
		UpdatedAt:      timestamppb.New(time.Unix(message.UpdatedAt, 0)),
	}
}

// convertToChatMessageWithDetails 将模型消息转换为API响应消息（包含详细信息）
func (h *ChatHandler) convertToChatMessageWithDetails(message *models.ChatMessage) *chat.ChatMessage {
	chatMsg := h.convertToChatMessage(message)

	// 如果是引用消息，获取被引用消息的详细信息
	if message.MessageType == models.MessageTypeQuote && message.QuoteMessageID != "" {
		quotedMsg, err := h.getQuotedMessageInfo(message.QuoteMessageID)
		if err == nil && quotedMsg != nil {
			quotedChatMsg := h.convertToChatMessage(quotedMsg)
			chatMsg.QuotedMessage = quotedChatMsg
		}
	}

	// 如果是联系人分享消息，获取联系人详细信息
	if message.MessageType == models.MessageTypeContact && message.ContactUserID > 0 {
		contactUser, err := h.getContactUserInfo(message.ContactUserID)
		if err == nil && contactUser != nil {
			apiPort := 8080
			if port := globalUtils.GetApiPort(); port > 0 {
				apiPort = port
			}
			chatMsg.ContactUser = &chat.UserInfo{
				Id:     uint32(contactUser.ID),
				Name:   contactUser.Name,
				Phone:  contactUser.Phone,
				Email:  contactUser.Email,
				Avatar: globalUtils.GetAvatarFullURL(contactUser.Avatar, fmt.Sprintf("http://127.0.0.1:%d", apiPort)),
			}
		}
	}

	return chatMsg
}


