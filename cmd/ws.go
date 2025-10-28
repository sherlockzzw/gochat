package cmd

import (
	"fmt"
	"gochat/internal/component"
	"gochat/internal/router"
	"gochat/utils"
	"os"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// wsServerCmd represents the websocket command
var wsServerCmd = &cobra.Command{
	Use:   "ws",
	Short: "Run WebSocket server (separate port)",
	Long:  "Start the WebSocket server on a separate port for real-time communication",
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered", "recover", r)
				debug.PrintStack()
				os.Exit(1)
			}
		}()

		// 设置WebSocket服务器（内部会初始化所有配置和数据库连接）
		component.SetWebSocketServer()

		// 设置Gin模式
		gin.SetMode(gin.ReleaseMode)

		// 创建路由
		r := gin.New()
		r.Use(gin.Logger(), gin.Recovery())

		// 注册WebSocket路由
		router.WebSocketRouter(r)

		// 启动WebSocket服务器
		wsPort := utils.GetWebSocketPort()
		fmt.Printf("WebSocket Server starting on port %d\n", wsPort)
		fmt.Printf("WebSocket连接地址: ws://localhost:%d/ws/connect?user_id=<用户ID>\n", wsPort)
		fmt.Printf("在线用户查询: http://localhost:%d/ws/online\n", wsPort)
		if err := r.Run(fmt.Sprintf(":%d", wsPort)); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(wsServerCmd)
}
