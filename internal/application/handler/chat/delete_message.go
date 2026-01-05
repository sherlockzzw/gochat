package chat

import (
	"gochat/api/api/chat"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// DeleteMessage 删除消息
func (h *ChatHandler) DeleteMessage(ctx *gin.Context) {
	req, err := analysis.BindParameter[chat.DeleteMessageRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.deleteMessageLogic(ctx, &req)
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

func (h *ChatHandler) deleteMessageLogic(ctx *gin.Context, req *chat.DeleteMessageRequest) (resp *chat.DeleteMessageResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 删除消息
	err = h.dao.DeleteMessage(req.GetMessageId(), userID, req.GetDeleteForBoth())
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return &chat.DeleteMessageResponse{
		Success: true,
	}, 0, nil
}




