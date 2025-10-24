package user

import (
	"gochat/internal/pkg/analysis"
	"gochat/models"

	"github.com/gin-gonic/gin"
)

type GetUserListRequest struct {
	Page     int32 `form:"page"`
	PageSize int32 `form:"page_size"`
}

type GetUserListResponse struct {
	Data []*models.UserBasic `json:"data"`
}

func (h *ControllerUser) GetUserList(ctx *gin.Context) {
	req, err := analysis.BindQuery[GetUserListRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, err := h.getUserListLogic(ctx, req)
	if err != nil {
		h.response.JsonError(ctx, err, err.Error())
		return
	}

	h.response.JsonSuccess(ctx, resp)
}

func (h *ControllerUser) getUserListLogic(ctx *gin.Context, req *GetUserListRequest) (resp *GetUserListResponse, err error) {
	users, err := h.dao.GetUserList()
	if err != nil {
		return nil, gin.Error{Err: gin.Error{}, Type: gin.ErrorTypePublic, Meta: "获取用户列表失败"}
	}

	resp = &GetUserListResponse{
		Data: users,
	}

	return resp, nil
}
