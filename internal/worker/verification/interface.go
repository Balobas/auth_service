package workerVerification

import (
	"context"
	"time"

	"github.com/balobas/auth_service/internal/entity"
)

type Config interface {
	SendVerificationInterval() time.Duration
	VerificationWorkerBatchSize() uint64
	EmailVerificationTemplate() string
	HttpVerificationScheme() string
}

type VerificationRepository interface {
	UpdateVerification(ctx context.Context, verification entity.Verification) error
	GetVerificationsInStatus(ctx context.Context, status entity.VerificationStatus, limit uint64) ([]entity.Verification, error)
}

type EmailNotifier interface {
	SendEmail(receiverEmail string, body []byte) error
}
