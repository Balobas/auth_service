package workerPublisher

import (
	"context"
	"time"

	"github.com/balobas/auth_service/internal/entity"
)

type Config interface {
	MqPublishMessagesInterval() time.Duration
	MqPublishMessagesBatchSize() int64
}

type Publisher interface {
	Publish(ctx context.Context, subjectName string, data []byte) error
}

type OutboxRepository interface {
	GetReadyMessagesForPublish(ctx context.Context, batchSize int64) ([]entity.MqMessage, error)
	UpdateMessage(ctx context.Context, msg entity.MqMessage) error
}
