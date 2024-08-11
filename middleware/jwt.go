package middleware

import (
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gochat/models"
	"gochat/utils"
	"gorm.io/gorm"
	"log"
	"time"
)

var identityKey = "id"
var encryptionKey = "examplekey946666" // 16字节密钥

type JwtMiddlewareWrapper struct {
	*jwt.GinJWTMiddleware
}

func JwtMiddleware() *JwtMiddlewareWrapper {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.UserBasic); ok {
				encryptedID, _ := utils.Encrypt(strconv.FormatUint(uint64(v.ID), 10), encryptionKey)
				encryptedName, _ := utils.Encrypt(v.Name, encryptionKey)
				return jwt.MapClaims{
					"id":   encryptedID,
					"name": encryptedName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			decryptedIDStr, _ := utils.Decrypt(claims["id"].(string), encryptionKey)
			decryptedID64, _ := strconv.ParseUint(decryptedIDStr, 10, 32)
			decryptedName, _ := utils.Decrypt(claims["name"].(string), encryptionKey)
			return &models.UserBasic{
				Model: gorm.Model{
					ID: uint(decryptedID64), // 将 uint64 转换为 uint
				},
				Name: decryptedName,
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			// 通过你自己的方法验证用户
			user, err := models.GetUserByName(loginVals.Username)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			return user, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*models.UserBasic); ok && v.Name == "admin" {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return &JwtMiddlewareWrapper{authMiddleware}
}

func (mw *JwtMiddlewareWrapper) GenerateToken(user *models.UserBasic) (string, time.Time, error) {
	token, expire, err := mw.TokenGenerator(user)
	if err != nil {
		return "", time.Time{}, err
	}
	return token, expire, nil
}

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
