package middleware

import (
	"github.com/gin-gonic/gin"
)

func MethodCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()
	}
}
