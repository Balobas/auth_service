package workerVerification

import (
	"bytes"
	"context"
	"log"
	"text/template"
	"time"

	"github.com/balobas/auth_service/internal/entity"
)

type Worker struct {
	cfg Config

	verificationRepo VerificationRepository
	emailNotifier    EmailNotifier
}

func New(
	cfg Config,
	verificationRepo VerificationRepository,
	emailNotifier EmailNotifier,
) *Worker {
	return &Worker{
		cfg:              cfg,
		verificationRepo: verificationRepo,
		emailNotifier:    emailNotifier,
	}
}

func (w *Worker) Run(ctx context.Context) {
	timer := time.NewTimer(w.cfg.SendVerificationInterval())
	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			log.Printf("ctx done %v\n", ctx.Err())
			return
		case <-timer.C:

			verifications, err := w.verificationRepo.GetVerificationsInStatus(
				ctx, entity.VerificationStatusCreated, w.cfg.VerificationWorkerBatchSize(),
			)
			if err != nil {
				log.Printf("error get verifications %v\n", err)
				timer.Reset(w.cfg.SendVerificationInterval())
				break
			}
			templ := w.cfg.EmailVerificationTemplate()
			for _, verification := range verifications {

				body, err := buildEmailBody(templ, verification.Token)
				if err != nil {
					log.Printf("failed to build email body %v\n", err)
					continue
				}

				if err := w.emailNotifier.SendEmail(verification.Email, body); err != nil {
					log.Printf("failed to send message to %s. %v\n", verification.Email, err)
					continue
				}
				verification.Status = entity.VerificationStatusWaiting
				verification.UpdatedAt = time.Now()

				if err := w.verificationRepo.UpdateVerification(ctx, verification); err != nil {
					log.Printf("failed to update verification (email: %s) %v", verification.Email, err)
					continue
				}
			}

			timer.Reset(w.cfg.SendVerificationInterval())
		}
	}
}

func buildEmailBody(tmpl string, token string) ([]byte, error) {
	t, err := template.New("msg").Parse(tmpl)
	if err != nil {
		return nil, err
	}
	buff := bytes.NewBuffer(make([]byte, 0, len(tmpl)+len(token)))

	if err := t.Execute(buff, struct{ Token string }{Token: token}); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
