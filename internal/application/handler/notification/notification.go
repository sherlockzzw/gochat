package notification

import (
	"gochat/internal/component"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/pkg/response"
	"gochat/utils"
)

type NotificationHandler struct {
	response *response.SvcRequest
	dao      *dao.NotificationDao
}

func NewNotificationHandler(server *component.ApiServer) *NotificationHandler {
	return &NotificationHandler{
		response: server.Result,
		dao:      dao.NewNotificationDao(utils.DB),
	}
}

