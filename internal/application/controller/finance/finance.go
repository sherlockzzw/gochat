package finance

import (
	"gochat/internal/infrastructure/dao"
	"gochat/internal/pkg/response"
	"gochat/utils"
)

type FinanceController struct {
	response    *response.SvcRequest
	balanceDao  *dao.BalanceDao
	configDao   *dao.SystemConfigDao
	adminLogDao *dao.AdminLogDao
	userDao     *dao.UserDao
}

func NewFinanceController() *FinanceController {
	return &FinanceController{
		response:    response.NewSvcRequest(),
		balanceDao:  dao.NewBalanceDao(utils.DB),
		configDao:   dao.NewSystemConfigDao(utils.DB),
		adminLogDao: dao.NewAdminLogDao(utils.DB),
		userDao:     dao.NewUserDao(utils.DB),
	}
}

