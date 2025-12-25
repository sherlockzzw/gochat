package conversation

import (
	"gochat/api/api/conversation"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetConversations 获取会话列表
func (h *ConversationHandler) GetConversations(ctx *gin.Context) {
	req, err := analysis.BindQuery[conversation.GetConversationsRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.getConversationsLogic(ctx, &req)
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

func (h *ConversationHandler) getConversationsLogic(ctx *gin.Context, req *conversation.GetConversationsRequest) (resp *conversation.GetConversationsResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	page := req.GetPage()
	pageSize := req.GetPageSize()
	includeHidden := req.GetIncludeHidden()

	// 获取会话列表（包含设置）
	conversations, total, err := h.dao.GetConversationsWithSettings(userID, page, pageSize, includeHidden)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 转换为响应格式
	conversationInfos := make([]*conversation.ConversationInfo, 0, len(conversations))
	for _, convWithSetting := range conversations {
		conv := convWithSetting.Conversation
		setting := convWithSetting.Setting

		info := &conversation.ConversationInfo{
			OtherUserId:    uint32(conv.OtherUserID),
			GroupId:        uint32(0), // 当前Conversation模型只支持私聊
			LastMessage:    conv.LastMessage,
			LastMessageType: int32(conv.LastMessageType),
			UnreadCount:     int32(conv.UnreadCount),
			LastMessageAt:   timestamppb.New(time.Unix(conv.LastMessageAt, 0)),
			CreatedAt:       timestamppb.New(time.Unix(conv.CreatedAt, 0)),
			UpdatedAt:       timestamppb.New(time.Unix(conv.UpdatedAt, 0)),
		}

		// 添加设置信息
		if setting != nil {
			info.IsPinned = setting.IsPinned
			info.IsMuted = setting.IsMuted
			info.IsHidden = setting.IsHidden
			// 使用设置中的未读数（更准确）
			info.UnreadCount = int32(setting.UnreadCount)
		}

		conversationInfos = append(conversationInfos, info)
	}

	resp = &conversation.GetConversationsResponse{
		Conversations: conversationInfos,
		TotalCount:    int32(total),
		CurrentPage:   page,
		PageSize:      pageSize,
	}

	return resp, 0, nil
}

