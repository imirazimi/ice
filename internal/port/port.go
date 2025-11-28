package port

import (
	"context"
	"ice/internal/todo"
)

// Repository abstracts persisting and retrieving todo items
type Repository interface {
	Create(ctx context.Context, item *todo.TodoItem) error
}

// TodoService abstracts the service for todo business logic
type TodoService interface {
	CreateTodo(ctx context.Context, item *todo.TodoItem) error
}

// RedisStreamPublisher abstracts publishing todo items to a Redis Stream
type RedisStreamPublisher interface {
	PublishTodo(ctx context.Context, item *todo.TodoItem) error
}
