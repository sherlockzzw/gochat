package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func MethodCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取预期的请求方法
		expectedMethod := c.GetString("expected_method")
		if expectedMethod == "" {
			c.Next()
			return
		}

		// 检查请求方法是否正确
		if c.Request.Method != expectedMethod {
			c.JSON(http.StatusMethodNotAllowed, gin.H{
				"error": "Method not allowed",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
