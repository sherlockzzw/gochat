package main

import (
	"gochat/router"
	"gochat/utils"
)

func main() {
	//初始化配置
	utils.InitConfig()
	utils.InitMysql()
	utils.InitRedis()
	//路由
	r := router.Router()

	r.Run(":8080")
}
