package utils

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"gochat/internal/infrastructure/websocket"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB      *gorm.DB
	RDB     *redis.Client
	MongoDB *mongo.Database
	WSHub   *websocket.Hub
)

// MysqlService MySQL服务
type MysqlService struct {
	DB *gorm.DB
}

// RedisService Redis服务
type RedisService struct {
	Client *redis.Client
}

// MongoService MongoDB服务
type MongoService struct {
	Database *mongo.Database
}

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

func InitRedis() {
	addr := fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port"))
	options := &redis.Options{
		Addr:     addr,
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	}

	RDB = redis.NewClient(options)

	ctx := context.Background()
	pong, err := RDB.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis")
		panic(err)
	}

	fmt.Println("Successfully connected to Redis:", pong)
}

func InitMongoDB() {
	uri := viper.GetString("mongodb.uri")
	database := viper.GetString("mongodb.database")

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println("Failed to connect to MongoDB")
		panic(err)
	}

	// 测试连接
	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println("Failed to ping MongoDB")
		panic(err)
	}

	MongoDB = client.Database(database)
	fmt.Println("Successfully connected to MongoDB")
}

// InitWebSocket 初始化WebSocket Hub
func InitWebSocket() {
	WSHub = websocket.NewHub()
	go WSHub.Run()
	fmt.Println("WebSocket Hub initialized")
}

// IsDevMode 判断是否为开发模式
func IsDevMode() bool {
	return viper.GetString("app.mode") == "dev"
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// GetApiPort 获取API服务端口
func GetApiPort() int {
	port := viper.GetInt("app.api_port")
	if port == 0 {
		port = 8080 // 默认端口
	}
	return port
}

// GetAdminPort 获取Admin服务端口
func GetAdminPort() int {
	port := viper.GetInt("app.admin_port")
	if port == 0 {
		port = 8081 // 默认端口
	}
	return port
}

// GetWebSocketPort 获取WebSocket服务端口
func GetWebSocketPort() int {
	port := viper.GetInt("websocket.port")
	if port == 0 {
		port = 8082 // 默认端口
	}
	return port
}
