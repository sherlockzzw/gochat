package balance

import (
	"gochat/internal/component"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/pkg/response"
	"gochat/utils"
)

type BalanceHandler struct {
	response *response.SvcRequest
	dao      *dao.BalanceDao
}

func NewBalanceHandler(server *component.ApiServer) *BalanceHandler {
	return &BalanceHandler{
		response: server.Result,
		dao:      dao.NewBalanceDao(utils.DB),
	}
}


