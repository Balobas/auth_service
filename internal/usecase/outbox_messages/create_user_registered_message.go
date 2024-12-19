package useCaseOutboxMessages

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type UserCreatedPayload struct {
	Uid   uuid.UUID `json:"uid"`
	Email string    `json:"email"`
}

func (uc *UseCaseOutboxMessages) CreateUserRegisteredMessage(ctx context.Context, user entity.User) error {
	log.Printf("UseCaseOutboxMessages.CreateUserRegisteredMessage: userUid %s", user.Uid)

	if !uc.cfg.EnableMqMessages() {
		log.Printf("mq messages disabled")
		return nil
	}

	if uuid.Equal(uuid.UUID{}, user.Uid) {
		log.Printf("error: empty user uid")
		return errors.New("empty user uid")
	}

	bts, err := json.Marshal(UserCreatedPayload{Uid: user.Uid, Email: user.Email})
	if err != nil {
		log.Printf("failed to marshal payload for userRegisteredMessage: %v", err)
		return errors.WithStack(err)
	}

	userRegisteredMessage := entity.MqMessage{
		Uid:         uuid.NewV4(),
		SubjectName: uc.cfg.UserRegisteredMessageSubject(),
		Payload:     bts,
		CreatedAt:   time.Now().UTC(),
	}

	if err := uc.outboxRepository.CreateMessage(ctx, userRegisteredMessage); err != nil {
		log.Printf("failed to create user registered message in outbox repository: %v", err)
		return errors.WithStack(err)
	}

	return nil
}
