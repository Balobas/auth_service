package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	sessionRepository "github.com/balobas/auth_service/internal/repository/postgres/session"
	repositoryUsers "github.com/balobas/auth_service/internal/repository/postgres/users"
	uuid "github.com/satori/go.uuid"
)

func TestCascade(t *testing.T) {
	ctx := context.Background()

	c := NewPgClient(t, ctx)

	usersRepo := repositoryUsers.New(c)
	user := entity.User{
		Uid:       uuid.NewV4(),
		Email:     "huisjkjsk@test",
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := usersRepo.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("failed to create user for test. err: %v", err)
	}

	sesRepo := sessionRepository.New(c)
	sUid := uuid.NewV4()
	err = sesRepo.CreateSession(ctx, entity.Session{
		Uid:       sUid,
		UserUid:   user.Uid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	if err := usersRepo.DeleteUser(ctx, user.Uid); err != nil {
		t.Fatalf(err.Error())
	}

	_, isFound, err := sesRepo.GetSessionByUid(ctx, sUid)
	if err != nil {
		t.Fatalf("failed to get session after delete by user uid: %v", err)
	}
	if isFound {
		t.Fatalf("session found after delete by user uid")
	}
}
