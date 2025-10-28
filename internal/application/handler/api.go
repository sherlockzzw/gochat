package handler

import (
	"gochat/internal/application/handler/user"
	"gochat/internal/component"
)

type API struct {
	UserHandler *user.UserHandler
}

// NewApi 用户端接口注册
func NewApi() *API {
	return &API{
		UserHandler: user.NewUserHandler(component.GetApiServer()),
	}
}
