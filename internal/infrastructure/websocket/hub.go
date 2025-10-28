package websocket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// Hub 维护活跃连接和广播消息
type Hub struct {
	// 注册的客户端
	clients map[*Client]bool

	// 用户ID到客户端的映射
	userClients map[uint]*Client

	// 注册客户端
	register chan *Client

	// 注销客户端
	unregister chan *Client

	// 广播消息
	broadcast chan []byte

	// 发送给特定用户的消息
	sendToUser chan *UserMessage

	// 互斥锁
	mutex sync.RWMutex

	// 日志
	logger *zap.Logger
}

// Client 表示一个WebSocket连接
type Client struct {
	// WebSocket连接
	conn *websocket.Conn

	// 发送消息的缓冲通道
	send chan []byte

	// 用户ID
	userID uint

	// Hub引用
	hub *Hub
}

// UserMessage 发送给特定用户的消息
type UserMessage struct {
	UserID  uint
	Message []byte
}

// NewHub 创建新的Hub
func NewHub() *Hub {
	logger, _ := zap.NewProduction()
	return &Hub{
		clients:     make(map[*Client]bool),
		userClients: make(map[uint]*Client),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcast:   make(chan []byte),
		sendToUser:  make(chan *UserMessage),
		logger:      logger,
	}
}

// Run 启动Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)

		case userMessage := <-h.sendToUser:
			h.sendToUserMessage(userMessage)
		}
	}
}

// registerClient 注册客户端
func (h *Hub) registerClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients[client] = true
	h.userClients[client.userID] = client

	h.logger.Info("Client registered",
		zap.Uint("userID", client.userID),
		zap.Int("totalClients", len(h.clients)))
}

// unregisterClient 注销客户端
func (h *Hub) unregisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		delete(h.userClients, client.userID)
		close(client.send)
	}

	h.logger.Info("Client unregistered",
		zap.Uint("userID", client.userID),
		zap.Int("totalClients", len(h.clients)))
}

// broadcastMessage 广播消息给所有客户端
func (h *Hub) broadcastMessage(message []byte) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

// sendToUserMessage 发送消息给特定用户
func (h *Hub) sendToUserMessage(userMessage *UserMessage) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if client, ok := h.userClients[userMessage.UserID]; ok {
		select {
		case client.send <- userMessage.Message:
		default:
			close(client.send)
			delete(h.clients, client)
			delete(h.userClients, userMessage.UserID)
		}
	}
}

// SendToUser 发送消息给特定用户
func (h *Hub) SendToUser(userID uint, message []byte) {
	select {
	case h.sendToUser <- &UserMessage{
		UserID:  userID,
		Message: message,
	}:
	default:
		h.logger.Warn("Failed to send message to user",
			zap.Uint("userID", userID))
	}
}

// GetOnlineUsers 获取在线用户数量
func (h *Hub) GetOnlineUsers() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.clients)
}

// IsUserOnline 检查用户是否在线
func (h *Hub) IsUserOnline(userID uint) bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	_, ok := h.userClients[userID]
	return ok
}

// GetOnlineUserIDs 获取在线用户ID列表
func (h *Hub) GetOnlineUserIDs() []uint {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	userIDs := make([]uint, 0, len(h.userClients))
	for userID := range h.userClients {
		userIDs = append(userIDs, userID)
	}
	return userIDs
}

// Client的读写方法
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// 这里可以处理接收到的消息
		_ = message
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("error: %v", err)
				return
			}
		}
	}
}
