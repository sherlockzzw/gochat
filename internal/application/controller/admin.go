package controller

import (
	"gochat/internal/application/controller/auth"
	"gochat/internal/application/controller/config"
	"gochat/internal/application/controller/content"
	"gochat/internal/application/controller/finance"
	"gochat/internal/application/controller/log"
	"gochat/internal/application/controller/permission"
	"gochat/internal/application/controller/statistics"
	"gochat/internal/application/controller/user"
	"gochat/internal/component"
)

type API struct {
	ControllerUser       *user.ControllerUser
	AuthController       *auth.AuthController
	PermissionController *permission.PermissionController
	ConfigController     *config.ConfigController
	FinanceController    *finance.FinanceController
	ContentController    *content.ContentController
	StatisticsController *statistics.StatisticsController
	LogController        *log.LogController
}

// NewAdmin 管理端接口注册
func NewAdmin() *API {
	return &API{
		ControllerUser:       user.NewControllerUser(component.GetApiServer()),
		AuthController:       auth.NewAuthController(),
		PermissionController: permission.NewPermissionController(),
		ConfigController:     config.NewConfigController(),
		FinanceController:    finance.NewFinanceController(),
		ContentController:    content.NewContentController(),
		StatisticsController: statistics.NewStatisticsController(),
		LogController:        log.NewLogController(),
	}
}
