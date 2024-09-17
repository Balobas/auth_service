package repo_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	repositoryUsers "github.com/balobas/auth_service/internal/repository/postgres/users"
	"github.com/balobas/auth_service/internal/repository/postgres/verification"
	uuid "github.com/satori/go.uuid"
)

func randStr() string {
	b := make([]byte, 12)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2:12]
}

func TestVerificationRepo(t *testing.T) {
	ctx := context.Background()

	c := NewPgClient(t, ctx)

	userUid := uuid.NewV4()
	ver := entity.Verification{
		UserUid:   userUid,
		Email:     fmt.Sprintf("%s@%s", randStr(), randStr()),
		Token:     randStr(),
		Status:    entity.VerificationStatusCreated,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	verRepo := verification.New(c)

	if err := verRepo.CreateVerification(ctx, ver); err == nil {
		log.Fatalf("created verification without user")
	}

	// user exists
	usersRepo := repositoryUsers.New(c)

	err := usersRepo.CreateUser(ctx, entity.User{
		Uid:       userUid,
		Email:     ver.Email,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Fatalf("failed to create user %v", err)
	}
	defer func() {
		// delete test user
		if err := usersRepo.DeleteUser(ctx, userUid); err != nil {
			log.Fatalf("failed to delete test user")
		}
	}()

	if err := verRepo.CreateVerification(ctx, ver); err != nil {
		log.Fatalf("failed to create verification: %v", err)
	}

	v, isFound, err := verRepo.GetVerificationByToken(ctx, ver.Token)
	if err != nil {
		log.Fatalf("failed to get verification by token: %v", err)
	}
	if !isFound {
		log.Fatalf("verification not found")
	}

	log.Printf("getted verification: %v", v)

	ver.Status = entity.VerificationStatusWaiting

	if err := verRepo.UpdateVerification(ctx, ver); err != nil {
		log.Fatalf("failed to update verification: %v", err)
	}

	v, isFound, err = verRepo.GetVerificationByToken(ctx, ver.Token)
	if err != nil {
		log.Fatalf("failed to get verification by token: %v", err)
	}
	if !isFound {
		log.Fatalf("verification not found")
	}

	if v.Status != ver.Status || v.Email != ver.Email {
		log.Fatalf("getted verification does not match: getted: %v, orig: %v", v, ver)
	}

	vrf, err := verRepo.GetVerificationsInStatus(ctx, entity.VerificationStatusWaiting, 2)
	if err != nil {
		log.Fatalf("failed to get verifications in status: %v", err)
	}

	log.Printf("verifications in status %s limit %d : %v", entity.VerificationStatusWaiting, 2, vrf)

	if err := verRepo.DeleteVerification(ctx, userUid); err != nil {
		log.Fatalf("failed to delete verification: %v", err)
	}

	_, isFound, err = verRepo.GetVerificationByToken(ctx, ver.Token)
	if err != nil {
		log.Fatalf("failed to get verification by token (deleted): %v", err)
	}
	if isFound {
		log.Fatalf("verification found")
	}

}
