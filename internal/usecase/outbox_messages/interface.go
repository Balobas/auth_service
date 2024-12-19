package useCaseOutboxMessages

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
)

type Config interface {
	UserRegisteredMessageSubject() string
	UserDeletedMessageSubject() string
	EnableMqMessages() bool
}

type OutboxRepository interface {
	CreateMessage(ctx context.Context, message entity.MqMessage) error
}
