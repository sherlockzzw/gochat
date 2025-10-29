package friend

import (
	"gochat/internal/component"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/pkg/response"
	"gochat/utils"
)

type FriendHandler struct {
	response *response.SvcRequest
	dao      *dao.FriendDao
}

func NewFriendHandler(server *component.ApiServer) *FriendHandler {
	return &FriendHandler{
		response: server.Result,
		dao:      dao.NewFriendDao(utils.DB),
	}
}
