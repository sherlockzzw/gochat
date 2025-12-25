package chat

import (
	"gochat/api/api/chat"
	"gochat/internal/application/handler/common"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// FavoriteMessage 收藏消息
func (h *ChatHandler) FavoriteMessage(ctx *gin.Context) {
	req, err := analysis.BindParameter[chat.FavoriteMessageRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.favoriteMessageLogic(ctx, &req)
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

func (h *ChatHandler) favoriteMessageLogic(ctx *gin.Context, req *chat.FavoriteMessageRequest) (resp *chat.FavoriteMessageResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, errCode, err := common.GetUserIDFromContext(ctx)
	if errCode != 0 {
		return nil, errCode, err
	}

	messageID := req.GetMessageId()
	if messageID == "" {
		return nil, code_msg.BadRequest, nil
	}

	// 验证消息是否存在
	messages, err := h.dao.GetMessagesByIDs([]string{messageID})
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if len(messages) == 0 {
		return nil, code_msg.BadRequest, nil
	}

	// 验证消息是否已撤回或已删除
	msg := messages[0]
	if msg.IsRecalled || msg.IsDeleted {
		return nil, code_msg.BadRequest, nil
	}

	// 检查是否已经收藏
	existing, err := h.favoriteDao.GetFavorite(userID, messageID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if existing != nil {
		// 已经收藏，直接返回成功
		return &chat.FavoriteMessageResponse{
			Success: true,
		}, 0, nil
	}

	// 创建收藏
	err = h.favoriteDao.CreateFavorite(userID, messageID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return &chat.FavoriteMessageResponse{
		Success: true,
	}, 0, nil
}

// UnfavoriteMessage 取消收藏消息
func (h *ChatHandler) UnfavoriteMessage(ctx *gin.Context) {
	req, err := analysis.BindParameter[chat.UnfavoriteMessageRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.unfavoriteMessageLogic(ctx, &req)
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

func (h *ChatHandler) unfavoriteMessageLogic(ctx *gin.Context, req *chat.UnfavoriteMessageRequest) (resp *chat.UnfavoriteMessageResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, errCode, err := common.GetUserIDFromContext(ctx)
	if errCode != 0 {
		return nil, errCode, err
	}

	messageID := req.GetMessageId()
	if messageID == "" {
		return nil, code_msg.BadRequest, nil
	}

	// 取消收藏
	err = h.favoriteDao.DeleteFavorite(userID, messageID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return &chat.UnfavoriteMessageResponse{
		Success: true,
	}, 0, nil
}

// GetFavoriteMessages 获取收藏消息列表
func (h *ChatHandler) GetFavoriteMessages(ctx *gin.Context) {
	req, err := analysis.BindQuery[chat.GetFavoriteMessagesRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.getFavoriteMessagesLogic(ctx, &req)
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

func (h *ChatHandler) getFavoriteMessagesLogic(ctx *gin.Context, req *chat.GetFavoriteMessagesRequest) (resp *chat.GetFavoriteMessagesResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, errCode, err := common.GetUserIDFromContext(ctx)
	if errCode != 0 {
		return nil, errCode, err
	}

	page := req.GetPage()
	pageSize := req.GetPageSize()

	// 获取收藏的消息ID列表
	messageIDs, total, err := h.favoriteDao.GetFavoriteMessageIDs(userID, page, pageSize)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 获取消息详情
	var messages []*models.ChatMessage
	if len(messageIDs) > 0 {
		messages, err = h.dao.GetMessagesByIDs(messageIDs)
		if err != nil {
			return nil, code_msg.ServerError, err
		}
	}

	// 转换为响应格式
	chatMessages := make([]*chat.ChatMessage, 0, len(messages))
	for _, msg := range messages {
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

	resp = &chat.GetFavoriteMessagesResponse{
		Messages:    chatMessages,
		TotalCount:  int32(total),
		CurrentPage: page,
		PageSize:    pageSize,
	}

	return resp, 0, nil
}

