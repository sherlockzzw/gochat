package middleware

import (
	"context"
	"encoding/json"
	jwt "github.com/appleboy/gin-jwt/v2"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"gochat/models"
	"gochat/utils"
	"log"
	"strconv"
	"time"
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
	expireHours := viper.GetInt("token.expire")
	expireTime := time.Duration(expireHours) * time.Hour

	var payloadFunc func(data interface{}) jwt.MapClaims

	switch modelType {
	case "UserBasic":
		payloadFunc = func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.UserBasic); ok {
				encryptedID, _ := utils.Encrypt(strconv.FormatUint(uint64(v.ID), 10), encryptionKey)
				encryptedName, _ := utils.Encrypt(v.Name, encryptionKey)

				redisKey := modelType + ":" + encryptedID
				userInfo := RedisUserInfo{
					ID:   encryptedID,
					Name: encryptedName,
				}
				userInfoJSON, _ := json.Marshal(userInfo)
				err := utils.Redis.Set(context.Background(), redisKey, userInfoJSON, expireTime).Err()
				if err != nil {
					log.Printf("Failed to store user info in Redis: %v", err)
				}
				return jwt.MapClaims{
					"id":   v.ID,
					"name": v.Name,
				}
			}
			return jwt.MapClaims{}
		}

	case "Admin":
		payloadFunc = func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.Admin); ok {

				return jwt.MapClaims{
					"id":   v.ID,
					"name": v.Name,
				}
			}
			return jwt.MapClaims{}
		}

	default:
		log.Fatal("Unknown model type")
	}

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       modelType,
		Key:         []byte(encryptionKey),
		Timeout:     expireTime,
		MaxRefresh:  expireTime,
		IdentityKey: identityKey,
		PayloadFunc: payloadFunc,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return &JwtMiddlewareWrapper{authMiddleware}
}
func (mw *JwtMiddlewareWrapper) GenerateToken(data interface{}, customExpire time.Duration) (string, int64, error) {
	token := jwtv4.NewWithClaims(jwtv4.SigningMethodHS256, jwtv4.MapClaims{
		"identity": data,
		"exp":      time.Now().Add(customExpire).Unix(),
	})
	tokenString, err := token.SignedString(mw.Key)
	if err != nil {
		return "", 0, err
	}
	return tokenString, time.Now().Add(customExpire).Unix(), nil
}
