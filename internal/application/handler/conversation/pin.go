package conversation

import (
	"gochat/api/api/conversation"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// PinConversation 设置会话置顶
func (h *ConversationHandler) PinConversation(ctx *gin.Context) {
	req, err := analysis.BindParameter[conversation.PinConversationRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.pinConversationLogic(ctx, &req)
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

func (h *ConversationHandler) pinConversationLogic(ctx *gin.Context, req *conversation.PinConversationRequest) (resp *conversation.PinConversationResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	otherUserID := int64(req.GetOtherUserId())
	groupID := int64(req.GetGroupId())
	isPinned := req.GetIsPinned()

	// 验证参数：私聊或群聊二选一
	if (otherUserID == 0 && groupID == 0) || (otherUserID > 0 && groupID > 0) {
		return nil, code_msg.BadRequest, nil
	}

	// 更新会话设置
	err = h.dao.UpdateConversationSetting(userID, otherUserID, groupID, map[string]interface{}{
		"is_pinned": isPinned,
	})
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	resp = &conversation.PinConversationResponse{
		Success: true,
	}

	return resp, 0, nil
}

