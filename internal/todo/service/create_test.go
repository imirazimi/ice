package service

import (
	"context"
	"errors"
	"ice/internal/todo"
	"testing"
	"time"
)

type mockRepo struct {
	err      error
	received *todo.TodoItem
}

func (m *mockRepo) Create(ctx context.Context, item *todo.TodoItem) error {
	m.received = item
	return m.err
}

type mockRedis struct {
	err      error
	received *todo.TodoItem
}

func (m *mockRedis) PublishTodo(ctx context.Context, item *todo.TodoItem) error {
	m.received = item
	return m.err
}

func TestCreateTodo_Success(t *testing.T) {
	repo := &mockRepo{}
	red := &mockRedis{}
	svc := NewTodoService(repo, red)
	item := &todo.TodoItem{ID: "id1", Description: "desc", DueDate: time.Now()}
	ctx := context.Background()

	err := svc.CreateTodo(ctx, item)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}
	if repo.received != item {
		t.Errorf("repo got different item")
	}
	if red.received != item {
		t.Errorf("redis got different item")
	}
}

func TestCreateTodo_RepoFail(t *testing.T) {
	repo := &mockRepo{err: errors.New("repo fail")}
	red := &mockRedis{}
	svc := NewTodoService(repo, red)
	item := &todo.TodoItem{ID: "id2", Description: "desc2"}
	ctx := context.Background()

	err := svc.CreateTodo(ctx, item)
	if err == nil {
		t.Fatal("want error, got nil")
	}
}

func TestCreateTodo_RedisFail(t *testing.T) {
	repo := &mockRepo{}
	red := &mockRedis{err: errors.New("redis fail")}
	svc := NewTodoService(repo, red)
	item := &todo.TodoItem{ID: "id3", Description: "desc3"}
	ctx := context.Background()

	err := svc.CreateTodo(ctx, item)
	if err == nil {
		t.Fatal("want error, got nil")
	}
}
