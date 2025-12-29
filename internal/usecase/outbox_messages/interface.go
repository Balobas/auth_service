package useCaseOutboxMessages

import (
	"context"

	outboxEntity "github.com/balobas/sport_city_common/entity/outbox"
)

type Config interface {
	UserRegisteredMessageSubject() string
	UserDeletedMessageSubject() string
	EnableMqMessages() bool
	UsersStreamName() string
}

type RiverClient interface {
	CreateTaskSendOutboxMessage(ctx context.Context, message outboxEntity.Message) error
}
