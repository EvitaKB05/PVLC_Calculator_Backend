// internal/redis/client.go
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Client - обертка для Redis клиента
// ДОБАВЛЕНО ДЛЯ ЛАБОРАТОРНОЙ РАБОТЫ 4
type Client struct {
	client *redis.Client
}

// NewRedisClient создает нового клиента Redis
func NewRedisClient(host string, port int, password string, db int) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Проверяем подключение
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Client{client: rdb}, nil
}

// SetSession сохраняет сессию пользователя в Redis
func (c *Client) SetSession(ctx context.Context, sessionID string, userData map[string]interface{}, expiration time.Duration) error {
	return c.client.HSet(ctx, "session:"+sessionID, userData).Err()
}

// GetSession получает сессию пользователя из Redis
func (c *Client) GetSession(ctx context.Context, sessionID string) (map[string]string, error) {
	return c.client.HGetAll(ctx, "session:"+sessionID).Result()
}

// DeleteSession удаляет сессию пользователя из Redis
func (c *Client) DeleteSession(ctx context.Context, sessionID string) error {
	return c.client.Del(ctx, "session:"+sessionID).Err()
}

// WriteJWTToBlacklist добавляет JWT токен в черный список
func (c *Client) WriteJWTToBlacklist(ctx context.Context, jwtStr string, jwtTTL time.Duration) error {
	return c.client.Set(ctx, "jwt_blacklist:"+jwtStr, true, jwtTTL).Err()
}

// CheckJWTInBlacklist проверяет, находится ли JWT токен в черном списке
func (c *Client) CheckJWTInBlacklist(ctx context.Context, jwtStr string) (bool, error) {
	result, err := c.client.Exists(ctx, "jwt_blacklist:"+jwtStr).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// Close закрывает соединение с Redis
func (c *Client) Close() error {
	return c.client.Close()
}
