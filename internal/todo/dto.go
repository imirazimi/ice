package todo

import "time"

type CreateTodoRequest struct {
	Description string    `json:"description" validate:"required,min=1"`
	DueDate     time.Time `json:"dueDate" validate:"required"`
}

type CreateTodoResponse struct {
	TodoItem TodoItem `json:"todoItem"`
}
