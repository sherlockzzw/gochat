package component

import (
	"fmt"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/pkg/response"
	"gochat/utils"

	"go.uber.org/zap"
)

var (
	setupServer     *SetupServer
	apiServer       *ApiServer
	adminServer     *AdminServer
	websocketServer *WebSocketServer
)

// initCoreServices 初始化核心服务和数据库连接
func initCoreServices() {
	utils.InitConfig()
	utils.InitMysql()
	utils.InitRedis()
	utils.InitMongoDB()
}

// initCoreServicesForSetup 初始化setup命令所需的服务（仅MySQL）
func initCoreServicesForSetup() {
	utils.InitConfig()
	utils.InitMysql()
	// Redis 和 MongoDB 在 setup 时不是必需的，可选初始化
	utils.InitRedis()
	utils.InitMongoDB()
}

// initWebSocketService 初始化WebSocket服务
func initWebSocketService() {
	// 实际创建WebSocket Hub
	if utils.WSHub == nil {
		hub := websocket.NewHub()
		utils.WSHub = hub
		// 启动Hub的goroutine
		go hub.Run()
		fmt.Println("WebSocket Hub initialized and started")
	}
}

// createLogger 创建日志实例
func createLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	return logger
}

// createSvcRequest 创建响应处理器
func createSvcRequest() *response.SvcRequest {
	return response.NewSvcRequest()
}

type ApiServer struct {
	MysqlSvc *utils.MysqlService
	RedisSvc *utils.RedisService
	Result   *response.SvcRequest
	Logger   *zap.Logger
}

func NewApiServer() *ApiServer {
	initCoreServices()
	initWebSocketService()

	return &ApiServer{
		MysqlSvc: &utils.MysqlService{DB: utils.DB},
		RedisSvc: &utils.RedisService{Client: utils.RDB},
		Result:   createSvcRequest(),
		Logger:   createLogger(),
	}
}

// SetupServer 数据库设置服务器
type SetupServer struct {
	MysqlSvc *utils.MysqlService
	RedisSvc *utils.RedisService
	MongoSvc *utils.MongoService
	Logger   *zap.Logger
}

// NewSetupServer 创建设置服务器
func NewSetupServer() *SetupServer {
	initCoreServicesForSetup()
	// WebSocket 在 setup 时不是必需的，跳过初始化

	return &SetupServer{
		MysqlSvc: &utils.MysqlService{DB: utils.DB},
		RedisSvc: &utils.RedisService{Client: utils.RDB}, // 可能为 nil
		MongoSvc: &utils.MongoService{DB: utils.MongoDB}, // 可能为 nil
		Logger:   createLogger(),
	}
}

// SetSetUpServer 设置全局设置服务器
func SetSetUpServer() {
	setupServer = NewSetupServer()
}

// GetSetUpServer 获取设置服务器
func GetSetUpServer() *SetupServer {
	if setupServer == nil {
		SetSetUpServer()
	}
	return setupServer
}

// AdminServer Admin服务器
type AdminServer struct {
	MysqlSvc *utils.MysqlService
	RedisSvc *utils.RedisService
	Result   *response.SvcRequest
	Logger   *zap.Logger
}

// NewAdminServer 创建Admin服务器
func NewAdminServer() *AdminServer {
	initCoreServices()

	return &AdminServer{
		MysqlSvc: &utils.MysqlService{DB: utils.DB},
		RedisSvc: &utils.RedisService{Client: utils.RDB},
		Result:   createSvcRequest(),
		Logger:   createLogger(),
	}
}

// SetApiServer 设置API服务器
func SetApiServer() {
	apiServer = NewApiServer()
}

// GetApiServer 获取API服务器
func GetApiServer() *ApiServer {
	if apiServer == nil {
		SetApiServer()
	}
	return apiServer
}

// SetAdminServer 设置Admin服务器
func SetAdminServer() {
	adminServer = NewAdminServer()
}

// GetAdminServer 获取Admin服务器
func GetAdminServer() *AdminServer {
	if adminServer == nil {
		SetAdminServer()
	}
	return adminServer
}

// WebSocketServer WebSocket服务器
type WebSocketServer struct {
	MysqlSvc *utils.MysqlService
	RedisSvc *utils.RedisService
	MongoSvc *utils.MongoService
	Logger   *zap.Logger
}

// NewWebSocketServer 创建WebSocket服务器
func NewWebSocketServer() *WebSocketServer {
	initCoreServices()
	initWebSocketService() // WebSocket服务需要WebSocket功能

	return &WebSocketServer{
		MysqlSvc: &utils.MysqlService{DB: utils.DB},
		RedisSvc: &utils.RedisService{Client: utils.RDB},
		MongoSvc: &utils.MongoService{DB: utils.MongoDB},
		Logger:   createLogger(),
	}
}

// SetWebSocketServer 设置WebSocket服务器
func SetWebSocketServer() {
	websocketServer = NewWebSocketServer()
}

// GetWebSocketServer 获取WebSocket服务器
func GetWebSocketServer() *WebSocketServer {
	if websocketServer == nil {
		SetWebSocketServer()
	}
	return websocketServer
}
