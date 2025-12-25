package websocket

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
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
	// 从查询参数获取token
	token := c.Query("token")
	if token == "" {
		h.logger.Error("Missing token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	// 从JWT token中解析用户ID
	userID, err := h.parseUserIDFromToken(token)
	if err != nil {
		h.logger.Error("Failed to parse user ID from token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// 如果token解析失败，尝试从query参数获取（向后兼容）
	if userID == 0 {
		userIDStr := c.Query("user_id")
		if userIDStr != "" {
			parsedID, parseErr := strconv.ParseUint(userIDStr, 10, 32)
			if parseErr == nil {
				userID = int64(parsedID)
			}
		}
	}

	if userID == 0 {
		h.logger.Error("Invalid user ID")
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
	// 增加send通道缓冲大小，提高并发性能
	// 256可能在高并发下不够，增加到1024
	client := &Client{
		conn:   conn,
		send:   make(chan []byte, 1024),
		userID: int64(userID),
		hub:    h.hub,
	}

	// 注册客户端
	fmt.Printf("准备注册WebSocket客户端，用户ID: %d\n", uint(userID))
	client.hub.register <- client

	// 启动读写协程
	go client.writePump()
	go client.readPump()

	fmt.Printf("WebSocket连接已建立，用户ID: %d\n", uint(userID))
	h.logger.Info("WebSocket connection established",
		zap.Int64("userID", int64(userID)))
}

// parseUserIDFromToken 从JWT token中解析用户ID
func (h *WebSocketHandler) parseUserIDFromToken(tokenString string) (int64, error) {
	// 获取加密密钥
	encryptionKey := viper.GetString("token.encryptionKey")
	if encryptionKey == "" {
		return 0, jwt.ErrSignatureInvalid
	}

	// 解析token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(encryptionKey), nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, jwt.ErrSignatureInvalid
	}

	// 获取claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, jwt.ErrSignatureInvalid
	}

	// 根据JWT中间件的实现，用户信息在"identity"字段中
	// identity字段存储的是UserBasic对象，解析后会是一个map
	if identity, exists := claims["identity"]; exists {
		if identityMap, ok := identity.(map[string]interface{}); ok {
			// UserBasic结构体中的ID字段是大写的"ID"
			if id, exists := identityMap["ID"]; exists {
				switch v := id.(type) {
				case float64:
					return int64(v), nil
				case int64:
					return v, nil
				case int:
					return int64(v), nil
				case uint:
					return int64(v), nil
				case uint64:
					return int64(v), nil
				case string:
					parsedID, err := strconv.ParseInt(v, 10, 64)
					if err == nil {
						return parsedID, nil
					}
				}
			}
		}
	}

	// 向后兼容：尝试从"id"字段获取（某些token格式可能直接包含id）
	if id, exists := claims["id"]; exists {
		switch v := id.(type) {
		case float64:
			return int64(v), nil
		case int64:
			return v, nil
		case int:
			return int64(v), nil
		case uint:
			return int64(v), nil
		case uint64:
			return int64(v), nil
		case string:
			parsedID, err := strconv.ParseInt(v, 10, 64)
			if err == nil {
				return parsedID, nil
			}
		}
	}

	return 0, jwt.ErrSignatureInvalid
}

// SendMessageToUser 发送消息给特定用户
func (h *WebSocketHandler) SendMessageToUser(userID uint, message []byte) {
	h.hub.SendToUser(int64(userID), message)
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
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		isOnline := h.hub.IsUserOnline(userID)
		c.JSON(http.StatusOK, gin.H{
			"user_id":   userID,
			"is_online": isOnline,
		})
	}
}
