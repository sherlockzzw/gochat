package handler

import (
	"gochat/internal/application/controller/user"
	"gochat/internal/component"
)

type API struct {
	UserHandler *user.ControllerUser
}

// NewApi 用户端接口注册
func NewApi() *API {
	return &API{
		UserHandler: user.NewControllerUser(component.GetApiServer()),
	}
}
