package cmd

import (
	"fmt"
	"gochat/internal/component"
	"gochat/models"

	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup database tables and run migrations",
	Long:  "Initialize database tables and run migrations. This command should be run before starting the services.",
	Run: func(cmd *cobra.Command, args []string) {
		setup(cmd, args)
	},
}

func setup(cmd *cobra.Command, args []string) {
	fmt.Println("🔧 开始数据库初始化...")

	// 设置组件服务器（内部会初始化所有配置和数据库连接）
	component.SetSetUpServer()

	// 获取MySQL服务
	msl := component.GetSetUpServer().MysqlSvc.GetDB()

	fmt.Println("📊 执行数据库迁移...")

	// 执行数据库迁移
	_ = msl.Set("gorm:table_options", "COMMENT='用户基础信息表'").AutoMigrate(&models.UserBasic{})
	_ = msl.Set("gorm:table_options", "COMMENT='管理员信息表'").AutoMigrate(&models.Admin{})
	_ = msl.Set("gorm:table_options", "COMMENT='聊天消息表'").AutoMigrate(&models.ChatMessage{})
	_ = msl.Set("gorm:table_options", "COMMENT='会话表'").AutoMigrate(&models.Conversation{})
	_ = msl.Set("gorm:table_options", "COMMENT='消息已读状态表'").AutoMigrate(&models.MessageReadStatus{})
	_ = msl.Set("gorm:table_options", "COMMENT='好友关系表'").AutoMigrate(&models.Friend{})
	_ = msl.Set("gorm:table_options", "COMMENT='好友申请表'").AutoMigrate(&models.FriendRequest{})
	_ = msl.Set("gorm:table_options", "COMMENT='群组表'").AutoMigrate(&models.Group{})
	_ = msl.Set("gorm:table_options", "COMMENT='群组成员表'").AutoMigrate(&models.GroupMember{})

	fmt.Println("✅ 数据库结构更新完成")
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
