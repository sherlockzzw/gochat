package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoDBService struct {
	Client   *mongo.Client
	Database *mongo.Database
	Logger   *zap.Logger
}

var MongoDB *MongoDBService

// NewMongoDBService 创建MongoDB服务
func NewMongoDBService(uri, database string) (*MongoDBService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// 测试连接
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	logger, _ := zap.NewProduction()

	service := &MongoDBService{
		Client:   client,
		Database: client.Database(database),
		Logger:   logger,
	}

	MongoDB = service
	return service, nil
}

// GetCollection 获取集合
func (m *MongoDBService) GetCollection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}

// Close 关闭连接
func (m *MongoDBService) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.Client.Disconnect(ctx)
}

// Ping 测试连接
func (m *MongoDBService) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.Client.Ping(ctx, nil)
}
