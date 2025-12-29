package config

import (
	"gochat/internal/infrastructure/dao"
	"gochat/internal/pkg/response"
	"gochat/utils"
)

type ConfigController struct {
	response     *response.SvcRequest
	configDao    *dao.SystemConfigDao
	adminLogDao  *dao.AdminLogDao
}

func NewConfigController() *ConfigController {
	return &ConfigController{
		response:    response.NewSvcRequest(),
		configDao:   dao.NewSystemConfigDao(utils.DB),
		adminLogDao: dao.NewAdminLogDao(utils.DB),
	}
}

