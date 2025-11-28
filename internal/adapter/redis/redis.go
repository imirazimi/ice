package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"ice/config"
	"ice/internal/todo"

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

func (r *RedisStreamClient) PublishTodo(ctx context.Context, item *todo.TodoItem) error {
	// Marshal todo item to JSON
	data, err := json.Marshal(map[string]interface{}{
		"id":          item.ID,
		"description": item.Description,
		"dueDate":     item.DueDate.Format("2006-01-02T15:04:05Z07:00"),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal todo item: %w", err)
	}

	// Add to Redis Stream
	args := &redis.XAddArgs{
		Stream: r.stream,
		Values: map[string]interface{}{
			"data": string(data),
		},
	}

	if err := r.client.XAdd(ctx, args).Err(); err != nil {
		return fmt.Errorf("failed to publish to redis stream: %w", err)
	}

	return nil
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
