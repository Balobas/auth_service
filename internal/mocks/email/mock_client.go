package emailMock

import (
	"context"
	"fmt"
	"log"
)

type Config interface {
	SenderEmail() string
	SenderPassword() string
	HostSMTP() string
	PortSMTP() string
}

type EmailClientMock struct {
	cfg Config
}

func NewClient(ctx context.Context, cfg Config) *EmailClientMock {
	return &EmailClientMock{
		cfg: cfg,
	}
}

func (c *EmailClientMock) SendEmail(receiverEmail string, body []byte) error {
	msg := fmt.Sprintf(`
		EmailClientMock.SendEmail

		From: %s
		To: %s
		To host %s:%s

		Body: %s

	`,
		c.cfg.SenderEmail(), receiverEmail,
		c.cfg.HostSMTP(), c.cfg.PortSMTP(),
		string(body),
	)
	log.Println(msg)

	return nil
}
