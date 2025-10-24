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

// adminServerCmd represents the admin command
var adminServerCmd = &cobra.Command{
	Use:   "admin",
	Short: "Run Admin server",
	Long:  "Start the Admin server for management endpoints",
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

		// 注册Admin路由
		router.AdminRouter(r)

		// 启动Admin服务器
		adminPort := utils.GetAdminPort()
		fmt.Printf("Admin Server starting on port %d\n", adminPort)
		if err := r.Run(fmt.Sprintf(":%d", adminPort)); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(adminServerCmd)
}
