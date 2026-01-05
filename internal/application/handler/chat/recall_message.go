package chat

import (
	"gochat/api/api/chat"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// RecallMessage 撤回消息
func (h *ChatHandler) RecallMessage(ctx *gin.Context) {
	req, err := analysis.BindParameter[chat.RecallMessageRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.recallMessageLogic(ctx, &req)
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

func (h *ChatHandler) recallMessageLogic(ctx *gin.Context, req *chat.RecallMessageRequest) (resp *chat.RecallMessageResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 撤回消息
	err = h.dao.RecallMessage(req.GetMessageId(), userID)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	return &chat.RecallMessageResponse{
		Success: true,
	}, 0, nil
}




