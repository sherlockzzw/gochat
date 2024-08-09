package router

import (
	"github.com/gin-gonic/gin"
	"gochat/service"
)

func Router() *gin.Engine {
	r := gin.Default()
	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		panic(err)
	}

	r.GET("/index", service.Index)
	return r
}
