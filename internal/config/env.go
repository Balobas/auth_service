package config

import (
	"os"
	"strconv"
)

type configEnv struct {
	ServiceName string

	SenderEmail    string
	SenderPassword string
	HostSMTP       string
	PortSMTP       string

	EnableMqMessages bool
	NatsUrl          string
	NatsClientName   string

	UsersStreamName              string
	UserRegisteredMessageSubject string
	UserDeletedMessageSubject    string

	TelegramBotToken string

	MetricsAddr string
}

func ParseEnv(cfg *configEnv) {
	cfg.ServiceName = os.Getenv("SERVICE_NAME")
	cfg.SenderEmail = os.Getenv("SENDER_EMAIL")
	cfg.SenderPassword = os.Getenv("SENDER_PASSWORD")
	cfg.HostSMTP = os.Getenv("HOST_SMTP")
	cfg.PortSMTP = os.Getenv("PORT_SMTP")

	cfg.EnableMqMessages, _ = strconv.ParseBool(os.Getenv("ENABLE_MQ_MESSAGES"))
	cfg.NatsUrl = os.Getenv("NATS_URL")
	cfg.NatsClientName = os.Getenv("NATS_CLIENT_NAME")
	cfg.UserRegisteredMessageSubject = os.Getenv("USER_REGISTERED_MESSAGE_SUBJECT")
	cfg.UserDeletedMessageSubject = os.Getenv("USER_DELETED_MESSAGE_SUBJECT")
	cfg.UsersStreamName = os.Getenv("USERS_STREAM_NAME")
	cfg.TelegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")

	cfg.MetricsAddr = os.Getenv("METRICS_ADDR")
	if cfg.MetricsAddr == "" {
		cfg.MetricsAddr = ":9091" // /metrics по умолчанию
	}
}
