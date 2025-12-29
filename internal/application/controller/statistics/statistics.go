package statistics

import (
	"gochat/internal/infrastructure/dao"
	"gochat/internal/pkg/response"
	"gochat/utils"
)

type StatisticsController struct {
	response *response.SvcRequest
	balanceDao *dao.BalanceDao
}

func NewStatisticsController() *StatisticsController {
	return &StatisticsController{
		response:   response.NewSvcRequest(),
		balanceDao: dao.NewBalanceDao(utils.DB),
	}
}

