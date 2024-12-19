package outboxRepository

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
)

func (r *OutboxRepository) CreateMessage(ctx context.Context, message entity.MqMessage) error {
	log.Printf("outboxRepository.CreateMessage: %s %s", message.Uid, message.SubjectName)

	if err := r.Create(ctx, pgEntity.NewMqMessageRow().FromEntity(message)); err != nil {
		log.Printf("failed to create outbox message %s: %v", message.Uid, err)
		return errors.WithStack(err)
	}
	return nil
}
