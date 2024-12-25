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

func (uc *UseCaseOutboxMessages) CreateUserDeletedMessage(ctx context.Context, user entity.User) error {
	log.Printf("UseCaseOutboxMessages.CreateUserDeletedMessage: userUid %s", user.Uid)

	if !uc.cfg.EnableMqMessages() {
		log.Printf("mq messages disabled")
		return nil
	}

	if uuid.Equal(uuid.UUID{}, user.Uid) {
		log.Printf("error: empty user uid")
		return errors.New("empty user uid")
	}

	msgUid := uuid.NewV4()

	payload := UserDeletedPayload{Uid: user.Uid, Email: user.Email}
	payload.MsgUid = msgUid

	bts, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal payload for userDeletedMessage: %v", err)
		return errors.WithStack(err)
	}

	userRegisteredMessage := entity.MqMessage{
		Uid:         msgUid,
		SubjectName: uc.cfg.UsersStreamName() + "." + uc.cfg.UserDeletedMessageSubject(),
		Payload:     bts,
		CreatedAt:   time.Now().UTC(),
	}

	if err := uc.outboxRepository.CreateMessage(ctx, userRegisteredMessage); err != nil {
		log.Printf("failed to create user registered message in outbox repository: %v", err)
		return errors.WithStack(err)
	}

	return nil
}
