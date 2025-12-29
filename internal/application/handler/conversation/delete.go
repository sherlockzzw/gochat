package conversation

import (
	"gochat/api/api/conversation"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// DeleteConversation 删除会话
func (h *ConversationHandler) DeleteConversation(ctx *gin.Context) {
	req, err := analysis.BindParameter[conversation.DeleteConversationRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.deleteConversationLogic(ctx, &req)
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

func (h *ConversationHandler) deleteConversationLogic(ctx *gin.Context, req *conversation.DeleteConversationRequest) (resp *conversation.DeleteConversationResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	otherUserID := int64(req.GetOtherUserId())
	groupID := int64(req.GetGroupId())

	// 验证参数：私聊或群聊二选一
	if (otherUserID == 0 && groupID == 0) || (otherUserID > 0 && groupID > 0) {
		return nil, code_msg.BadRequest, nil
	}

	// 删除会话
	err = h.dao.DeleteConversation(userID, otherUserID, groupID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	resp = &conversation.DeleteConversationResponse{
		Success: true,
	}

	return resp, 0, nil
}


