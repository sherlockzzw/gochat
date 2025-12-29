package content

import (
	"gochat/internal/infrastructure/dao"
	"gochat/internal/pkg/response"
	"gochat/utils"
)

type ContentController struct {
	response          *response.SvcRequest
	violationDao      *dao.ContentViolationDao
	adminLogDao       *dao.AdminLogDao
	userDao           *dao.UserDao
}

func NewContentController() *ContentController {
	return &ContentController{
		response:     response.NewSvcRequest(),
		violationDao: dao.NewContentViolationDao(utils.DB),
		adminLogDao:  dao.NewAdminLogDao(utils.DB),
		userDao:      dao.NewUserDao(utils.DB),
	}
}

