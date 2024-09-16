package config

import "os"

type configEnv struct {
	SenderEmail    string
	SenderPassword string
	HostSMTP       string
	PortSMTP       string
}

func ParseEnv(cfg *configEnv) {
	cfg.SenderEmail = os.Getenv("SENDER_EMAIL")
	cfg.SenderPassword = os.Getenv("SENDER_PASSWORD")
	cfg.HostSMTP = os.Getenv("HOST_SMTP")
	cfg.PortSMTP = os.Getenv("PORT_SMTP")
}
