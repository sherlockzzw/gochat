package notification

import (
	"gochat/api/api/notification"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// DeleteNotification 删除通知
func (h *NotificationHandler) DeleteNotification(ctx *gin.Context) {
	req, err := analysis.BindParameter[notification.DeleteNotificationRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.deleteNotificationLogic(ctx, &req)
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

func (h *NotificationHandler) deleteNotificationLogic(ctx *gin.Context, req *notification.DeleteNotificationRequest) (resp *notification.DeleteNotificationResponse, errCode code_msg.BusinessCode, err error) {
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

	// 删除通知
	deletedCount, err := h.dao.DeleteNotifications(userID, ids)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	resp = &notification.DeleteNotificationResponse{
		Success:      true,
		DeletedCount: int32(deletedCount),
	}

	return resp, 0, nil
}


