package permission

import (
	"gochat/internal/infrastructure/dao"
	"gochat/internal/pkg/response"
	"gochat/utils"
)

type PermissionController struct {
	response      *response.SvcRequest
	permissionDao *dao.PermissionDao
	adminLogDao   *dao.AdminLogDao
}

func NewPermissionController() *PermissionController {
	return &PermissionController{
		response:      response.NewSvcRequest(),
		permissionDao: dao.NewPermissionDao(utils.DB),
		adminLogDao:   dao.NewAdminLogDao(utils.DB),
	}
}

