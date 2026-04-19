package jwtManager

import (
	"testing"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	uuid "github.com/satori/go.uuid"
)

func TestManager(t *testing.T) {
	kp := &KeysProviderMock{}

	manager := New(kp)

	tokenInfo := entity.TokenInfo{
		UserUid:    uuid.FromStringOrNil("c61f4445-02d5-4afd-8ad5-49b26a3834c1"),
		DeviceUid: uuid.FromStringOrNil("c61f4445-02d5-4afd-8ad5-49b26a3834c1"),
		Email:      "hui",
		Roles:      []string{"admin"},
		SessionUid: uuid.FromStringOrNil("c61f4445-02d5-4afd-8ad5-49b26a3834c1"),
		ExpiredAt:  time.Now().Add(100 * time.Hour).Unix(),
	}

	token, err := manager.NewToken(tokenInfo, 100*time.Hour)
	if err != nil {
		t.Fatalf("failed to create new token: %v", err)
		return
	}

	t.Log(token)
}

type KeysProviderMock struct{}

func (kpm *KeysProviderMock) GetPrivateKey() ([]byte, error) {
	return []byte("huihuihui"), nil
}
