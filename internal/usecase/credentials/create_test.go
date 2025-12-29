package useCaseCredentials

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestGenPwd(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("events_service@events_service"), bcrypt.DefaultCost)
	t.Log(string(passwordHash))
	t.Log(err)
}
