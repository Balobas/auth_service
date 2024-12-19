package config

import (
	"os"
	"strconv"
)

type configEnv struct {
	SenderEmail    string
	SenderPassword string
	HostSMTP       string
	PortSMTP       string

	EnableMqMessages bool
	NatsUrl          string
	NatsClientName   string

	UserRegisteredMessageSubject string
	UserDeletedMessageSubject    string
}

func ParseEnv(cfg *configEnv) {
	cfg.SenderEmail = os.Getenv("SENDER_EMAIL")
	cfg.SenderPassword = os.Getenv("SENDER_PASSWORD")
	cfg.HostSMTP = os.Getenv("HOST_SMTP")
	cfg.PortSMTP = os.Getenv("PORT_SMTP")

	cfg.EnableMqMessages, _ = strconv.ParseBool(os.Getenv("ENABLE_MQ_MESSAGES"))
	cfg.NatsUrl = os.Getenv("NATS_URL")
	cfg.NatsClientName = os.Getenv("NATS_CLIENT_NAME")
	cfg.UserRegisteredMessageSubject = os.Getenv("USER_REGISTERED_MESSAGE_SUBJECT")
	cfg.UserDeletedMessageSubject = os.Getenv("USER_DELETED_MESSAGE_SUBJECT")
}
