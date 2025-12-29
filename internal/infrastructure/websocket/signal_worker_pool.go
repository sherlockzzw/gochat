package websocket

import (
	"context"
	"log"
	"sync"
	"time"
)

// SignalWorkerPool 信令处理工作池
// 使用goroutine池处理信令消息，避免阻塞WebSocket主循环
type SignalWorkerPool struct {
	workers    int                 // 工作协程数量
	jobQueue   chan *SignalJob     // 任务队列
	wg         sync.WaitGroup      // 等待组
	ctx        context.Context     // 上下文
	cancel     context.CancelFunc  // 取消函数
	processor  CallSignalProcessor // 信令处理器
	maxRetries int                 // 最大重试次数
}

// SignalJob 信令任务
type SignalJob struct {
	Message    []byte    // 信令消息
	UserID     int64     // 用户ID
	RetryCount int       // 重试次数
	CreatedAt  time.Time // 创建时间
}

// NewSignalWorkerPool 创建信令处理工作池
func NewSignalWorkerPool(workers int, queueSize int, processor CallSignalProcessor) *SignalWorkerPool {
	if workers <= 0 {
		workers = 10 // 默认10个worker
	}
	if queueSize <= 0 {
		queueSize = 1000 // 默认队列大小1000
	}

	ctx, cancel := context.WithCancel(context.Background())

	pool := &SignalWorkerPool{
		workers:    workers,
		jobQueue:   make(chan *SignalJob, queueSize),
		ctx:        ctx,
		cancel:     cancel,
		processor:  processor,
		maxRetries: 3,
	}

	// 启动worker
	pool.start()

	return pool
}

// start 启动工作池
func (p *SignalWorkerPool) start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

// worker 工作协程
func (p *SignalWorkerPool) worker(id int) {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return
		case job := <-p.jobQueue:
			p.processJob(job, id)
		}
	}
}

// processJob 处理任务
func (p *SignalWorkerPool) processJob(job *SignalJob, workerID int) {
	// 检查任务是否超时（超过5秒的任务丢弃）
	if time.Since(job.CreatedAt) > 5*time.Second {
		log.Printf("Worker %d: Job timeout, dropped", workerID)
		return
	}

	// 处理信令
	err := p.processor.ProcessCallSignal(job.Message, job.UserID)
	if err != nil {
		log.Printf("Worker %d: Failed to process signal: %v", workerID, err)

		// 重试机制
		if job.RetryCount < p.maxRetries {
			job.RetryCount++
			// 指数退避：1s, 2s, 4s
			backoff := time.Duration(1<<uint(job.RetryCount-1)) * time.Second
			time.Sleep(backoff)

			// 重新加入队列
			select {
			case p.jobQueue <- job:
			default:
				log.Printf("Worker %d: Job queue full, dropped retry", workerID)
			}
		} else {
			log.Printf("Worker %d: Max retries reached, dropped job", workerID)
		}
	}
}

// Submit 提交任务
func (p *SignalWorkerPool) Submit(message []byte, userID int64) bool {
	job := &SignalJob{
		Message:    message,
		UserID:     userID,
		RetryCount: 0,
		CreatedAt:  time.Now(),
	}

	select {
	case p.jobQueue <- job:
		return true
	default:
		// 队列满，记录日志但不阻塞
		log.Printf("SignalWorkerPool: Queue full, dropped signal from user %d", userID)
		return false
	}
}

// Stop 停止工作池
func (p *SignalWorkerPool) Stop() {
	p.cancel()
	close(p.jobQueue)
	p.wg.Wait()
}

// Stats 获取统计信息
func (p *SignalWorkerPool) Stats() map[string]interface{} {
	return map[string]interface{}{
		"workers":    p.workers,
		"queue_size": len(p.jobQueue),
		"queue_cap":  cap(p.jobQueue),
	}
}

