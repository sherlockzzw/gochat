package utils

import (
	"fmt"
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	e := viper.ReadInConfig()
	if e != nil {
		fmt.Println("配置文件读取错误")
		panic(e)
	}
	fmt.Println("config:", viper.Get("mysql"))
}
func InitMysql() {

}
