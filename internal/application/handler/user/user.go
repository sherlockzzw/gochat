package user

import (
	"gochat/internal/component"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/pkg/response"
	"gochat/internal/pkg/verification"
	"gochat/utils"
)

type UserHandler struct {
	response      *response.SvcRequest
	dao           *dao.UserDao
	verification  *verification.VerificationService
}

func NewUserHandler(server *component.ApiServer) *UserHandler {
	return &UserHandler{
		response:     server.Result,
		dao:          dao.NewUserDao(utils.DB),
		verification: verification.NewVerificationService(utils.RDB),
	}
}
