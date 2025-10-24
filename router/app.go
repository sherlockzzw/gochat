package router

import (
	"gochat/docs"
	"gochat/internal/router"
	"gochat/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 使用新的路由结构
	router.ApiRouter(r)
	router.AdminRouter(r)

	// 保留原有的index接口
	auth := r.Group("/")
	auth.Use(gin.Recovery())
	{
		auth.GET("/index", service.Index)
	}

	return r
}
