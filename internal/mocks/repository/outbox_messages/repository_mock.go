package outboxMessagesMockRepository

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
)

type RepositoryOutboxMessagesMock struct {
}

func New() *RepositoryOutboxMessagesMock {
	return &RepositoryOutboxMessagesMock{}
}

func (r *RepositoryOutboxMessagesMock) CreateMessage(ctx context.Context, message entity.MqMessage) error {
	log.Printf("RepositoryOutboxMessagesMock.CreateMessage")
	return nil
}

func (r *RepositoryOutboxMessagesMock) GetReadyMessagesForPublish(ctx context.Context, batchSize int64) ([]entity.MqMessage, error) {
	log.Printf("RepositoryOutboxMessagesMock.GetReadyMessagesForPublish")
	return nil, nil
}

func (r *RepositoryOutboxMessagesMock) UpdateMessage(ctx context.Context, msg entity.MqMessage) error {
	log.Printf("RepositoryOutboxMessagesMock.UpdateMessage")
	return nil
}
