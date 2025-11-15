package redis

import (
	"context"
	"task-splitter/pkg/logger"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheAdapter адаптер для работы с Redis
type CacheAdapter struct {
	client *redis.Client
	log    logger.Logger
}

// NewCacheAdapter создает новый адаптер кэша
func NewCacheAdapter(client *redis.Client, log logger.Logger) *CacheAdapter {
	return &CacheAdapter{
		client: client,
		log:    log,
	}
}

// Set устанавливает значение в кэш
func (a *CacheAdapter) Set(ctx context.Context, key string, value interface{}, ttlSeconds int) error {
	ttl := time.Duration(ttlSeconds) * time.Second
	if err := a.client.Set(ctx, key, value, ttl).Err(); err != nil {
		a.log.ErrorContext(ctx, "Ошибка установки значения в кэш", "key", key, "error", err)
		return err
	}

	a.log.DebugContext(ctx, "Значение установлено в кэш", "key", key, "ttl", ttlSeconds)
	return nil
}

// Get получает значение из кэша
func (a *CacheAdapter) Get(ctx context.Context, key string) (string, error) {
	value, err := a.client.Get(ctx, key).Result()
	if err == redis.Nil {
		a.log.DebugContext(ctx, "Значение не найдено в кэше", "key", key)
		return "", nil
	}
	if err != nil {
		a.log.ErrorContext(ctx, "Ошибка получения значения из кэша", "key", key, "error", err)
		return "", err
	}

	a.log.DebugContext(ctx, "Значение получено из кэша", "key", key)
	return value, nil
}

// Delete удаляет значение из кэша
func (a *CacheAdapter) Delete(ctx context.Context, key string) error {
	if err := a.client.Del(ctx, key).Err(); err != nil {
		a.log.ErrorContext(ctx, "Ошибка удаления значения из кэша", "key", key, "error", err)
		return err
	}

	a.log.DebugContext(ctx, "Значение удалено из кэша", "key", key)
	return nil
}

// Exists проверяет существование ключа в кэше
func (a *CacheAdapter) Exists(ctx context.Context, key string) (bool, error) {
	count, err := a.client.Exists(ctx, key).Result()
	if err != nil {
		a.log.ErrorContext(ctx, "Ошибка проверки существования ключа", "key", key, "error", err)
		return false, err
	}

	exists := count > 0
	a.log.DebugContext(ctx, "Проверка существования ключа", "key", key, "exists", exists)
	return exists, nil
}

