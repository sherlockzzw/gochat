package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
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

	// 最大连接数限制，0表示无限制
	maxConnections int

	// 心跳配置
	pingInterval time.Duration // 发送ping的间隔
	pongWait     time.Duration // 等待pong的超时时间
	pingMessage  []byte        // ping消息

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

	// 心跳相关
	lastPongTime   time.Time // 最后一次收到pong的时间
	heartbeatTimer *time.Timer
	pingTicker     *time.Ticker
	mu             sync.Mutex // 保护心跳相关字段
}

// UserMessage 发送给特定用户的消息
type UserMessage struct {
	UserID  uint
	Message []byte
}

// NewHub 创建新的Hub
func NewHub() *Hub {
	logger, _ := zap.NewProduction()

	// 从配置文件读取最大连接数，默认为5000
	maxConnections := viper.GetInt("websocket.max_connections")
	if maxConnections <= 0 {
		maxConnections = 5000 // 默认值
	}

	// 从配置文件读取心跳配置
	pingIntervalSeconds := viper.GetInt("websocket.ping_interval")
	if pingIntervalSeconds <= 0 {
		pingIntervalSeconds = 30 // 默认30秒
	}
	pingInterval := time.Duration(pingIntervalSeconds) * time.Second

	pongWaitSeconds := viper.GetInt("websocket.pong_wait")
	if pongWaitSeconds <= 0 {
		pongWaitSeconds = 60 // 默认60秒（允许2次心跳间隔）
	}
	pongWait := time.Duration(pongWaitSeconds) * time.Second

	// 构造ping消息
	pingMsg := map[string]interface{}{
		"type": "pong", // 后端返回pong响应前端的ping
	}
	pingMessage, _ := json.Marshal(pingMsg)

	logger.Info("WebSocket Hub initialized",
		zap.Int("maxConnections", maxConnections),
		zap.Duration("pingInterval", pingInterval),
		zap.Duration("pongWait", pongWait))

	return &Hub{
		clients:        make(map[*Client]bool),
		userClients:    make(map[uint]*Client),
		register:       make(chan *Client, 1000),
		unregister:     make(chan *Client, 1000),
		broadcast:      make(chan []byte, 1000),
		sendToUser:     make(chan *UserMessage, 10000),
		maxConnections: maxConnections,
		pingInterval:   pingInterval,
		pongWait:       pongWait,
		pingMessage:    pingMessage,
		logger:         logger,
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

	// 检查连接数限制
	currentConnections := len(h.clients)
	if h.maxConnections > 0 && currentConnections >= h.maxConnections {
		// 如果用户已经在线，允许替换旧连接
		if oldClient, exists := h.userClients[client.userID]; exists {
			if oldClient != client {
				fmt.Printf("用户 %d 已有连接，替换旧连接（连接数已达上限: %d）\n",
					client.userID, h.maxConnections)
				delete(h.clients, oldClient)
				close(oldClient.send)
				// 继续注册新连接
			} else {
				// 同一个客户端，直接返回
				return
			}
		} else {
			// 连接数已达上限，拒绝新连接
			fmt.Printf("拒绝新连接：连接数已达上限 %d，当前连接数: %d\n",
				h.maxConnections, currentConnections)
			h.logger.Warn("Connection rejected: max connections reached",
				zap.Int("maxConnections", h.maxConnections),
				zap.Int("currentConnections", currentConnections),
				zap.Uint("userID", client.userID))
			// 关闭连接
			close(client.send)
			client.conn.Close()
			return
		}
	}

	// 如果用户已经在线，先注销旧连接（防止同一用户多个连接）
	if oldClient, exists := h.userClients[client.userID]; exists {
		if oldClient != client { // 避免关闭自己
			fmt.Printf("用户 %d 已有连接，先注销旧连接\n", client.userID)
			delete(h.clients, oldClient)
			close(oldClient.send)
		}
	}

	h.clients[client] = true
	h.userClients[client.userID] = client

	onlineCount := len(h.clients)
	// 只在连接数较少或特定条件下打印详细日志，避免日志过多影响性能
	if onlineCount%100 == 0 || onlineCount < 10 {
		fmt.Printf("用户 %d 已注册WebSocket连接，总连接数: %d",
			client.userID, onlineCount)
		if h.maxConnections > 0 {
			fmt.Printf(" (上限: %d)", h.maxConnections)
		}
		fmt.Println()
	}

	h.logger.Info("Client registered",
		zap.Uint("userID", client.userID),
		zap.Int("totalClients", onlineCount),
		zap.Int("maxConnections", h.maxConnections))
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
// 这个方法在Hub的Run goroutine中执行，确保线程安全
func (h *Hub) sendToUserMessage(userMessage *UserMessage) {
	// 验证用户ID
	if userMessage.UserID == 0 {
		h.logger.Warn("Invalid userID in message")
		return
	}

	h.mutex.RLock()
	client, ok := h.userClients[userMessage.UserID]
	onlineCount := len(h.userClients)
	h.mutex.RUnlock()

	if !ok {
		// 用户不在线，记录日志但不打印所有在线用户（避免性能问题）
		if onlineCount%100 == 0 || onlineCount < 10 {
			// 只在特定条件下打印详细日志，避免日志过多
			fmt.Printf("错误: 用户 %d 不在线，当前在线用户数量: %d\n",
				userMessage.UserID, onlineCount)
		}
		h.logger.Warn("User not found in online clients",
			zap.Uint("userID", userMessage.UserID),
			zap.Int("onlineUsers", onlineCount))
		return
	}

	// 尝试发送消息，使用非阻塞方式
	select {
	case client.send <- userMessage.Message:
		// 消息发送成功
		if onlineCount%100 == 0 || onlineCount < 10 {
			fmt.Printf("消息已成功发送给用户 %d，消息长度: %d\n",
				userMessage.UserID, len(userMessage.Message))
		}
		h.logger.Debug("Message sent to user",
			zap.Uint("userID", userMessage.UserID),
			zap.Int("messageLength", len(userMessage.Message)))
	default:
		// 客户端发送通道已满，说明客户端处理速度慢
		// 可以选择：1. 关闭连接 2. 丢弃消息 3. 等待
		// 这里选择关闭连接，让客户端重连
		fmt.Printf("警告: 用户 %d 的发送通道已满，关闭连接\n", userMessage.UserID)
		h.logger.Warn("Failed to send message to user, channel full",
			zap.Uint("userID", userMessage.UserID))

		// 需要加写锁来删除客户端
		h.mutex.Lock()
		if _, exists := h.clients[client]; exists {
			delete(h.clients, client)
			delete(h.userClients, userMessage.UserID)
			close(client.send)
		}
		h.mutex.Unlock()
	}
}

// SendToUser 发送消息给特定用户
// 这个方法是非阻塞的，如果通道满了会记录警告但不会阻塞调用者
func (h *Hub) SendToUser(userID uint, message []byte) {
	// 验证用户ID有效性
	if userID == 0 {
		h.logger.Warn("Invalid userID: 0")
		return
	}

	// 检查消息大小，防止过大的消息
	if len(message) > 1024*1024 { // 1MB限制
		h.logger.Warn("Message too large",
			zap.Uint("userID", userID),
			zap.Int("messageLength", len(message)))
		return
	}

	select {
	case h.sendToUser <- &UserMessage{
		UserID:  userID,
		Message: message,
	}:
		// 消息已成功入队
	default:
		// 通道满了，记录警告
		// 在高并发场景下，可以考虑使用更大的缓冲或者实现消息丢弃策略
		h.logger.Warn("Failed to queue message for user, channel full",
			zap.Uint("userID", userID),
			zap.Int("messageLength", len(message)),
			zap.Int("channelCapacity", cap(h.sendToUser)))
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
		c.stopHeartbeat()
	}()

	// 设置读超时
	c.conn.SetReadDeadline(time.Now().Add(c.hub.pongWait))
	c.conn.SetPongHandler(func(string) error {
		// 收到pong，更新最后pong时间并延长读超时
		c.mu.Lock()
		c.lastPongTime = time.Now()
		c.mu.Unlock()
		c.conn.SetReadDeadline(time.Now().Add(c.hub.pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// 处理心跳消息
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err == nil {
			if msgType, ok := msg["type"].(string); ok {
				if msgType == "ping" {
					// 前端发送ping，后端回复pong
					c.mu.Lock()
					c.lastPongTime = time.Now()
					c.mu.Unlock()
					c.conn.SetReadDeadline(time.Now().Add(c.hub.pongWait))
					// 发送pong响应
					pongMsg := map[string]interface{}{
						"type": "pong",
					}
					pongData, _ := json.Marshal(pongMsg)
					select {
					case c.send <- pongData:
					default:
						// 发送通道已满，关闭连接
						return
					}
					continue
				}
			}
		}

		// 其他消息可以在这里处理
		_ = message
	}
}

func (c *Client) writePump() {
	// 启动心跳定时器
	c.startHeartbeat()
	defer func() {
		c.stopHeartbeat()
		c.conn.Close()
	}()

	// 设置写超时
	ticker := time.NewTicker(54 * time.Second) // 比ping间隔稍长，用于写超时检测
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("error: %v", err)
				return
			}

		case <-ticker.C:
			// 定期检查写超时
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// startHeartbeat 启动心跳检测
func (c *Client) startHeartbeat() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.lastPongTime = time.Now()

	// 启动心跳超时检测goroutine
	go func() {
		ticker := time.NewTicker(c.hub.pongWait / 2) // 每半个超时时间检查一次
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.mu.Lock()
				lastPong := c.lastPongTime
				c.mu.Unlock()

				// 检查是否超时
				if time.Since(lastPong) > c.hub.pongWait {
					c.hub.logger.Warn("Heartbeat timeout, closing connection",
						zap.Uint("userID", c.userID),
						zap.Duration("timeSinceLastPong", time.Since(lastPong)))
					c.hub.unregister <- c
					c.conn.Close()
					return
				}
			}
		}
	}()
}

// stopHeartbeat 停止心跳检测
func (c *Client) stopHeartbeat() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.pingTicker != nil {
		c.pingTicker.Stop()
		c.pingTicker = nil
	}

	if c.heartbeatTimer != nil {
		c.heartbeatTimer.Stop()
		c.heartbeatTimer = nil
	}
}
