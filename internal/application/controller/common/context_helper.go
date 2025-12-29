package common

import (
	"gochat/internal/pkg/code_msg"
	"gochat/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetUserIDFromContext 从上下文中获取用户ID（统一处理，支持前台用户）
// 返回: (userID, errCode, error)
func GetUserIDFromContext(ctx *gin.Context) (int64, code_msg.BusinessCode, error) {
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return 0, code_msg.Unauthorized, err
	}
	if userID <= 0 {
		return 0, code_msg.Unauthorized, nil
	}
	return userID, 0, nil
}
