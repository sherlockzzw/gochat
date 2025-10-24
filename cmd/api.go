package cmd

import (
	"fmt"
	"gochat/internal/router"
	"gochat/utils"
	"os"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// apiServerCmd represents the api command
var apiServerCmd = &cobra.Command{
	Use:   "api",
	Short: "Run API server",
	Long:  "Start the API server for user endpoints",
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered", "recover", r)
				debug.PrintStack()
				os.Exit(1)
			}
		}()

		// 初始化配置
		utils.InitConfig()
		utils.InitMysql()
		utils.InitRedis()

		// 设置Gin模式
		gin.SetMode(gin.ReleaseMode)

		// 创建路由
		r := gin.New()
		r.Use(gin.Logger(), gin.Recovery())

		// 注册API路由
		router.ApiRouter(r)

		// 启动API服务器
		apiPort := utils.GetApiPort()
		fmt.Printf("API Server starting on port %d\n", apiPort)
		if err := r.Run(fmt.Sprintf(":%d", apiPort)); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(apiServerCmd)
}
