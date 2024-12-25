package useCaseOutboxMessages

import uuid "github.com/satori/go.uuid"

type BaseMsgPayload struct {
	MsgUid uuid.UUID `json:"msg_uid"`
}

type UserCreatedPayload struct {
	BaseMsgPayload
	Uid   uuid.UUID `json:"uid"`
	Email string    `json:"email"`
}

type UserDeletedPayload struct {
	BaseMsgPayload
	Uid   uuid.UUID `json:"uid"`
	Email string    `json:"email"`
}
