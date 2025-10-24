package user

import (
	"gochat/utils"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type GetUserInfoRequest struct {
	ID int64 `form:"id" binding:"required"`
}

type GetUserInfoResponse struct {
	Message string      `json:"message"`
	User    interface{} `json:"user"`
}

func (h *ControllerUser) GetUserInfo(ctx *gin.Context) {
	// 从JWT中获取用户信息
	claims := jwt.ExtractClaims(ctx)
	userID := claims["id"].(string)
	userName := claims["name"].(string)

	decryptedID, _ := utils.Decrypt(userID, viper.GetString("token.encryptionKey"))
	decryptedName, _ := utils.Decrypt(userName, viper.GetString("token.encryptionKey"))

	resp := &GetUserInfoResponse{
		Message: "欢迎来到服务",
		User: gin.H{
			"id":   decryptedID,
			"name": decryptedName,
		},
	}

	h.response.JsonSuccess(ctx, resp)
}
