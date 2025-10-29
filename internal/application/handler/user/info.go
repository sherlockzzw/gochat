package user

import (
	"fmt"
	"gochat/api/api/user"
	"gochat/internal/pkg/analysis"
	globalUtils "gochat/utils"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) GetUserInfo(ctx *gin.Context) {
	req, err := analysis.BindQuery[user.GetUserInfoRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, err := h.getUserInfoLogic(ctx, req)
	if err != nil {
		h.response.JsonError(ctx, err, err.Error())
		return
	}

	h.response.JsonSuccess(ctx, resp)
}

func (h *UserHandler) getUserInfoLogic(ctx *gin.Context, req user.GetUserInfoRequest) (resp *user.GetUserInfoResponse, err error) {
	// 获取用户信息
	userModel, err := h.dao.GetUserByID(uint(req.GetId()))
	if err != nil {
		h.response.JsonNotFound(ctx, "用户不存在")
		return nil, gin.Error{Err: gin.Error{}, Type: gin.ErrorTypePublic, Meta: "用户不存在"}
	}

	// 获取API端口配置
	apiPort := 8080 // 默认端口
	if port := globalUtils.GetApiPort(); port > 0 {
		apiPort = port
	}

	resp = &user.GetUserInfoResponse{
		Message: "获取成功",
		User: &user.UserInfo{
			Id:         int64(userModel.ID),
			Name:       userModel.Name,
			Phone:      userModel.Phone,
			Email:      userModel.Email,
			Avatar:     globalUtils.GetAvatarFullURL(userModel.Avatar, fmt.Sprintf("http://127.0.0.1:%d", apiPort)), // 返回完整头像URL
			ClientIp:   userModel.ClientIp,
			ClientPort: userModel.ClientPort,
			DeviceInfo: userModel.DeviceInfo,
		},
	}

	return resp, nil
}
