package outbox

import "time"

type OutboxItem struct {
	ID        int64
	Topic     string
	Payload   string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
