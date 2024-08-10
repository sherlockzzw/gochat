package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Configuration file was incorrectly read")
		panic(err)
	}
	//fmt.Println("config:", viper.Get("mysql"))
}

func InitMysql() {
	//sql语句打印
	sqlLog := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.database"),
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: sqlLog})
	if err != nil {
		fmt.Println("Database connection error")
		panic(err)
	}

	// 测试连接
	sqlDB, err := DB.DB()
	if err != nil {
		fmt.Println("Failed to get SQL DB")
		panic(err)
	}
	err = sqlDB.Ping()
	if err != nil {
		fmt.Println("Failed to ping database")
		panic(err)
	}

	fmt.Println("Successfully connected to MySQL")
}

var RedisClient *redis.Client

func InitRedis() {
	addr := fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port"))
	options := &redis.Options{
		Addr:     addr,
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	}

	RedisClient = redis.NewClient(options)

	ctx := context.Background()
	pong, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis")
		panic(err)
	}

	fmt.Println("Successfully connected to Redis:", pong)
}
