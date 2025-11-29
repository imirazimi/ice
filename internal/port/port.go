package port

import (
	"context"
	"ice/internal/outbox"
	"ice/internal/todo"
)

// Repository abstracts persisting and retrieving todo items
type TodoRepository interface {
	Create(ctx context.Context, item *todo.TodoItem) error
}

// TodoService abstracts the service for todo business logic
type TodoService interface {
	CreateTodo(ctx context.Context, item *todo.TodoItem) error
}

// RedisStreamPublisher abstracts publishing todo items to a Redis Stream
type RedisStreamPublisher interface {
	Publish(ctx context.Context, stream string, data interface{}) error
}

type OutboxWriter interface {
	Write(ctx context.Context, topic string, event any) error
}

type OutboxRepository interface {
	Insert(ctx context.Context, msg *outbox.OutboxItem) error
	FetchPending(ctx context.Context, limit int) ([]outbox.OutboxItem, error)
	MarkSent(ctx context.Context, id int64) error
	MarkFailed(ctx context.Context, id int64) error
}
