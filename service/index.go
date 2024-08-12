package service

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gochat/utils"
	"net/http"
)

func Index(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID := claims["id"].(string)
	userName := claims["name"].(string)

	decryptedID, _ := utils.Decrypt(userID, viper.GetString("token.encryptionKey"))
	decryptedName, _ := utils.Decrypt(userName, viper.GetString("token.encryptionKey"))

	// 返回用户信息
	c.JSON(http.StatusOK, gin.H{
		"message": "欢迎来到服务",
		"user": gin.H{
			"id":   decryptedID,
			"name": decryptedName,
		},
	})
}
