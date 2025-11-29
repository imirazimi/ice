package service

import (
	"context"
	"encoding/json"
	"ice/internal/outbox"
	"ice/internal/port"
	"ice/pkg/logger"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	repo      port.OutboxRepository
	publisher port.RedisStreamPublisher
}

func NewService(repo port.OutboxRepository, pub port.RedisStreamPublisher) *Service {
	return &Service{repo: repo, publisher: pub}
}

func (s *Service) Write(ctx context.Context, topic string, event any) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return s.repo.Insert(ctx, &outbox.OutboxItem{
		Topic:   topic,
		Payload: string(body),
	})
}

func (s *Service) StartProcessor(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				logger.Get().Info("Outbox processor stopped")
				return

			case <-ticker.C:
				s.process(ctx)
			}
		}
	}()
}

func (s *Service) process(ctx context.Context) {
	msgs, err := s.repo.FetchPending(ctx, 30)
	if err != nil {
		logger.Get().Error("failed to fetch pending outbox", zap.Error(err))
		return
	}

	for _, msg := range msgs {
		var data any
		json.Unmarshal([]byte(msg.Payload), &data)

		if err := s.publisher.Publish(ctx, msg.Topic, data); err != nil {
			logger.Get().Error("failed to publish", zap.Error(err))
			s.repo.MarkFailed(ctx, msg.ID)
			continue
		}

		s.repo.MarkSent(ctx, msg.ID)
	}
}
