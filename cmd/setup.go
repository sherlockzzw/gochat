package cmd

import (
	"fmt"
	"gochat/models"
	"gochat/utils"

	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup database tables",
	Long:  "Initialize database tables and run migrations",
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化配置
		utils.InitConfig()
		utils.InitMysql()
		utils.InitRedis()
		utils.InitMongoDB()
		utils.InitWebSocket()

		// 自动迁移数据库表
		err := utils.DB.AutoMigrate(
			&models.UserBasic{},
			&models.Admin{},
		)
		if err != nil {
			fmt.Printf("Database migration failed: %v\n", err)
			return
		}

		fmt.Println("Database setup completed successfully!")
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
