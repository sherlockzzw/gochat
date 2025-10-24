package user

import (
	"gochat/internal/pkg/analysis"
	"gochat/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Token   string      `json:"token"`
	Expire  int64       `json:"expire"`
}

func (h *ControllerUser) Login(ctx *gin.Context) {
	req, err := analysis.BindParameter[LoginRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, err := h.loginLogic(ctx, req)
	if err != nil {
		h.response.JsonError(ctx, err, err.Error())
		return
	}

	h.response.JsonSuccess(ctx, resp)
}

func (h *ControllerUser) loginLogic(ctx *gin.Context, req *LoginRequest) (resp *LoginResponse, err error) {
	// 获取用户
	user, err := h.dao.GetUserByName(req.Name)
	if err != nil {
		return nil, err
	}
	if user == nil {
		h.response.JsonNotFound(ctx, "用户不存在")
		return nil, gin.Error{Err: gin.Error{}, Type: gin.ErrorTypePublic, Meta: "用户不存在"}
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		h.response.JsonUnauthorized(ctx, "密码错误")
		return nil, gin.Error{Err: gin.Error{}, Type: gin.ErrorTypePublic, Meta: "密码错误"}
	}

	// 生成token
	authMiddleware := middleware.JwtMiddleware("UserBasic")
	expireTime := time.Duration(24) * time.Hour
	token, expire, err := authMiddleware.GenerateToken(user, expireTime)
	if err != nil {
		h.response.JsonError(ctx, err, "Token生成失败")
		return nil, gin.Error{Err: gin.Error{}, Type: gin.ErrorTypePublic, Meta: "Token生成失败"}
	}

	resp = &LoginResponse{
		Message: "登录成功",
		Data:    user,
		Token:   token,
		Expire:  expire,
	}

	return resp, nil
}
