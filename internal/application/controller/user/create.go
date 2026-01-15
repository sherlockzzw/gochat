package user

import (
	"gochat/api/admin/user"
	"gochat/internal/infrastructure/models"
	"gochat/internal/pkg/analysis"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser 创建用户
func (h *ControllerUser) CreateUser(ctx *gin.Context) {
	req, err := analysis.BindParameter[user.CreateUserRequest](ctx, h.response)
	if err != nil {
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.response.JsonError(ctx, err, err.Error())
		return
	}

	now := time.Now().Unix()
	userBasic := &models.UserBasic{
		Name:            req.Name,
		Password:        string(hashedPassword),
		Phone:           req.Phone,
		Email:           req.Email,
		ClientIp:        req.ClientIp,
		PrivacySettings: "[]",
		CreatedAt:       now,
		UpdatedAt:       now,
		LoginAccount:    req.Name,
	}

	if err := h.dao.CreateUser(userBasic); err != nil {
		h.response.JsonError(ctx, err, err.Error())
		return
	}

	resp := &user.CreateUserResponse{
		Code:    0,
		Message: "创建成功",
		Data:    "success",
	}

	h.response.JsonSuccess(ctx, resp)
}
