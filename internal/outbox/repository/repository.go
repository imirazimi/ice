package repository

import (
	"context"
	"ice/internal/adapter/mysql"
	"ice/internal/outbox"
)

type Repository struct {
	db *mysql.MySQL
}

func NewRepository(db *mysql.MySQL) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Insert(ctx context.Context, msg *outbox.OutboxItem) error {
	_, err := r.db.DB().ExecContext(ctx,
		`INSERT INTO outbox (topic, payload, status)
		 VALUES (?, ?, 'pending')`,
		msg.Topic, msg.Payload,
	)
	return err
}

func (r *Repository) FetchPending(ctx context.Context, limit int) ([]outbox.OutboxItem, error) {
	rows, err := r.db.DB().QueryContext(ctx,
		`SELECT id, topic, payload FROM outbox
		 WHERE status='pending'
		 ORDER BY id ASC
		 LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []outbox.OutboxItem
	for rows.Next() {
		var m outbox.OutboxItem
		rows.Scan(&m.ID, &m.Topic, &m.Payload)
		list = append(list, m)
	}
	return list, nil
}

func (r *Repository) MarkSent(ctx context.Context, id int64) error {
	_, err := r.db.DB().ExecContext(ctx,
		`UPDATE outbox SET status='sent' WHERE id=?`,
		id)
	return err
}

func (r *Repository) MarkFailed(ctx context.Context, id int64) error {
	_, err := r.db.DB().ExecContext(ctx,
		`UPDATE outbox SET status='failed' WHERE id=?`,
		id)
	return err
}
