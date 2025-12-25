package chat

import (
	"fmt"
	"gochat/api/api/chat"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/code_msg"
)

// validateMessageType 验证消息类型相关的业务逻辑
func (h *ChatHandler) validateMessageType(req *chat.SendMessageRequest, fromUserID int64) (code_msg.BusinessCode, error) {
	messageType := req.GetMessageType()

	// 验证引用消息
	if messageType == chat.MessageType_QUOTE {
		if req.GetQuoteMessageId() == "" {
			return code_msg.BadRequest, fmt.Errorf("引用消息ID不能为空")
		}

		// 验证被引用的消息是否存在
		messages, err := h.dao.GetMessagesByIDs([]string{req.GetQuoteMessageId()})
		if err != nil {
			return code_msg.ServerError, err
		}
		if len(messages) == 0 {
			return code_msg.BadRequest, fmt.Errorf("被引用的消息不存在")
		}

		quotedMsg := messages[0]
		// 验证消息是否已撤回或已删除
		if quotedMsg.IsRecalled {
			return code_msg.BadRequest, fmt.Errorf("被引用的消息已撤回")
		}
		if quotedMsg.IsDeleted {
			return code_msg.BadRequest, fmt.Errorf("被引用的消息已删除")
		}

		// 验证消息是否在当前会话中（私聊或群聊）
		// 如果是私聊，验证消息是否在双方会话中
		if req.GetGroupId() == 0 {
			toUserID := int64(req.GetToUserId())
			if (quotedMsg.FromUserID != fromUserID && quotedMsg.FromUserID != toUserID) ||
				(quotedMsg.ToUserID != fromUserID && quotedMsg.ToUserID != toUserID) {
				return code_msg.BadRequest, fmt.Errorf("被引用的消息不在当前会话中")
			}
		} else {
			// 如果是群聊，验证消息是否在该群组中
			groupID := int64(req.GetGroupId())
			if quotedMsg.GroupID != groupID {
				return code_msg.BadRequest, fmt.Errorf("被引用的消息不在当前群组中")
			}
		}
	}

	// 验证联系人分享消息
	if messageType == chat.MessageType_CONTACT {
		if req.GetContactUserId() == 0 {
			return code_msg.BadRequest, fmt.Errorf("联系人ID不能为空")
		}

		contactUserID := int64(req.GetContactUserId())
		// 验证联系人是否存在
		contactUser, err := h.userDao.GetUserByID(contactUserID)
		if err != nil {
			return code_msg.ServerError, err
		}
		if contactUser == nil {
			return code_msg.BadRequest, fmt.Errorf("联系人不存在")
		}

		// 验证是否为好友关系（只有好友才能分享）
		friendRelation, err := h.friendDao.CheckIsFriend(fromUserID, contactUserID)
		if err != nil {
			return code_msg.ServerError, err
		}
		if friendRelation == nil {
			return code_msg.NotFriend, fmt.Errorf("只能分享好友联系人")
		}
	}

	// 验证合并消息
	if messageType == chat.MessageType_MERGE {
		mergeMessageIDs := req.GetMergeMessageIds()
		if len(mergeMessageIDs) == 0 {
			return code_msg.BadRequest, fmt.Errorf("合并消息ID列表不能为空")
		}

		// 验证合并消息数量限制
		imageCount := 0
		videoCount := 0
		fileCount := 0
		otherCount := 0

		// 获取所有被合并的消息
		messages, err := h.dao.GetMessagesByIDs(mergeMessageIDs)
		if err != nil {
			return code_msg.ServerError, err
		}
		if len(messages) != len(mergeMessageIDs) {
			return code_msg.BadRequest, fmt.Errorf("部分被合并的消息不存在")
		}

		// 统计各类型消息数量并验证
		for _, msg := range messages {
			// 验证消息是否已撤回或已删除
			if msg.IsRecalled || msg.IsDeleted {
				return code_msg.BadRequest, fmt.Errorf("被合并的消息中包含已撤回或已删除的消息")
			}

			// 验证消息是否在当前会话中
			if req.GetGroupId() == 0 {
				toUserID := int64(req.GetToUserId())
				if (msg.FromUserID != fromUserID && msg.FromUserID != toUserID) ||
					(msg.ToUserID != fromUserID && msg.ToUserID != toUserID) {
					return code_msg.BadRequest, fmt.Errorf("被合并的消息不在当前会话中")
				}
			} else {
				groupID := int64(req.GetGroupId())
				if msg.GroupID != groupID {
					return code_msg.BadRequest, fmt.Errorf("被合并的消息不在当前群组中")
				}
			}

			// 统计消息类型
			switch msg.MessageType {
			case models.MessageTypeImage:
				imageCount++
			case models.MessageTypeVideo:
				videoCount++
			case models.MessageTypeFile:
				fileCount++
			default:
				otherCount++
			}
		}

		// 验证数量限制：10图+5视频+3文件
		if imageCount > 10 {
			return code_msg.BadRequest, fmt.Errorf("合并消息中图片数量不能超过10张")
		}
		if videoCount > 5 {
			return code_msg.BadRequest, fmt.Errorf("合并消息中视频数量不能超过5个")
		}
		if fileCount > 3 {
			return code_msg.BadRequest, fmt.Errorf("合并消息中文件数量不能超过3个")
		}

		// 验证是否包含不支持合并的消息类型（红包、转账等不能合并）
		for _, msg := range messages {
			if msg.MessageType == models.MessageTypeRedPacket || msg.MessageType == models.MessageTypeTransfer {
				return code_msg.BadRequest, fmt.Errorf("红包和转账消息不能合并")
			}
		}
	}

	return 0, nil
}

// getQuotedMessageInfo 获取被引用消息的详细信息（用于响应）
func (h *ChatHandler) getQuotedMessageInfo(quoteMessageID string) (*models.ChatMessage, error) {
	if quoteMessageID == "" {
		return nil, nil
	}

	messages, err := h.dao.GetMessagesByIDs([]string{quoteMessageID})
	if err != nil {
		return nil, err
	}
	if len(messages) == 0 {
		return nil, nil
	}

	return messages[0], nil
}

// getContactUserInfo 获取联系人详细信息（用于响应）
func (h *ChatHandler) getContactUserInfo(contactUserID int64) (*models.UserBasic, error) {
	if contactUserID == 0 {
		return nil, nil
	}

	return h.userDao.GetUserByID(contactUserID)
}

