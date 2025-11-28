package service

import (
	"context"
	"ice/internal/port"
	"ice/internal/todo"
)

type TodoServiceImpl struct {
	repo  port.Repository
	redis port.RedisStreamPublisher
}

func NewTodoService(repo port.Repository, redis port.RedisStreamPublisher) *TodoServiceImpl {
	return &TodoServiceImpl{repo: repo, redis: redis}
}

// TODO: need to implement Outbox Pattern to ensure data consistency between MySQL and Redis.
func (s *TodoServiceImpl) CreateTodo(ctx context.Context, item *todo.TodoItem) error {
	if err := s.repo.Create(ctx, item); err != nil {
		return err
	}
	if err := s.redis.PublishTodo(ctx, item); err != nil {
		return err
	}
	return nil
}
