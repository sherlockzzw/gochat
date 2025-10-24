package controller

import (
	"gochat/internal/application/controller/user"
	"gochat/internal/component"
)

type API struct {
	ControllerUser *user.ControllerUser
}

// NewAdmin 管理端接口注册
func NewAdmin() *API {
	return &API{
		ControllerUser: user.NewControllerUser(component.GetApiServer()),
	}
}
