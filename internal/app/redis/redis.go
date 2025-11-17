package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"lab31/internal/app/config"

	"github.com/go-redis/redis/v8"
)

const servicePrefix = "lab31."

type Client struct {
	client *redis.Client
}

func New(ctx context.Context, cfg config.RedisConfig) (*Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Password:    cfg.Password,
		Username:    cfg.User,
		Addr:        fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DB:          0,
		DialTimeout: cfg.DialTimeout,
		ReadTimeout: cfg.ReadTimeout,
	})

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("cant ping redis: %w", err)
	}
	return &Client{client: redisClient}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

// AddToBlacklist добавляет JWT в черный список с TTL
func (c *Client) AddToBlacklist(ctx context.Context, token string, ttl time.Duration) error {
	key := servicePrefix + "blacklist:" + token
	// Устанавливаем ключ со значением "true" и временем жизни (TTL)
	return c.client.Set(ctx, key, true, ttl).Err()
}

// IsBlacklisted проверяет, находится ли JWT в черном списке
func (c *Client) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := servicePrefix + "blacklist:" + token
	_, err := c.client.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return false, nil // Не в черном списке (NotFound)
	}
	if err != nil {
		return false, err // Внутренняя ошибка redis
	}
	return true, nil // Найдено, в черном списке
}

func (c *Client) GetBlacklistEntry(ctx context.Context, token string) *redis.StringCmd {
	key := servicePrefix + "blacklist:" + token
	return c.client.Get(ctx, key)
}
