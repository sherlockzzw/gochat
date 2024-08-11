package router

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gochat/docs"
	"gochat/middleware"
	"gochat/service"
)

func Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/index", service.Index)
	r.GET("/user/list", service.UserList)

	r.POST("/user/login", service.UserLogin)
	auth := r.Group("/")
	auth.Use(middleware.JwtMiddleware().MiddlewareFunc())
	{
		auth.POST("/user/add", service.CreateUser)
	}

	return r
}
