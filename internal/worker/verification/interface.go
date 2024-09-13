package workerVerification

import (
	"context"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	uuid "github.com/satori/go.uuid"
)

type Config interface {
	SendVerificationInterval() time.Duration
	VerificationWorkerBatchSize() uint64
	EmailVerificationTemplate() string
}

type VerificationRepository interface {
	UpdateVerification(ctx context.Context, verification entity.Verification) error
	GetUserVerification(ctx context.Context, userUid uuid.UUID) (entity.Verification, error)
	GetVerificationsInStatus(ctx context.Context, status entity.VerificationStatus, limit uint64) ([]entity.Verification, error)
}

type EmailNotifier interface {
	SendEmail(receiverEmail string, body []byte) error
}
