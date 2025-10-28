package websocket

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有来源，生产环境应该检查来源
		return true
	},
}

// WebSocketHandler WebSocket处理器
type WebSocketHandler struct {
	hub    *Hub
	logger *zap.Logger
}

// NewWebSocketHandler 创建WebSocket处理器
func NewWebSocketHandler(hub *Hub) *WebSocketHandler {
	logger, _ := zap.NewProduction()
	return &WebSocketHandler{
		hub:    hub,
		logger: logger,
	}
}

// HandleWebSocket 处理WebSocket连接
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// 从查询参数获取用户ID
	userIDStr := c.Query("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid user ID", zap.String("userID", userIDStr), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 升级HTTP连接为WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("Failed to upgrade connection", zap.Error(err))
		return
	}

	// 创建客户端
	client := &Client{
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: uint(userID),
		hub:    h.hub,
	}

	// 注册客户端
	client.hub.register <- client

	// 启动读写协程
	go client.writePump()
	go client.readPump()

	h.logger.Info("WebSocket connection established",
		zap.Uint("userID", uint(userID)))
}

// SendMessageToUser 发送消息给特定用户
func (h *WebSocketHandler) SendMessageToUser(userID uint, message []byte) {
	h.hub.SendToUser(userID, message)
}

// GetOnlineUsers 获取在线用户信息
func (h *WebSocketHandler) GetOnlineUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDs := h.hub.GetOnlineUserIDs()
		c.JSON(http.StatusOK, gin.H{
			"online_users": userIDs,
			"count":        len(userIDs),
		})
	}
}

// IsUserOnline 检查用户是否在线
func (h *WebSocketHandler) IsUserOnline() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("user_id")
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		isOnline := h.hub.IsUserOnline(uint(userID))
		c.JSON(http.StatusOK, gin.H{
			"user_id":   uint(userID),
			"is_online": isOnline,
		})
	}
}
