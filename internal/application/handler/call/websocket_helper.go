package call

import (
	"gochat/internal/infrastructure/websocket"
	globalUtils "gochat/utils"
)

// getWebSocketHub 获取WebSocket Hub（带类型断言）
func getWebSocketHub() *websocket.Hub {
	if hub := globalUtils.GetWebSocketHub(); hub != nil {
		if wsHub, ok := hub.(*websocket.Hub); ok {
			return wsHub
		}
	}
	return nil
}


