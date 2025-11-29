package service

import (
	"context"
	"ice/internal/port"
	"ice/internal/todo"
)

type Service struct {
	repo   port.TodoRepository
	outbox port.OutboxWriter
}

func NewService(repo port.TodoRepository, outbox port.OutboxWriter) *Service {
	return &Service{repo: repo, outbox: outbox}
}

func (s *Service) CreateTodo(ctx context.Context, item *todo.TodoItem) error {
	if err := s.repo.Create(ctx, item); err != nil {
		return err
	}
	return s.outbox.Write(ctx, "todo_stream", item)
}
