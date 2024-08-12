package middleware

import (
	"context"
	"encoding/json"
	"github.com/spf13/viper"
	"log"
	"strconv"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"gochat/models"
	"gochat/utils"
)

type JwtMiddlewareWrapper struct {
	*jwt.GinJWTMiddleware
}

type RedisUserInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func JwtMiddleware(modelType string) *JwtMiddlewareWrapper {
	var identityKey = "id"
	var encryptionKey = viper.GetString("token.encryptionKey")
	var timeout time.Duration
	var payloadFunc func(data interface{}) jwt.MapClaims

	switch modelType {
	case "UserBasic":
		timeout = 24 * time.Hour
		payloadFunc = func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.UserBasic); ok {
				encryptedID, _ := utils.Encrypt(strconv.FormatUint(uint64(v.ID), 10), encryptionKey)
				encryptedName, _ := utils.Encrypt(v.Name, encryptionKey)
				redisKey := modelType + ":" + strconv.FormatUint(uint64(v.ID), 10)
				userInfo := RedisUserInfo{
					ID:   encryptedID,
					Name: encryptedName,
				}
				userInfoJSON, _ := json.Marshal(userInfo)
				err := utils.Redis.Set(context.Background(), redisKey, userInfoJSON, timeout).Err()
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
		timeout = 24 * time.Hour
		payloadFunc = func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.Admin); ok {
				encryptedID, _ := utils.Encrypt(strconv.FormatUint(uint64(v.ID), 10), encryptionKey)
				encryptedName, _ := utils.Encrypt(v.Name, encryptionKey)
				redisKey := modelType + ":" + encryptedID
				userInfo := RedisUserInfo{
					ID:   encryptedID,
					Name: encryptedName,
				}
				userInfoJSON, _ := json.Marshal(userInfo)
				err := utils.Redis.Set(context.Background(), redisKey, userInfoJSON, timeout).Err()
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
		Realm:       modelType,
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

func (mw *JwtMiddlewareWrapper) GenerateToken(data interface{}) (string, int64, error) {
	token, expire, err := mw.TokenGenerator(data)
	if err != nil {
		return "", 0, err
	}
	expireTimestamp := expire.Unix()
	return token, expireTimestamp, nil
}
