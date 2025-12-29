package auth

import (
	"gochat/internal/infrastructure/dao"
	"gochat/internal/pkg/response"
	"gochat/utils"
)

type AuthController struct {
	response *response.SvcRequest
	adminDao *dao.AdminDao
}

func NewAuthController() *AuthController {
	return &AuthController{
		response: response.NewSvcRequest(),
		adminDao: dao.NewAdminDao(utils.DB),
	}
}

