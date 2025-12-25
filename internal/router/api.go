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

	// 注册WebSocket路由（需要在认证路由之前）
	WebSocketRouter(r)

	publicRouter(route, apiHandler)
	privateRouter(route, apiHandler)
}

// staticRouter 静态资源路由
func staticRouter(r *gin.Engine) {
	// 用户头像静态资源
	r.Static("/static/avatar", "./resources/user/avatar")
	// 上传文件静态资源
	r.Static("/static/upload", "./resources/upload")
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
	userRoute.POST("profile", api.UserHandler.UpdateUserProfile)

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
		chatRoute.POST("message/recall", api.ChatHandler.RecallMessage)
		chatRoute.POST("message/delete", api.ChatHandler.DeleteMessage)
		chatRoute.POST("message/forward", api.ChatHandler.ForwardMessage)
		chatRoute.POST("message/favorite", api.ChatHandler.FavoriteMessage)
		chatRoute.POST("message/unfavorite", api.ChatHandler.UnfavoriteMessage)
		chatRoute.GET("message/favorites", api.ChatHandler.GetFavoriteMessages)

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

		// 获取好友详情
		friendRoute.GET("detail", api.FriendHandler.GetFriendDetail)

		// 设置好友分组
		friendRoute.POST("set-group", api.FriendHandler.SetFriendGroup)
	}

	// 群组相关接口
	groupRoute := r.Group("group")
	{
		// 创建群组
		groupRoute.POST("create", api.GroupHandler.CreateGroup)

		// 获取群组信息
		groupRoute.GET("info", api.GroupHandler.GetGroupInfo)

		// 获取用户群组列表
		groupRoute.GET("list", api.GroupHandler.GetUserGroups)

		// 添加群成员
		groupRoute.POST("members/add", api.GroupHandler.AddGroupMembers)

		// 移除群成员
		groupRoute.DELETE("members/remove", api.GroupHandler.RemoveGroupMember)

		// 检查用户是否在群组中
		groupRoute.GET("members/check", api.GroupHandler.CheckMemberInGroup)

		// 更新群组信息（仅群主）
		groupRoute.POST("update", api.GroupHandler.UpdateGroupInfo)

		// 群主权限：任命管理员
		groupRoute.POST("admin/appoint", api.GroupHandler.AppointAdmin)

		// 群主权限：移除管理员
		groupRoute.POST("admin/remove", api.GroupHandler.RemoveAdmin)

		// 群主/管理员权限：禁言成员
		groupRoute.POST("mute", api.GroupHandler.MuteMember)

		// 群主/管理员权限：解除禁言
		groupRoute.POST("unmute", api.GroupHandler.UnmuteMember)

		// 群主权限：解散群组
		groupRoute.POST("dissolve", api.GroupHandler.DissolveGroup)

		// 群主权限：转让群组
		groupRoute.POST("transfer", api.GroupHandler.TransferGroup)

		// 成员权限：退出群组
		groupRoute.POST("leave", api.GroupHandler.LeaveGroup)

		// 管理员权限：删除消息
		groupRoute.POST("message/delete", api.GroupHandler.DeleteGroupMessage)
	}

	// 资金相关接口
	balanceRoute := r.Group("balance")
	{
		// 余额查询
		balanceRoute.GET("", api.BalanceHandler.GetBalance)

		// 充值
		balanceRoute.POST("recharge", api.BalanceHandler.Recharge)
		balanceRoute.GET("recharge/list", api.BalanceHandler.GetRechargeRequests)

		// 提现
		balanceRoute.POST("withdraw", api.BalanceHandler.Withdraw)
		balanceRoute.GET("withdraw/list", api.BalanceHandler.GetWithdrawRequests)

		// 红包
		balanceRoute.POST("redpacket/private/send", api.BalanceHandler.SendPrivateRedPacket)
		balanceRoute.POST("redpacket/group/send", api.BalanceHandler.SendGroupRedPacket)
		balanceRoute.POST("redpacket/receive", api.BalanceHandler.ReceiveRedPacket)
		balanceRoute.GET("redpacket/detail", api.BalanceHandler.GetRedPacketDetail)

		// 转账
		balanceRoute.POST("transfer/private/send", api.BalanceHandler.SendPrivateTransfer)
		balanceRoute.POST("transfer/group/send", api.BalanceHandler.SendGroupTransfer)
		balanceRoute.POST("transfer/receive", api.BalanceHandler.ReceiveTransfer)
		balanceRoute.POST("transfer/cancel", api.BalanceHandler.CancelTransfer)
		balanceRoute.GET("transfer/detail", api.BalanceHandler.GetTransferDetail)

		// 资金流水
		balanceRoute.GET("flows", api.BalanceHandler.GetBalanceFlows)
	}

	// 通话相关接口
	callRoute := r.Group("call")
	{
		// 发起私聊通话
		callRoute.POST("start-private", api.CallHandler.StartPrivateCall)

		// 发起群聊通话
		callRoute.POST("start-group", api.CallHandler.StartGroupCall)

		// 接受通话
		callRoute.POST("accept", api.CallHandler.AcceptCall)

		// 拒绝通话
		callRoute.POST("reject", api.CallHandler.RejectCall)

		// 取消通话
		callRoute.POST("cancel", api.CallHandler.CancelCall)

		// 结束通话
		callRoute.POST("end", api.CallHandler.EndCall)

		// 获取通话记录
		callRoute.GET("records", api.CallHandler.GetCallRecords)
	}

	// 通知相关接口
	notificationRoute := r.Group("notification")
	{
		// 获取通知列表
		notificationRoute.GET("list", api.NotificationHandler.GetNotifications)

		// 标记通知已读
		notificationRoute.POST("read", api.NotificationHandler.MarkAsRead)

		// 标记所有通知已读
		notificationRoute.POST("read-all", api.NotificationHandler.MarkAllAsRead)

		// 获取未读通知数量
		notificationRoute.GET("unread-count", api.NotificationHandler.GetUnreadCount)

		// 删除通知
		notificationRoute.POST("delete", api.NotificationHandler.DeleteNotification)
	}

	// 会话相关接口
	conversationRoute := r.Group("conversation")
	{
		// 获取会话列表
		conversationRoute.GET("list", api.ConversationHandler.GetConversations)

		// 设置会话置顶
		conversationRoute.POST("pin", api.ConversationHandler.PinConversation)

		// 设置会话静音
		conversationRoute.POST("mute", api.ConversationHandler.MuteConversation)

		// 设置会话隐藏
		conversationRoute.POST("hide", api.ConversationHandler.HideConversation)

		// 删除会话
		conversationRoute.POST("delete", api.ConversationHandler.DeleteConversation)

		// 清除未读数
		conversationRoute.POST("clear-unread", api.ConversationHandler.ClearUnreadCount)
	}
}
