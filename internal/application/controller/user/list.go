package user

import (
	admin_user "gochat/api/admin/user"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/pkg/analysis"
	"gochat/utils"

	"github.com/gin-gonic/gin"
)

func (h *ControllerUser) GetUserList(ctx *gin.Context) {
	req, err := analysis.BindQuery[admin_user.GetUserListRequest](ctx, h.response)
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

func (h *ControllerUser) getUserListLogic(ctx *gin.Context, req admin_user.GetUserListRequest) (resp *admin_user.GetUserListResponse, err error) {
	page := int(req.GetPage())
	if page <= 0 {
		page = 1
	}
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 20
	}

	users, _, err := h.dao.GetUserListWithPagination(page, pageSize)
	if err != nil {
		return nil, gin.Error{Err: gin.Error{}, Type: gin.ErrorTypePublic, Meta: "获取用户列表失败"}
	}

	// 转换为protobuf结构体
	var userInfos []*admin_user.UserInfo
	balanceDao := dao.NewBalanceDao(utils.DB)
	for _, u := range users {
		// 获取用户余额
		balance, _ := balanceDao.GetBalance(int64(u.ID))
		balanceAmount := int64(0)
		if balance != nil {
			balanceAmount = balance.Balance
		}

		userInfos = append(userInfos, &admin_user.UserInfo{
			Id:         int64(u.ID),
			Name:       u.Name,
			Phone:      u.Phone,
			Email:      u.Email,
			ClientIp:   u.ClientIp,
			ClientPort: u.ClientPort,
			DeviceInfo: u.DeviceInfo,
			Status:     1, // TODO: 从用户表获取状态
			Balance:    balanceAmount,
		})
	}

	resp = &admin_user.GetUserListResponse{
		Data: userInfos,
	}

	return resp, nil
}
