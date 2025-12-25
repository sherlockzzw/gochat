package chat

import (
	"gochat/internal/component"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/pkg/response"
	"gochat/utils"
)

type ChatHandler struct {
	response     *response.SvcRequest
	dao          *dao.ChatDao
	friendDao    *dao.FriendDao
	groupDao     *dao.GroupDao
	muteDao      *dao.GroupMuteDao
	userDao      *dao.UserDao
	favoriteDao  *dao.MessageFavoriteDao
}

func NewChatHandler(server *component.ApiServer) *ChatHandler {
	return &ChatHandler{
		response:    server.Result,
		dao:         dao.NewChatDao(utils.DB, utils.MongoDB),
		friendDao:   dao.NewFriendDao(utils.DB),
		groupDao:    dao.NewGroupDao(utils.DB),
		muteDao:     dao.NewGroupMuteDao(utils.DB),
		userDao:     dao.NewUserDao(utils.DB),
		favoriteDao: dao.NewMessageFavoriteDao(utils.DB),
	}
}
