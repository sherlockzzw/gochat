package user

import (
	"fmt"
	"gochat/api/api/user"
	"gochat/internal/pkg/analysis"
	"gochat/internal/pkg/code_msg"
	"gochat/internal/infrastructure/models"
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
	// 根据登录账号/手机号/邮箱获取用户
	var userModel *models.UserBasic
	
	loginAccount := req.GetLoginAccount()
	// 尝试按登录账号查询
	userModel, err = h.dao.GetUserByLoginAccount(loginAccount)
	if err != nil {
		return nil, code_msg.ServerError, err
	}
	// 如果登录账号不存在，尝试按手机号查询
	if userModel == nil {
		userModel, err = h.dao.GetUserByPhone(loginAccount)
		if err != nil {
			return nil, code_msg.ServerError, err
		}
	}
	// 如果手机号也不存在，尝试按邮箱查询
	if userModel == nil {
		userModel, err = h.dao.GetUserByEmail(loginAccount)
		if err != nil {
			return nil, code_msg.ServerError, err
		}
	}
	
	if userModel == nil {
		return nil, code_msg.UserNotExists, nil
	}

	// 检查账户是否被锁定
	now := time.Now().Unix()
	if userModel.LockedUntil > 0 && userModel.LockedUntil > now {
		// 账户被锁定，计算剩余锁定时间（秒）
		remainingSeconds := userModel.LockedUntil - now
		remainingMinutes := remainingSeconds / 60
		return nil, code_msg.AccountLocked, fmt.Errorf("账户已被锁定，请在 %d 分钟后重试", remainingMinutes+1)
	}

	// 如果锁定时间已过，重置锁定状态
	if userModel.LockedUntil > 0 && userModel.LockedUntil <= now {
		err = h.dao.UpdateUser(userModel.ID, map[string]interface{}{
			"login_fail_count": 0,
			"locked_until":      0,
		})
		if err != nil {
			return nil, code_msg.ServerError, err
		}
		userModel.LoginFailCount = 0
		userModel.LockedUntil = 0
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(req.GetPassword()))
	if err != nil {
		// 密码错误，增加失败计数
		failCount := userModel.LoginFailCount + 1
		updates := map[string]interface{}{
			"login_fail_count": failCount,
		}

		// 如果失败次数达到5次，锁定账户1小时
		if failCount >= 5 {
			lockUntil := now + 3600 // 锁定1小时（3600秒）
			updates["locked_until"] = lockUntil
		}

		err = h.dao.UpdateUser(userModel.ID, updates)
		if err != nil {
			return nil, code_msg.ServerError, err
		}

		// 如果达到锁定条件，返回锁定错误
		if failCount >= 5 {
			return nil, code_msg.AccountLocked, fmt.Errorf("登录失败次数过多，账户已被锁定1小时")
		}

		return nil, code_msg.PasswordError, nil
	}

	// 登录成功，重置失败计数和锁定状态
	err = h.dao.UpdateUser(userModel.ID, map[string]interface{}{
		"login_fail_count": 0,
		"locked_until":      0,
		"login_time":       now,
		"last_active_time": now,
	})
	if err != nil {
		return nil, code_msg.ServerError, err
	}

	// 生成token（根据自动登录天数设置有效期）
	authMiddleware := middleware.JwtMiddleware("UserBasic")
	
	// 获取自动登录天数（默认7天，范围7-30天）
	autoLoginDays := int(req.GetAutoLoginDays())
	if autoLoginDays == 0 {
		autoLoginDays = 7 // 默认7天
	}
	if autoLoginDays < 7 {
		autoLoginDays = 7
	}
	if autoLoginDays > 30 {
		autoLoginDays = 30
	}

	expireTime := time.Duration(autoLoginDays) * 24 * time.Hour
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
