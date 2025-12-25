package call

import (
	"gochat/internal/component"
	"gochat/internal/infrastructure/dao"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/pkg/response"
	"gochat/utils"
)

type CallHandler struct {
	response *response.SvcRequest
	dao      *dao.CallDao
}

func NewCallHandler(server *component.ApiServer) *CallHandler {
	handler := &CallHandler{
		response: server.Result,
		dao:      dao.NewCallDao(utils.DB),
	}

	// 注册信令处理器到WebSocket层
	websocket.SetCallSignalProcessor(handler)

	// 初始化通话超时管理器
	if wsHub := getWebSocketHub(); wsHub != nil {
		// 从配置文件读取超时时间，默认60秒
		timeoutSeconds := 60 // 可以从viper读取
		timeoutManager := websocket.NewCallTimeoutManager(timeoutSeconds, handler.dao, wsHub)
		websocket.SetCallTimeoutManager(timeoutManager)
		// 设置到utils以便其他模块访问
		utils.SetCallTimeoutManager(timeoutManager)
	}

	return handler
}

