package log

import (
	"gochat/internal/infrastructure/dao"
	"gochat/internal/pkg/response"
	"gochat/utils"
)

type LogController struct {
	response   *response.SvcRequest
	adminLogDao *dao.AdminLogDao
}

func NewLogController() *LogController {
	return &LogController{
		response:    response.NewSvcRequest(),
		adminLogDao: dao.NewAdminLogDao(utils.DB),
	}
}

