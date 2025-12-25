package handler

import (
	"gochat/internal/application/handler/balance"
	"gochat/internal/application/handler/chat"
	"gochat/internal/application/handler/friend"
	"gochat/internal/application/handler/group"
	"gochat/internal/application/handler/user"
	"gochat/internal/component"
)

type API struct {
	UserHandler    *user.UserHandler
	ChatHandler    *chat.ChatHandler
	FriendHandler  *friend.FriendHandler
	GroupHandler   *group.GroupHandler
	BalanceHandler *balance.BalanceHandler
}

// NewApi 用户端接口注册
func NewApi() *API {
	server := component.GetApiServer()
	return &API{
		UserHandler:    user.NewUserHandler(server),
		ChatHandler:    chat.NewChatHandler(server),
		FriendHandler:  friend.NewFriendHandler(server),
		GroupHandler:   group.NewGroupHandler(server),
		BalanceHandler: balance.NewBalanceHandler(server),
	}
}
