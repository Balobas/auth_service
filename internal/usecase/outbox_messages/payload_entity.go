package useCaseOutboxMessages

import uuid "github.com/satori/go.uuid"

type BaseMsgPayload struct {
	MsgUid uuid.UUID `json:"msg_uid"`
}

type UserCreatedPayload struct {
	BaseMsgPayload
	Uid      uuid.UUID `json:"uid"`
	Email    string    `json:"email"`
	Name     string    `json:"name"`     // отображаемое имя (из Telegram first_name), пусто для email-регистрации
	Nickname string    `json:"nickname"` // желаемый ник (из Telegram username), пусто если нет
}

type UserDeletedPayload struct {
	BaseMsgPayload
	Uid   uuid.UUID `json:"uid"`
	Email string    `json:"email"`
}
