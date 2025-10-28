package router

import (
	"gochat/internal/application/controller"
	"gochat/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRouter(r *gin.Engine) {
	apiController := controller.NewAdmin()
	route := r.Group("admin")

	adminPublicRouter(route, apiController)
	adminPrivateRouter(route, apiController)
}

func adminPublicRouter(route *gin.RouterGroup, api *controller.API) {
	// 公开接口，无需认证
}

func adminPrivateRouter(r *gin.RouterGroup, api *controller.API) {
	// 使用JWT认证中间件
	r.Use(middleware.JwtMiddleware("UserBasic").MiddlewareFunc())

	// 用户管理
	userRoute := r.Group("user")
	userRoute.GET("list", api.ControllerUser.GetUserList)
}
