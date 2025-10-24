package user

import (
	"gochat/internal/component"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/pkg/response"

	"go.uber.org/zap"
)

type ControllerUser struct {
	apiServer *component.ApiServer
	response  *response.SvcRequest
	dao       *dao.UserDao
	logger    *zap.Logger
}

func NewControllerUser(apiServer *component.ApiServer) *ControllerUser {
	return &ControllerUser{
		apiServer: apiServer,
		response:  apiServer.Result,
		dao:       dao.NewUserDao(apiServer.MysqlSvc.DB),
		logger:    apiServer.Logger,
	}
}
