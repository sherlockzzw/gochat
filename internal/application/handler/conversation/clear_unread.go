package conversation

import (
	"gochat/api/api/conversation"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// ClearUnreadCount 清除未读数
func (h *ConversationHandler) ClearUnreadCount(ctx *gin.Context) {
	req, err := analysis.BindParameter[conversation.ClearUnreadCountRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.clearUnreadCountLogic(ctx, &req)
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

func (h *ConversationHandler) clearUnreadCountLogic(ctx *gin.Context, req *conversation.ClearUnreadCountRequest) (resp *conversation.ClearUnreadCountResponse, errCode code_msg.BusinessCode, err error) {
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

	// 清除未读数
	clearedCount, err := h.dao.ClearUnreadCount(userID, otherUserID, groupID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	resp = &conversation.ClearUnreadCountResponse{
		Success:      true,
		ClearedCount: int32(clearedCount),
	}

	return resp, 0, nil
}


