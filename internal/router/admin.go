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
	route.POST("login", api.AuthController.Login)
}

func adminPrivateRouter(r *gin.RouterGroup, api *controller.API) {
	// 使用JWT认证中间件（TODO: 应该使用Admin的JWT中间件）
	r.Use(middleware.JwtMiddleware("UserBasic").MiddlewareFunc())

	// 认证相关
	r.POST("logout", api.AuthController.Logout)

	// 用户管理
	userRoute := r.Group("user")
	{
		userRoute.GET("list", api.ControllerUser.GetUserList)
		userRoute.GET("info", api.ControllerUser.GetUserInfo)
		userRoute.POST("add", api.ControllerUser.CreateUser)
		userRoute.PUT("update", api.ControllerUser.UpdateUser)
		userRoute.POST("disable", api.ControllerUser.DisableUser)
	}

	// 权限管理
	permissionRoute := r.Group("permission")
	{
		permissionRoute.GET("groups", api.PermissionController.GetPermissionGroups)
		permissionRoute.POST("group", api.PermissionController.CreatePermissionGroup)
		permissionRoute.PUT("group/:id", api.PermissionController.UpdatePermissionGroup)
		permissionRoute.DELETE("group/:id", api.PermissionController.DeletePermissionGroup)
		permissionRoute.POST("assign", api.PermissionController.AssignUserToGroup)
		permissionRoute.POST("feature/toggle", api.PermissionController.ToggleFeature)
	}

	// 功能配置
	configRoute := r.Group("config")
	{
		configRoute.GET("emoji", api.ConfigController.GetEmojiConfig)
		configRoute.POST("emoji", api.ConfigController.UpdateEmojiConfig)
		configRoute.GET("translate", api.ConfigController.GetTranslateConfig)
		configRoute.POST("translate", api.ConfigController.UpdateTranslateConfig)
		configRoute.GET("group", api.ConfigController.GetGroupConfig)
		configRoute.POST("group", api.ConfigController.UpdateGroupConfig)
		configRoute.GET("redpacket", api.ConfigController.GetRedPacketConfig)
		configRoute.POST("redpacket", api.ConfigController.UpdateRedPacketConfig)
	}

	// 资金管理
	financeRoute := r.Group("finance")
	{
		financeRoute.GET("config", api.FinanceController.GetFinanceConfig)
		financeRoute.POST("config", api.FinanceController.UpdateFinanceConfig)
		financeRoute.GET("records", api.FinanceController.GetFinanceRecords)
		financeRoute.GET("recharge/pending", api.FinanceController.GetPendingRecharges)
		financeRoute.POST("recharge/audit", api.FinanceController.AuditRecharge)
		financeRoute.GET("withdraw/pending", api.FinanceController.GetPendingWithdraws)
		financeRoute.POST("withdraw/audit", api.FinanceController.AuditWithdraw)
		financeRoute.POST("balance/adjust", api.FinanceController.AdjustBalance)
	}

	// 内容安全
	contentRoute := r.Group("content")
	{
		contentRoute.GET("violations", api.ContentController.GetViolations)
		contentRoute.POST("violation/delete", api.ContentController.DeleteViolation)
		contentRoute.GET("abnormal", api.ContentController.GetAbnormalOperations)
		contentRoute.POST("abnormal/block", api.ContentController.BlockAbnormalOperation)
	}

	// 数据统计
	statisticsRoute := r.Group("statistics")
	{
		statisticsRoute.GET("usage", api.StatisticsController.GetUsageStatistics)
		statisticsRoute.GET("finance", api.StatisticsController.GetFinanceStatistics)
		statisticsRoute.GET("terminal", api.StatisticsController.GetTerminalStatistics)
	}

	// 操作日志
	logRoute := r.Group("log")
	{
		logRoute.GET("list", api.LogController.GetLogs)
	}
}
