package common

import (
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetUserIDFromContext 从上下文中获取用户ID（统一处理）
// 返回: (userID, errCode, error)
func GetUserIDFromContext(ctx *gin.Context) (int64, code_msg.BusinessCode, error) {
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return 0, code_msg.ServerError, err
	}
	return userID, 0, nil
}

// ValidateUserID 验证用户ID是否有效
func ValidateUserID(userID int64) bool {
	return userID > 0
}


