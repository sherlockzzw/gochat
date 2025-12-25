package websocket

import (
	"log"
)

// CallSignalProcessor 通话信令处理器接口（由handler层实现）
type CallSignalProcessor interface {
	ProcessCallSignal(signalData []byte, fromUserID int64) error
}

// Hub的CallSignalProcessor，由外部设置
var globalCallSignalProcessor CallSignalProcessor
var globalSignalWorkerPool *SignalWorkerPool

// SetCallSignalProcessor 设置通话信令处理器
func SetCallSignalProcessor(processor CallSignalProcessor) {
	globalCallSignalProcessor = processor
	// 初始化worker pool（10个worker，队列大小1000）
	if globalSignalWorkerPool == nil {
		globalSignalWorkerPool = NewSignalWorkerPool(10, 1000, processor)
	}
}

// GetSignalWorkerPool 获取信令工作池
func GetSignalWorkerPool() *SignalWorkerPool {
	return globalSignalWorkerPool
}

// HandleCallSignal 处理通话信令消息（异步处理）
func (c *Client) HandleCallSignal(message []byte) {
	// 使用worker pool异步处理，不阻塞WebSocket主循环
	if globalSignalWorkerPool != nil {
		// 提交到工作池，如果队列满则丢弃（避免阻塞）
		if !globalSignalWorkerPool.Submit(message, c.userID) {
			log.Printf("Signal worker pool queue full, dropped signal from user %d", c.userID)
		}
		return
	}

	// 如果没有worker pool，降级到同步处理
	if globalCallSignalProcessor != nil {
		if err := globalCallSignalProcessor.ProcessCallSignal(message, c.userID); err != nil {
			log.Printf("Failed to process call signal: %v", err)
		}
		return
	}

	// 如果没有处理器，记录日志
	log.Printf("Call signal received but no processor set: user_id=%d, message=%s",
		c.userID, string(message))
}
