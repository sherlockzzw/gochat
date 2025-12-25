package common

import (
	"gochat/internal/infrastructure/websocket"
	globalUtils "gochat/utils"
)

// GetWebSocketHub 获取WebSocket Hub（带类型断言）
func GetWebSocketHub() *websocket.Hub {
	if hub := globalUtils.GetWebSocketHub(); hub != nil {
		if wsHub, ok := hub.(*websocket.Hub); ok {
			return wsHub
		}
	}
	return nil
}

// IsUserOnline 检查用户是否在线
func IsUserOnline(userID int64) bool {
	wsHub := GetWebSocketHub()
	if wsHub == nil {
		return false
	}
	return wsHub.IsUserOnline(userID)
}

// GetOnlineUserIDs 获取在线用户ID列表
func GetOnlineUserIDs() []int64 {
	wsHub := GetWebSocketHub()
	if wsHub == nil {
		return []int64{}
	}
	return wsHub.GetOnlineUserIDs()
}

// SendToUser 发送消息给指定用户
func SendToUser(userID int64, message []byte) {
	wsHub := GetWebSocketHub()
	if wsHub == nil {
		return
	}
	wsHub.SendToUser(userID, message)
}

