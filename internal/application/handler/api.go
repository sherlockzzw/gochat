package handler

import (
	"gochat/internal/application/handler/chat"
	"gochat/internal/application/handler/user"
	"gochat/internal/component"
)

type API struct {
	UserHandler *user.UserHandler
	ChatHandler *chat.ChatHandler
}

// NewApi 用户端接口注册
func NewApi() *API {
	server := component.GetApiServer()
	return &API{
		UserHandler: user.NewUserHandler(server),
		ChatHandler: chat.NewChatHandler(server),
	}
}
