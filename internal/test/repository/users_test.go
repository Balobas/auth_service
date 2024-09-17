package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	repositoryUsers "github.com/balobas/auth_service/internal/repository/postgres/users"
	uuid "github.com/satori/go.uuid"
)

func TestUsersRepo(t *testing.T) {
	ctx := context.Background()

	c := NewPgClient(t, ctx)

	usersRepo := repositoryUsers.New(c)

	user := entity.User{
		Uid:         uuid.NewV4(),
		Email:       "test@m.ru",
		Role:        "user",
		Permissions: []entity.UserPermission{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := usersRepo.CreateUser(ctx, user); err != nil {
		t.Fatal(err)
	}

	getUsr, isFound, err := usersRepo.GetUserByUid(ctx, user.Uid)
	if err != nil {
		t.Fatal(err)
	}
	if !isFound {
		t.Fatal("user not found by uid")
	}

	if getUsr.Email != user.Email || getUsr.Role != user.Role || !getUsr.CreatedAt.Equal(user.CreatedAt) || !getUsr.UpdatedAt.Equal(user.UpdatedAt) {
		t.Fatalf("getUsr by uid not match user: %v \n %v", user, getUsr)
	}

	getUsr, isFound, err = usersRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		t.Fatal(err)
	}
	if !isFound {
		t.Fatal("user not found by email")
	}

	if !uuid.Equal(getUsr.Uid, user.Uid) || getUsr.Role != user.Role || !getUsr.CreatedAt.Equal(user.CreatedAt) || !getUsr.UpdatedAt.Equal(user.UpdatedAt) {
		t.Fatalf("getUsr by email not match user: %v \n %v", user, getUsr)
	}

	user.Email = "huhuhuhuhuhuhu"
	user.UpdatedAt = time.Now()

	if err := usersRepo.UpdateUser(ctx, user); err != nil {
		t.Fatal(err)
	}

	getUsr, isFound, err = usersRepo.GetUserByUid(ctx, user.Uid)
	if err != nil {
		t.Fatal(err)
	}
	if !isFound {
		t.Fatal("user not found by uid")
	}

	if getUsr.Email != user.Email || getUsr.Role != user.Role || !getUsr.CreatedAt.Equal(user.CreatedAt) || !getUsr.UpdatedAt.Equal(user.UpdatedAt) {
		t.Fatalf("getUsr by uid not match updated user: %v \n %v", user, getUsr)
	}

	if err := usersRepo.DeleteUser(ctx, user.Uid); err != nil {
		t.Fatal(err)
	}

	_, isFound, err = usersRepo.GetUserByUid(ctx, user.Uid)
	if err != nil {
		t.Fatal(err)
	}
	if isFound {
		t.Fatal("deleted user found by uid")
	}

}
