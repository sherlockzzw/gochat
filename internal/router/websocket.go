package router

import (
	"fmt"
	"gochat/internal/component"
	"gochat/internal/infrastructure/websocket"
	"gochat/utils"

	"github.com/gin-gonic/gin"
)

// WebSocketRouter 注册WebSocket路由
func WebSocketRouter(r *gin.Engine) {
	hub, ok := utils.WSHub.(*websocket.Hub)
	if !ok || hub == nil {

		fmt.Println("Warning: WebSocket Hub not initialized, attempting to initialize...")
		component.GetApiServer()
		hub, ok = utils.WSHub.(*websocket.Hub)
		if !ok || hub == nil {
			fmt.Println("Error: Failed to initialize WebSocket Hub")
			return
		}
	}

	// 创建WebSocket处理器
	wsHandler := websocket.NewWebSocketHandler(hub)

	// WebSocket路由组
	ws := r.Group("/ws")
	{
		// 连接WebSocket
		ws.GET("/connect", wsHandler.HandleWebSocket)

		// 获取在线用户列表
		ws.GET("/online", wsHandler.GetOnlineUsers())

		// 检查特定用户是否在线
		ws.GET("/online/:user_id", wsHandler.IsUserOnline())

		// 健康检查
		ws.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":       "ok",
				"service":      "websocket",
				"online_users": hub.GetOnlineUsers(),
			})
		})
	}
}
