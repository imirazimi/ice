package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"ice/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStreamClient struct {
	client *redis.Client
	stream string
}

func NewRedisStreamClient(cfg config.RedisConfig) (*RedisStreamClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RedisStreamClient{
		client: client,
		stream: "todos",
	}, nil
}

func (r *RedisStreamClient) Publish(ctx context.Context, stream string, data interface{}) error {
	// تبدیل data به JSON
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// ارسال به Redis Stream
	args := &redis.XAddArgs{
		Stream: stream,
		Values: map[string]interface{}{
			"payload": string(payload),
		},
		MaxLen: 0, // می‌توانید محدودیت طول Stream را تنظیم کنید
	}

	// با timeout کوتاه برای جلوگیری از بلاک شدن
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err = r.client.XAdd(ctx, args).Result()
	return err
}

func (r *RedisStreamClient) Client() *redis.Client {
	return r.client
}

func (r *RedisStreamClient) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}
