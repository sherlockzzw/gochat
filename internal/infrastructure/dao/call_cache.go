package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"gochat/internal/infrastructure/models"
	globalUtils "gochat/utils"
	"time"
)

// CallCache 通话相关缓存操作
type CallCache struct {
	redisClient *globalUtils.RedisService
	ttl         time.Duration // 缓存过期时间
}

// NewCallCache 创建通话缓存
func NewCallCache() *CallCache {
	ttl := 5 * time.Minute // 默认5分钟过期

	return &CallCache{
		redisClient: &globalUtils.RedisService{Client: globalUtils.RDB},
		ttl:         ttl,
	}
}

// GetRoomFromCache 从缓存获取房间信息
func (c *CallCache) GetRoomFromCache(roomID int64) (*models.CallRoom, error) {
	if c.redisClient.Client == nil {
		return nil, fmt.Errorf("Redis not available")
	}

	key := fmt.Sprintf("call:room:%d", roomID)
	ctx := context.Background()

	val, err := c.redisClient.Client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var room models.CallRoom
	if err := json.Unmarshal([]byte(val), &room); err != nil {
		return nil, err
	}

	return &room, nil
}

// SetRoomToCache 设置房间信息到缓存
func (c *CallCache) SetRoomToCache(room *models.CallRoom) error {
	if c.redisClient.Client == nil {
		return nil // Redis不可用时静默失败，降级到数据库
	}

	key := fmt.Sprintf("call:room:%d", room.ID)
	ctx := context.Background()

	data, err := json.Marshal(room)
	if err != nil {
		return err
	}

	return c.redisClient.Client.Set(ctx, key, data, c.ttl).Err()
}

// DeleteRoomFromCache 从缓存删除房间信息
func (c *CallCache) DeleteRoomFromCache(roomID int64) error {
	if c.redisClient.Client == nil {
		return nil
	}

	key := fmt.Sprintf("call:room:%d", roomID)
	ctx := context.Background()

	return c.redisClient.Client.Del(ctx, key).Err()
}

// GetParticipantsFromCache 从缓存获取参与者列表
func (c *CallCache) GetParticipantsFromCache(roomID int64) ([]*models.CallParticipant, error) {
	if c.redisClient.Client == nil {
		return nil, fmt.Errorf("Redis not available")
	}

	key := fmt.Sprintf("call:participants:%d", roomID)
	ctx := context.Background()

	val, err := c.redisClient.Client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var participants []*models.CallParticipant
	if err := json.Unmarshal([]byte(val), &participants); err != nil {
		return nil, err
	}

	return participants, nil
}

// SetParticipantsToCache 设置参与者列表到缓存
func (c *CallCache) SetParticipantsToCache(roomID int64, participants []*models.CallParticipant) error {
	if c.redisClient.Client == nil {
		return nil
	}

	key := fmt.Sprintf("call:participants:%d", roomID)
	ctx := context.Background()

	data, err := json.Marshal(participants)
	if err != nil {
		return err
	}

	return c.redisClient.Client.Set(ctx, key, data, c.ttl).Err()
}

// DeleteParticipantsFromCache 从缓存删除参与者列表
func (c *CallCache) DeleteParticipantsFromCache(roomID int64) error {
	if c.redisClient.Client == nil {
		return nil
	}

	key := fmt.Sprintf("call:participants:%d", roomID)
	ctx := context.Background()

	return c.redisClient.Client.Del(ctx, key).Err()
}

// InvalidateRoomCache 使房间相关缓存失效
func (c *CallCache) InvalidateRoomCache(roomID int64) {
	// 异步删除，不阻塞
	go func() {
		c.DeleteRoomFromCache(roomID)
		c.DeleteParticipantsFromCache(roomID)
	}()
}

