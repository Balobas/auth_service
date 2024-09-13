package email

import (
	"context"
	"net/smtp"
)

type Config interface {
	SenderEmail() string
	SenderPassword() string
	HostSMTP() string
	PortSMTP() string
}

type SmtpClient struct {
	cfg  Config
	auth smtp.Auth
}

func NewClient(ctx context.Context, cfg Config) *SmtpClient {
	return &SmtpClient{
		cfg:  cfg,
		auth: smtp.PlainAuth("", cfg.SenderEmail(), cfg.SenderPassword(), cfg.HostSMTP()),
	}
}

func (c *SmtpClient) SendEmail(receiverEmail string, body []byte) error {
	return smtp.SendMail(
		c.cfg.HostSMTP()+":"+c.cfg.PortSMTP(),
		c.auth,
		c.cfg.SenderEmail(),
		[]string{receiverEmail},
		body,
	)
}
