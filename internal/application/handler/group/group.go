package group

import (
	"gochat/internal/component"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/pkg/response"
	"gochat/utils"
)

type GroupHandler struct {
	response *response.SvcRequest
	dao      *dao.GroupDao
}

func NewGroupHandler(server *component.ApiServer) *GroupHandler {
	return &GroupHandler{
		response: server.Result,
		dao:      dao.NewGroupDao(utils.DB),
	}
}

