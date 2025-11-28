package redis

import (
	"context"
	"ice/config"
	"ice/internal/todo"
	"testing"
	"time"
)

func TestRedisStreamClient_PublishTodo(t *testing.T) {
	// Skip if redis is not available
	cfg := config.RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
	client, err := NewRedisStreamClient(cfg)
	if err != nil {
		t.Skip("Redis not available, skipping test")
		return
	}
	defer client.Close()

	item := &todo.TodoItem{
		ID:          "test-id",
		Description: "test description",
		DueDate:     time.Now(),
	}
	err = client.PublishTodo(context.Background(), item)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}
