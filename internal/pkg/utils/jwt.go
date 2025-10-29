package utils

import (
	"fmt"
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// GetCurrentUserID 从JWT token中获取当前用户ID
func GetCurrentUserID(ctx *gin.Context) (uint, error) {
	claims := jwt.ExtractClaims(ctx)

	// 检查claims是否存在
	if claims == nil {
		return 0, fmt.Errorf("no claims found in JWT token")
	}

	// 从identity字段中获取用户信息
	identityInterface, exists := claims["identity"]
	if !exists {
		return 0, fmt.Errorf("identity field not found in JWT claims")
	}

	// 检查identity字段是否为nil
	if identityInterface == nil {
		return 0, fmt.Errorf("identity field is nil in JWT claims")
	}

	// 将identity转换为map
	identityMap, ok := identityInterface.(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("identity field is not a map")
	}

	// 从identity map中获取ID字段
	userIDInterface, exists := identityMap["ID"]
	if !exists {
		return 0, fmt.Errorf("ID field not found in identity map")
	}

	// 检查ID字段是否为nil
	if userIDInterface == nil {
		return 0, fmt.Errorf("ID field is nil in identity map")
	}

	// 直接转换为uint类型（从调试信息看，ID是数字类型）
	if userID, ok := userIDInterface.(uint); ok {
		return userID, nil
	}

	// 如果是其他数字类型，尝试转换
	if userID, ok := userIDInterface.(uint64); ok {
		return uint(userID), nil
	}

	if userID, ok := userIDInterface.(int); ok {
		if userID < 0 {
			return 0, fmt.Errorf("ID field is negative")
		}
		return uint(userID), nil
	}

	if userID, ok := userIDInterface.(float64); ok {
		if userID < 0 {
			return 0, fmt.Errorf("ID field is negative")
		}
		return uint(userID), nil
	}

	// 如果是字符串，尝试解析
	if userIDStr, ok := userIDInterface.(string); ok {
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			return 0, fmt.Errorf("failed to parse ID as uint: %v", err)
		}
		return uint(userID), nil
	}

	return 0, fmt.Errorf("ID field has unexpected type: %T", userIDInterface)
}
