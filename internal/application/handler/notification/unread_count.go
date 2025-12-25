package notification

import (
	"gochat/api/api/notification"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetUnreadCount 获取未读通知数量
func (h *NotificationHandler) GetUnreadCount(ctx *gin.Context) {
	req, err := analysis.BindQuery[notification.GetUnreadCountRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.getUnreadCountLogic(ctx, &req)
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

func (h *NotificationHandler) getUnreadCountLogic(ctx *gin.Context, req *notification.GetUnreadCountRequest) (resp *notification.GetUnreadCountResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 获取未读数量
	count, err := h.dao.GetUnreadCount(userID, "")
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	resp = &notification.GetUnreadCountResponse{
		Count: int32(count),
	}

	return resp, 0, nil
}

