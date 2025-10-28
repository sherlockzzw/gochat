package chat

import (
	"gochat/internal/component"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/pkg/response"
	"gochat/utils"
)

type ChatHandler struct {
	response *response.SvcRequest
	dao      *dao.ChatDao
}

func NewChatHandler(server *component.ApiServer) *ChatHandler {
	return &ChatHandler{
		response: server.Result,
		dao:      dao.NewChatDao(utils.DB, utils.MongoDB),
	}
}
