package conversation

import (
	"gochat/api/api/conversation"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// HideConversation 设置会话隐藏
func (h *ConversationHandler) HideConversation(ctx *gin.Context) {
	req, err := analysis.BindParameter[conversation.HideConversationRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.hideConversationLogic(ctx, &req)
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

func (h *ConversationHandler) hideConversationLogic(ctx *gin.Context, req *conversation.HideConversationRequest) (resp *conversation.HideConversationResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	otherUserID := int64(req.GetOtherUserId())
	groupID := int64(req.GetGroupId())
	isHidden := req.GetIsHidden()

	// 验证参数：私聊或群聊二选一
	if (otherUserID == 0 && groupID == 0) || (otherUserID > 0 && groupID > 0) {
		return nil, code_msg.BadRequest, nil
	}

	// 更新会话设置
	err = h.dao.UpdateConversationSetting(userID, otherUserID, groupID, map[string]interface{}{
		"is_hidden": isHidden,
	})
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	resp = &conversation.HideConversationResponse{
		Success: true,
	}

	return resp, 0, nil
}


