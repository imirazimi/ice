package repository

import (
	"context"
	"ice/internal/todo"
)

func (r *Repository) Create(ctx context.Context, item *todo.TodoItem) error {
	_, err := r.mysql.DB().ExecContext(ctx,
		"INSERT INTO todos (id, description, due_date) VALUES (?, ?, ?)",
		item.ID, item.Description, item.DueDate,
	)
	return err
}
