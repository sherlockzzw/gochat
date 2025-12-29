package notification

import (
	"gochat/api/api/notification"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// MarkAsRead 标记通知为已读
func (h *NotificationHandler) MarkAsRead(ctx *gin.Context) {
	req, err := analysis.BindParameter[notification.MarkAsReadRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.markAsReadLogic(ctx, &req)
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

func (h *NotificationHandler) markAsReadLogic(ctx *gin.Context, req *notification.MarkAsReadRequest) (resp *notification.MarkAsReadResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	notificationIDs := req.GetNotificationIds()
	if len(notificationIDs) == 0 {
		return nil, code_msg.BadRequest, nil
	}

	// 转换为int64
	ids := make([]int64, 0, len(notificationIDs))
	for _, id := range notificationIDs {
		ids = append(ids, int64(id))
	}

	// 标记为已读
	updatedCount, err := h.dao.MarkAsRead(userID, ids)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	resp = &notification.MarkAsReadResponse{
		Success:      true,
		UpdatedCount: int32(updatedCount),
	}

	return resp, 0, nil
}

// MarkAllAsRead 标记所有通知为已读
func (h *NotificationHandler) MarkAllAsRead(ctx *gin.Context) {
	req, err := analysis.BindParameter[notification.MarkAllAsReadRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.markAllAsReadLogic(ctx, &req)
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

func (h *NotificationHandler) markAllAsReadLogic(ctx *gin.Context, req *notification.MarkAllAsReadRequest) (resp *notification.MarkAllAsReadResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 标记所有通知为已读
	updatedCount, err := h.dao.MarkAllAsRead(userID, "")
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	resp = &notification.MarkAllAsReadResponse{
		Success:      true,
		UpdatedCount: int32(updatedCount),
	}

	return resp, 0, nil
}


