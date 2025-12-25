package user

import (
	"gochat/api/api/user"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/verification"
	"gochat/internal/infrastructure/models"
	"gochat/utils"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (h *UserHandler) Register(ctx *gin.Context) {
	req, err := analysis.BindParameter[user.RegisterRequest](ctx, h.response)
	if err != nil {
		return
	}

	resp, code, err := h.registerLogic(ctx, &req)
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

func (h *UserHandler) registerLogic(ctx *gin.Context, req *user.RegisterRequest) (resp *user.UserInfo, errCode code_msg.BusinessCode, err error) {
	// 1. 检查是否同意协议
	if !req.GetAgreeTerms() || !req.GetAgreePrivacy() {
		return nil, code_msg.AgreementNotAccepted, nil
	}

	// 2. 检查登录账号是否已存在
	existingUser, err := h.dao.GetUserByLoginAccount(req.GetLoginAccount())
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if existingUser != nil {
		return nil, code_msg.LoginAccountExists, nil
	}

	// 3. 验证手机号或邮箱
	var (
		target   string
		codeType string
	)
	if req.GetPhone() != "" {
		target = req.GetPhone()
		codeType = verification.CodeTypePhone
		// 检查手机号是否已被使用
		phoneUser, err := h.dao.GetUserByPhone(req.GetPhone())
		if err != nil {
			return nil, code_msg.ServerError, err
		}
		if phoneUser != nil {
			return nil, code_msg.PhoneExists, nil
		}
	} else if req.GetEmail() != "" {
		target = req.GetEmail()
		codeType = verification.CodeTypeEmail
		// 检查邮箱是否已被使用
		emailUser, err := h.dao.GetUserByEmail(req.GetEmail())
		if err != nil {
			return nil, code_msg.ServerError, err
		}
		if emailUser != nil {
			return nil, code_msg.EmailExists, nil
		}
	} else {
		return nil, code_msg.PhoneOrEmailRequired, nil
	}

	// 4. 验证验证码类型是否匹配
	if req.GetCodeType() != codeType {
		return nil, code_msg.CodeTypeMismatch, nil
	}

	// 5. 验证验证码
	valid, err := h.verification.VerifyCode(ctx.Request.Context(), codeType, target, req.GetCode())
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	if !valid {
		return nil, code_msg.CodeError, nil
	}

	// 6. 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), 12)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 7. 创建用户
	now := time.Now().Unix()
	userModel := &models.UserBasic{
		LoginAccount:   req.GetLoginAccount(),
		Name:           req.GetName(),
		Password:       string(hashedPassword),
		Phone:          req.GetPhone(),
		PhoneVerified:  req.GetPhone() != "",
		Email:          req.GetEmail(),
		ClientIp:       req.GetClientIp(),
		ClientPort:     req.GetClientPort(),
		DeviceInfo:     req.GetDeviceInfo(),
		LoginFailCount: 0,
		LockedUntil:    0,
		LastActiveTime: now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err = h.dao.CreateUser(userModel)
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 8. 创建设备记录(如果提供了设备信息)
	if req.GetDeviceId() != "" && req.GetDeviceType() != "" {
		deviceDao := dao.NewDeviceDao(utils.DB)
		device := &models.UserDevice{
			UserID:       userModel.ID,
			DeviceType:   req.GetDeviceType(),
			DeviceID:     req.GetDeviceId(),
			DeviceName:   req.GetDeviceName(),
			LastActiveAt: now,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		_ = deviceDao.CreateDevice(device) // 忽略错误，不影响注册
	}

	// 9. 返回用户信息
	userInfo := &user.UserInfo{
		Id:             int64(userModel.ID),
		Name:           userModel.Name,
		LoginAccount:   userModel.LoginAccount,
		Phone:          userModel.Phone,
		PhoneVerified:  userModel.PhoneVerified,
		Email:          userModel.Email,
		ClientIp:       userModel.ClientIp,
		ClientPort:     userModel.ClientPort,
		DeviceInfo:     userModel.DeviceInfo,
		LoginFailCount: int32(userModel.LoginFailCount),
		LockedUntil:    userModel.LockedUntil,
		LastActiveTime: userModel.LastActiveTime,
		CreatedAt:      time.Unix(userModel.CreatedAt, 0).Format("2006-01-02 15:04:05"),
		UpdatedAt:      time.Unix(userModel.UpdatedAt, 0).Format("2006-01-02 15:04:05"),
	}

	return userInfo, 0, nil
}
