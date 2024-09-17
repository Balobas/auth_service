package repo_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/internal/repository/postgres/credentials"
	repositoryUsers "github.com/balobas/auth_service/internal/repository/postgres/users"
	uuid "github.com/satori/go.uuid"
)

func TestCredsRepo(t *testing.T) {
	ctx := context.Background()

	c := NewPgClient(t, ctx)

	credsRepo := credentials.New(c)

	userUid := uuid.NewV4()
	creds := entity.UserCredentials{
		UserUid:      userUid,
		PasswordHash: []byte("hasdaisdjklajds"),
	}

	if err := credsRepo.CreateCredentials(ctx, creds); err == nil {
		t.Fatalf("creds created without existed user")
	}

	// user exist
	usersRepo := repositoryUsers.New(c)
	err := usersRepo.CreateUser(ctx, entity.User{
		Uid:       userUid,
		Email:     "hui@test",
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

	if err := credsRepo.CreateCredentials(ctx, creds); err != nil {
		t.Fatalf("failed to create creds: %v", err)
	}

	cr, isFound, err := credsRepo.GetByUserUid(ctx, userUid)
	if err != nil {
		log.Fatalf("failed to get creds: %v", err)
	}
	if !isFound {
		log.Fatalf("creds not found")
	}

	if string(cr.PasswordHash) != string(creds.PasswordHash) {
		log.Fatalf("creds not equal")
	}

	creds.PasswordHash = []byte("ppppoasjijuoihaiow")
	if err := credsRepo.UpdateCredentials(ctx, creds); err != nil {
		log.Fatalf("failed to update creds: %v", err)
	}

	cr, isFound, err = credsRepo.GetByUserUid(ctx, userUid)
	if err != nil {
		log.Fatalf("failed to get creds after update: %v", err)
	}
	if !isFound {
		log.Fatalf("creds not found after update")
	}

	if string(cr.PasswordHash) != string(creds.PasswordHash) {
		log.Fatalf("creds not equal after update")
	}

	if err := credsRepo.DeleteByUserUid(ctx, userUid); err != nil {
		log.Fatalf("failed to delete creds: %v", err)
	}

	_, isFound, err = credsRepo.GetByUserUid(ctx, userUid)
	if err != nil {
		log.Fatalf("err get not existed: %v", err)
	}
	if isFound {
		log.Fatalf("creds exists after delete")
	}
}
