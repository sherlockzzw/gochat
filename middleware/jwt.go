package middleware

import (
	"context"
	jwt "github.com/appleboy/gin-jwt/v2"
	"gochat/models"
	"gochat/utils"
	"log"
	"strconv"
	"time"
)

var identityKey = "id"
var encryptionKey = "examplekey946666" // 16字节密钥

type JwtMiddlewareWrapper struct {
	*jwt.GinJWTMiddleware
}

func JwtMiddleware(modelType string) *JwtMiddlewareWrapper {
	var timeout time.Duration
	var payloadFunc func(data interface{}) jwt.MapClaims

	switch modelType {
	case "UserBasic":
		timeout = time.Hour // 前台用户过期时间
		payloadFunc = func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.UserBasic); ok {
				encryptedID, _ := utils.Encrypt(strconv.FormatUint(uint64(v.ID), 10), encryptionKey)
				encryptedName, _ := utils.Encrypt(v.Name, encryptionKey)

				err := utils.Redis.Set(context.Background(), encryptedID, v.ID, timeout).Err()
				if err != nil {
					log.Printf("Failed to store user info in Redis: %v", err)
				}

				return jwt.MapClaims{
					"id":   encryptedID,
					"name": encryptedName,
				}
			}
			return jwt.MapClaims{}
		}

	case "Admin":
		timeout = 2 * time.Hour // 后台用户过期时间
		payloadFunc = func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.Admin); ok {
				encryptedID, _ := utils.Encrypt(strconv.FormatUint(uint64(v.ID), 10), encryptionKey)
				encryptedName, _ := utils.Encrypt(v.Name, encryptionKey)

				err := utils.Redis.Set(context.Background(), encryptedID, v.ID, timeout).Err()
				if err != nil {
					log.Printf("Failed to store user info in Redis: %v", err)
				}

				return jwt.MapClaims{
					"id":   encryptedID,
					"name": encryptedName,
				}
			}
			return jwt.MapClaims{}
		}

	default:
		log.Fatal("Unknown model type")
	}

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     timeout,
		MaxRefresh:  timeout,
		IdentityKey: identityKey,
		PayloadFunc: payloadFunc,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return &JwtMiddlewareWrapper{authMiddleware}
}

func (mw *JwtMiddlewareWrapper) GenerateToken(user interface{}) (string, int64, error) {
	token, expire, err := mw.TokenGenerator(user)
	if err != nil {
		return "", 0, err
	}
	expireTimestamp := expire.Unix()
	return token, expireTimestamp, nil
}
