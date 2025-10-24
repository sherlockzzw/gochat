package component

import (
	"gochat/internal/pkg/response"
	"gochat/utils"

	"go.uber.org/zap"
)

type ApiServer struct {
	MysqlSvc *utils.MysqlService
	RedisSvc *utils.RedisService
	Result   *response.SvcRequest
	Logger   *zap.Logger
}

func NewApiServer() *ApiServer {
	// 初始化日志
	logger, _ := zap.NewProduction()

	// 初始化响应处理器
	svcRequest := response.NewSvcRequest()

	return &ApiServer{
		MysqlSvc: &utils.MysqlService{DB: utils.DB},
		RedisSvc: &utils.RedisService{Client: utils.RDB},
		Result:   svcRequest,
		Logger:   logger,
	}
}

func GetApiServer() *ApiServer {
	return NewApiServer()
}
