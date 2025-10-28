package user

import (
	"gochat/api/api/user"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (h *UserHandler) Login(ctx *gin.Context) {
	req, err := analysis.BindParameter[user.LoginRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.loginLogic(ctx, req)
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

func (h *UserHandler) loginLogic(ctx *gin.Context, req user.LoginRequest) (resp *user.LoginResponse, errCode code_msg.BusinessCode, err error) {
	// 获取用户
	userModel, err := h.dao.GetUserByName(req.GetName())
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if userModel == nil {
		return nil, code_msg.UserNotExists, nil
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(req.GetPassword()))
	if err != nil {
		return nil, code_msg.PasswordError, nil
	}

	// 生成token
	authMiddleware := middleware.JwtMiddleware("UserBasic")
	expireTime := time.Duration(24) * time.Hour
	token, expire, err := authMiddleware.GenerateToken(userModel, expireTime)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	resp = &user.LoginResponse{
		Data: &user.UserInfo{
			Id:         int64(userModel.ID),
			Name:       userModel.Name,
			Phone:      userModel.Phone,
			Email:      userModel.Email,
			ClientIp:   userModel.ClientIp,
			ClientPort: userModel.ClientPort,
			DeviceInfo: userModel.DeviceInfo,
		},
		Token:  token,
		Expire: expire,
	}

	return resp, 0, nil
}
