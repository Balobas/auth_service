package outboxRepository

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
)

func (r *OutboxRepository) UpdateMessage(ctx context.Context, message entity.MqMessage) error {
	log.Printf("outboxRepository.UpdateMessage: %s %s", message.Uid, message.SubjectName)

	msgRow := pgEntity.NewMqMessageRow().FromEntity(message)

	if err := r.Update(ctx, msgRow, msgRow.ConditionUidEqual()); err != nil {
		log.Printf("failed to update outbox message %s: %v", message.Uid, err)
		return errors.WithStack(err)
	}
	return nil
}
