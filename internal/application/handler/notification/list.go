package notification

import (
	"gochat/api/api/notification"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetNotifications 获取通知列表
func (h *NotificationHandler) GetNotifications(ctx *gin.Context) {
	req, err := analysis.BindQuery[notification.GetNotificationsRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.getNotificationsLogic(ctx, &req)
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

func (h *NotificationHandler) getNotificationsLogic(ctx *gin.Context, req *notification.GetNotificationsRequest) (resp *notification.GetNotificationsResponse, errCode code_msg.BusinessCode, err error) {
	// 从JWT中获取当前用户ID
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	page := req.GetPage()
	pageSize := req.GetPageSize()
	notificationType := req.GetType()
	isReadStr := req.GetIsRead()

	// 处理isRead参数（从字符串转换为bool指针）
	var isReadPtr *bool
	if isReadStr != "" {
		if isReadStr == "true" {
			read := true
			isReadPtr = &read
		} else if isReadStr == "false" {
			read := false
			isReadPtr = &read
		}
	}

	// 获取通知列表
	notifications, total, err := h.dao.GetNotifications(userID, page, pageSize, notificationType, isReadPtr)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 获取未读数量
	unreadCount, err := h.dao.GetUnreadCount(userID, notificationType)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 转换为响应格式
	notificationInfos := make([]*notification.NotificationInfo, 0, len(notifications))
	for _, n := range notifications {
		info := &notification.NotificationInfo{
			Id:        uint64(n.ID),
			Type:      n.Type,
			Title:     n.Title,
			Content:   n.Content,
			Amount:    n.Amount,
			RelatedId: uint64(n.RelatedID),
			IsRead:    n.IsRead,
			CreatedAt: timestamppb.New(time.Unix(n.CreatedAt, 0)),
		}
		notificationInfos = append(notificationInfos, info)
	}

	resp = &notification.GetNotificationsResponse{
		Notifications: notificationInfos,
		TotalCount:    int32(total),
		CurrentPage:   page,
		PageSize:      pageSize,
		UnreadCount:   int32(unreadCount),
	}

	return resp, 0, nil
}

