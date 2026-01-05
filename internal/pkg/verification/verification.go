package verification

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	// 验证码类型
	CodeTypePhone = "phone"
	CodeTypeEmail = "email"

	// Redis Key前缀
	CodeKeyPrefix = "verification:code:"
	CodeResendKeyPrefix = "verification:resend:"

	// 验证码有效期(5分钟)
	CodeExpireTime = 5 * time.Minute
	// 重发间隔(60秒)
	ResendInterval = 60 * time.Second
	// 验证码长度
	CodeLength = 6
)

type VerificationService struct {
	rdb *redis.Client
}

func NewVerificationService(rdb *redis.Client) *VerificationService {
	return &VerificationService{rdb: rdb}
}

// GenerateCode 生成验证码
func (s *VerificationService) GenerateCode() string {
	rand.Seed(time.Now().UnixNano())
	code := ""
	for i := 0; i < CodeLength; i++ {
		code += fmt.Sprintf("%d", rand.Intn(10))
	}
	return code
}

// SendCode 发送验证码(存储到Redis)
func (s *VerificationService) SendCode(ctx context.Context, codeType, target string) (string, error) {
	// 检查重发间隔
	resendKey := CodeResendKeyPrefix + codeType + ":" + target
	exists, err := s.rdb.Exists(ctx, resendKey).Result()
	if err != nil {
		return "", err
	}
	if exists > 0 {
		ttl, _ := s.rdb.TTL(ctx, resendKey).Result()
		return "", fmt.Errorf("请等待%d秒后重试", int(ttl.Seconds())+1)
	}

	// 生成验证码
	code := s.GenerateCode()

	// 存储验证码(5分钟有效期)
	codeKey := CodeKeyPrefix + codeType + ":" + target
	err = s.rdb.Set(ctx, codeKey, code, CodeExpireTime).Err()
	if err != nil {
		return "", err
	}

	// 设置重发间隔(60秒)
	err = s.rdb.Set(ctx, resendKey, "1", ResendInterval).Err()
	if err != nil {
		return "", err
	}

	// TODO: 实际发送验证码(短信/邮件)
	// 这里只返回验证码，实际项目中需要调用短信/邮件服务
	if codeType == CodeTypePhone {
		// 发送短信验证码
		// smsService.Send(target, code)
	} else if codeType == CodeTypeEmail {
		// 发送邮件验证码
		// emailService.Send(target, code)
	}

	return code, nil
}

// VerifyCode 验证验证码
func (s *VerificationService) VerifyCode(ctx context.Context, codeType, target, code string) (bool, error) {
	codeKey := CodeKeyPrefix + codeType + ":" + target
	storedCode, err := s.rdb.Get(ctx, codeKey).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil // 验证码不存在或已过期
		}
		return false, err
	}

	if storedCode != code {
		return false, nil // 验证码错误
	}

	// 验证成功后删除验证码(防止重复使用)
	s.rdb.Del(ctx, codeKey)

	return true, nil
}

// DeleteCode 删除验证码(验证成功后调用)
func (s *VerificationService) DeleteCode(ctx context.Context, codeType, target string) error {
	codeKey := CodeKeyPrefix + codeType + ":" + target
	return s.rdb.Del(ctx, codeKey).Err()
}





