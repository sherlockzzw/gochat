package user

import (
	"gochat/api/api/user"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (h *UserHandler) Register(ctx *gin.Context) {
	req, err := analysis.BindParameter[user.RegisterRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.registerLogic(ctx, req)
	if code != 0 {
		h.response.JsonErrorFixation(ctx, code)
		return
	}
	if err != nil {
		h.response.JsonError(ctx, err, err.Error())
		return
	}

	h.response.JsonSuccess(ctx, resp)
}

func (h *UserHandler) registerLogic(ctx *gin.Context, req user.RegisterRequest) (resp *user.UserInfo, errCode code_msg.BusinessCode, err error) {
	// 检查用户是否已存在
	existingUser, err := h.dao.GetUserByName(req.GetName())
	if err == nil && existingUser != nil {
		return nil, code_msg.UserNameExisted, nil
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), 12)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 创建用户
	userModel := &models.UserBasic{
		Name:       req.GetName(),
		Password:   string(hashedPassword),
		Phone:      req.GetPhone(),
		Email:      req.GetEmail(),
		ClientIp:   req.GetClientIp(),
		ClientPort: req.GetClientPort(),
		DeviceInfo: req.GetDeviceInfo(),
	}

	err = h.dao.CreateUser(userModel)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 返回用户信息
	resp = &user.UserInfo{
		Id:         int64(userModel.ID),
		Name:       userModel.Name,
		Phone:      userModel.Phone,
		Email:      userModel.Email,
		ClientIp:   userModel.ClientIp,
		ClientPort: userModel.ClientPort,
		DeviceInfo: userModel.DeviceInfo,
	}

	return resp, 0, nil
}
