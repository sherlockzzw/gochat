package user

import (
	admin_user "gochat/api/admin/user"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/pkg/analysis"
	"gochat/utils"

	"github.com/gin-gonic/gin"
)

// GetUserInfo 获取用户信息
func (h *ControllerUser) GetUserInfo(ctx *gin.Context) {
	req, err := analysis.BindQuery[admin_user.GetUserInfoRequest](ctx, h.response)
	if err != nil {
		return
	}

	u, err := h.dao.GetUserByID(req.GetId())
	if err != nil {
		h.response.JsonError(ctx, err, err.Error())
		return
	}
	if u == nil {
		h.response.JsonErrorFixation(ctx, 404)
		return
	}

	// 获取用户余额
	balanceDao := dao.NewBalanceDao(utils.DB)
	balance, _ := balanceDao.GetBalance(req.GetId())
	balanceAmount := int64(0)
	if balance != nil {
		balanceAmount = balance.Balance
	}

	userInfo := &admin_user.UserInfo{
		Id:         int64(u.ID),
		Name:       u.Name,
		Phone:      u.Phone,
		Email:      u.Email,
		ClientIp:   u.ClientIp,
		ClientPort: u.ClientPort,
		DeviceInfo: u.DeviceInfo,
		Status:     1, // TODO: 从用户表获取状态
		Balance:    balanceAmount,
	}

	resp := &admin_user.GetUserInfoResponse{
		Message: "获取成功",
		User:    userInfo,
	}

	h.response.JsonSuccess(ctx, resp)
}

