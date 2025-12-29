package auth

import (
	"gochat/api/admin/auth"

	"github.com/gin-gonic/gin"
)

// Logout 管理员登出
func (c *AuthController) Logout(ctx *gin.Context) {
	// TODO: 实现token失效逻辑
	c.response.JsonSuccess(ctx, &auth.LogoutResponse{})
}

