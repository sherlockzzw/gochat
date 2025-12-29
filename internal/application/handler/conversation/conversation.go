package conversation

import (
	"gochat/internal/component"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/pkg/response"
	"gochat/utils"
)

type ConversationHandler struct {
	response *response.SvcRequest
	dao      *dao.ConversationDao
}

func NewConversationHandler(server *component.ApiServer) *ConversationHandler {
	return &ConversationHandler{
		response: server.Result,
		dao:      dao.NewConversationDao(utils.DB),
	}
}


