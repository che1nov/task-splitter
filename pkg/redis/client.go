package redis

import (
	"context"
	"fmt"
	"task-splitter/config"

	"github.com/redis/go-redis/v9"
)

// NewClient создает новый клиент Redis
func NewClient(cfg config.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Проверяем подключение
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("⚠️  Failed to connect to Redis: %v\n", err)
		return nil
	}

	fmt.Println("✅ Redis connected successfully")
	return rdb
}
