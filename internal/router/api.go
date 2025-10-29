package router

import (
	"gochat/internal/application/handler"
	"gochat/middleware"

	"github.com/gin-gonic/gin"
)

func ApiRouter(r *gin.Engine) {
	apiHandler := handler.NewApi()
	route := r.Group("api")

	// 静态资源服务
	staticRouter(r)

	publicRouter(route, apiHandler)
	privateRouter(route, apiHandler)
}

// staticRouter 静态资源路由
func staticRouter(r *gin.Engine) {
	// 用户头像静态资源
	r.Static("/static/avatar", "./resources/user/avatar")
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

	// 好友相关接口
	friendRoute := r.Group("friend")
	{
		// 添加好友
		friendRoute.POST("add", api.FriendHandler.AddFriend)

		// 获取好友列表
		friendRoute.GET("list", api.FriendHandler.GetFriendList)

		// 获取好友申请列表
		friendRoute.GET("requests", api.FriendHandler.GetFriendRequests)

		// 处理好友申请
		friendRoute.POST("handle-request", api.FriendHandler.HandleFriendRequest)

		// 删除好友
		friendRoute.DELETE("delete", api.FriendHandler.DeleteFriend)

		// 检查好友关系
		friendRoute.GET("check", api.FriendHandler.CheckFriend)

		// 设置好友备注
		friendRoute.POST("set-remark", api.FriendHandler.SetFriendRemark)

		// 屏蔽好友
		friendRoute.POST("block", api.FriendHandler.BlockFriend)
	}
}
