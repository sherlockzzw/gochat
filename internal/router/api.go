package router

import (
	"gochat/internal/application/handler"
	"gochat/middleware"

	"github.com/gin-gonic/gin"
)

func ApiRouter(r *gin.Engine) {
	apiHandler := handler.NewApi()
	route := r.Group("api")

	publicRouter(route, apiHandler)
	privateRouter(route, apiHandler)
}

func publicRouter(route *gin.RouterGroup, api *handler.API) {
	// 公开接口，无需认证
	route.POST("login", api.UserHandler.Login)
	route.POST("user/register", api.UserHandler.Register)
}

func privateRouter(r *gin.RouterGroup, api *handler.API) {
	// 使用JWT认证中间件
	r.Use(middleware.JwtMiddleware("UserBasic").MiddlewareFunc())

	// 用户相关接口
	userRoute := r.Group("user")
	userRoute.GET("info", api.UserHandler.GetUserInfo)

	// 聊天相关接口
	chatRoute := r.Group("chat")
	{
		// 用户搜索
		chatRoute.GET("search", api.ChatHandler.SearchUser)

		// 消息相关
		chatRoute.POST("message/send", api.ChatHandler.SendMessage)
		chatRoute.GET("message/history", api.ChatHandler.GetMessageHistory)
		chatRoute.POST("message/read", api.ChatHandler.MarkMessageRead)
		chatRoute.GET("message/unread", api.ChatHandler.GetUnreadCount)

		// 文件上传
		chatRoute.POST("upload", api.ChatHandler.UploadFile)
	}
}
