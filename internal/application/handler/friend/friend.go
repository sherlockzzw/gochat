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
	userDao  *dao.UserDao
}

func NewFriendHandler(server *component.ApiServer) *FriendHandler {
	return &FriendHandler{
		response: server.Result,
		dao:      dao.NewFriendDao(utils.DB),
		userDao:  dao.NewUserDao(utils.DB),
	}
}
