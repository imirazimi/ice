package todo

import "time"

type CreateTodoRequest struct {
	Description string    `json:"description" validate:"required,min=1" example:"test task"`
	DueDate     time.Time `json:"dueDate" validate:"required" example:"2025-01-01T06:00:00Z"`
}

type CreateTodoResponse struct {
	TodoItem TodoItem `json:"todoItem"`
}
