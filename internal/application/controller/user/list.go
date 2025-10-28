package user

import (
	"gochat/api/admin/user"
	"gochat/internal/pkg/analysis"

	"github.com/gin-gonic/gin"
)

func (h *ControllerUser) GetUserList(ctx *gin.Context) {
	req, err := analysis.BindQuery[user.GetUserListRequest](ctx, h.response)
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

func (h *ControllerUser) getUserListLogic(ctx *gin.Context, req user.GetUserListRequest) (resp *user.GetUserListResponse, err error) {
	users, err := h.dao.GetUserList()
	if err != nil {
		return nil, gin.Error{Err: gin.Error{}, Type: gin.ErrorTypePublic, Meta: "获取用户列表失败"}
	}

	// 转换为protobuf结构体
	var userInfos []*user.UserInfo
	for _, u := range users {
		userInfos = append(userInfos, &user.UserInfo{
			Id:         int64(u.ID),
			Name:       u.Name,
			Phone:      u.Phone,
			Email:      u.Email,
			ClientIp:   u.ClientIp,
			ClientPort: u.ClientPort,
			DeviceInfo: u.DeviceInfo,
		})
	}

	resp = &user.GetUserListResponse{
		Data: userInfos,
	}

	return resp, nil
}
