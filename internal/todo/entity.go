package todo

import (
	"time"
)

// TodoItem is the core domain entity for a todo item
// Contains UUID, description, and due date
type TodoItem struct {
	ID          string    // UUID
	Description string    // Description
	DueDate     time.Time // Due date
}
