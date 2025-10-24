package user

import (
	"gochat/internal/pkg/analysis"
	"gochat/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	Name       string `json:"name" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	ClientIp   string `json:"client_ip"`
	ClientPort string `json:"client_port"`
	DeviceInfo string `json:"device_info"`
}

type CreateUserResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (h *ControllerUser) CreateUser(ctx *gin.Context) {
	req, err := analysis.BindParameter[CreateUserRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, err := h.createUserLogic(ctx, req)
	if err != nil {
		h.response.JsonError(ctx, err, err.Error())
		return
	}

	h.response.JsonSuccess(ctx, resp)
}

func (h *ControllerUser) createUserLogic(ctx *gin.Context, req *CreateUserRequest) (resp *CreateUserResponse, err error) {
	// 检查用户是否已存在
	existingUser, err := h.dao.GetUserByName(req.Name)
	if err == nil && existingUser != nil {
		return nil, gin.Error{Err: gin.Error{}, Type: gin.ErrorTypePublic, Meta: "用户已经存在"}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, gin.Error{Err: gin.Error{}, Type: gin.ErrorTypePublic, Meta: "密码加密失败"}
	}

	// 创建用户
	user := &models.UserBasic{
		Name:       req.Name,
		Password:   string(hashedPassword),
		Phone:      req.Phone,
		Email:      req.Email,
		ClientIp:   req.ClientIp,
		ClientPort: req.ClientPort,
		DeviceInfo: req.DeviceInfo,
	}

	err = h.dao.CreateUser(user)
	if err != nil {
		return nil, gin.Error{Err: gin.Error{}, Type: gin.ErrorTypePublic, Meta: "创建失败"}
	}

	resp = &CreateUserResponse{
		Code:    200,
		Message: "创建成功",
		Data:    "",
	}

	return resp, nil
}
