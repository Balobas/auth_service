package repo_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	sessionRepository "github.com/balobas/auth_service/internal/repository/postgres/session"
	repositoryUsers "github.com/balobas/auth_service/internal/repository/postgres/users"
	uuid "github.com/satori/go.uuid"
)

func TestSessionRepo(t *testing.T) {

	ctx := context.Background()

	c := NewPgClient(t, ctx)

	sRepo := sessionRepository.New(c)

	userUid := uuid.NewV4()

	session := entity.Session{
		Uid:       uuid.NewV4(),
		UserUid:   userUid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := sRepo.CreateSession(ctx, session); err == nil {
		t.Fatalf("created session without user")
	}

	usersRepo := repositoryUsers.New(c)

	err := usersRepo.CreateUser(ctx, entity.User{
		Uid:       userUid,
		Email:     fmt.Sprintf("%s@%s", randStr(), randStr()),
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("failed to create user %v", err)
	}
	defer func() {
		// delete test user
		if err := usersRepo.DeleteUser(ctx, userUid); err != nil {
			t.Fatalf("failed to delete test user")
		}
	}()

	if err := sRepo.CreateSession(ctx, session); err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	log.Printf("successfully create session")
	ses, isFound, err := sRepo.GetSessionByUid(ctx, session.Uid)
	if err != nil {
		t.Fatalf("failed to get session: %v", err)
	}
	if !isFound {
		t.Fatalf("session not found")
	}

	log.Printf("getted by uid: %v", ses)
	log.Printf("successfully get session by uid")

	ses, isFound, err = sRepo.GetSessionByUserUid(ctx, session.UserUid)
	if err != nil {
		t.Fatalf("failed to get session by user: %v", err)
	}
	if !isFound {
		t.Fatalf("session not found by user")
	}

	log.Printf("getted by user %v", ses)
	log.Printf("successfully get session by user uid")

	updTime := time.Now()
	if err := sRepo.UpdateSession(ctx, session.Uid, updTime); err != nil {
		t.Fatalf("failed to update session: %v", err)
	}
	log.Printf("successfully update session")

	ses, isFound, err = sRepo.GetSessionByUid(ctx, session.Uid)
	if err != nil {
		t.Fatalf("failed to get session after update: %v", err)
	}
	if !isFound {
		t.Fatalf("session not found after update")
	}
	log.Printf("successfully get session after update")

	if !ses.UpdatedAt.Equal(updTime) {
		t.Fatalf("upd time not equal after update: upd: %s , getted: %s", updTime.String(), ses.UpdatedAt.String())
	}

	if err := sRepo.DeleteSessionByUid(ctx, session.Uid); err != nil {
		t.Fatalf("failed to delete session: %v", err)
	}
	log.Printf("successfully delete session by uid")

	_, isFound, err = sRepo.GetSessionByUid(ctx, session.Uid)
	if err != nil {
		t.Fatalf("failed to get session after update: %v", err)
	}
	if isFound {
		t.Fatalf("session found after delete")
	}

	if err := sRepo.CreateSession(ctx, session); err != nil {
		t.Fatalf("failed to create session")
	}
	log.Printf("successfully create session again")

	if err := sRepo.DeleteSessionByUserUid(ctx, session.UserUid); err != nil {
		t.Fatalf("failed to delete session by user uid: %v", err)
	}

	_, isFound, err = sRepo.GetSessionByUid(ctx, session.Uid)
	if err != nil {
		t.Fatalf("failed to get session after delete by user uid: %v", err)
	}
	if isFound {
		t.Fatalf("session found after delete by user uid")
	}
	log.Printf("successfully delete session by user uid")
}
