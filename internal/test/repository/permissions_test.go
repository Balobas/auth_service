package repo_test

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/internal/repository/postgres/permissions"
	repositoryUsers "github.com/balobas/auth_service/internal/repository/postgres/users"
	"github.com/jackc/pgx/v4"
	uuid "github.com/satori/go.uuid"
)

func TestPermissionsRepo(t *testing.T) {
	ctx := context.Background()

	c := NewPgClient(t, ctx)

	permsRepo := permissions.New(c)

	// юзера не существует
	userUid := uuid.NewV4()
	perms := []entity.UserPermission{entity.UserPermissionNotVerified}

	err := permsRepo.CreateUserPermissions(ctx, userUid, perms)
	if err == nil {
		log.Fatalf("permissions created, but it does not wanted")
	}

	// юзер есть
	usersRepo := repositoryUsers.New(c)
	err = usersRepo.CreateUser(ctx, entity.User{
		Uid:       userUid,
		Email:     "hui@test",
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("failed to create user for test. err: %v", err)
	}
	defer func() {
		// delete test user
		if err := usersRepo.DeleteUser(ctx, userUid); err != nil {
			t.Fatalf("failed to delete test user")
		}
	}()

	err = permsRepo.CreateUserPermissions(ctx, userUid, perms)
	if err != nil {
		t.Fatalf("failed to create permissions: %v", err)
	}

	p, err := permsRepo.GetUserPermissions(ctx, userUid)
	if err != nil {
		t.Fatalf("failed to get user permissions: %v", err)
	}

	if !isPermsEqual(p, perms) {
		t.Fatalf("getted perms not equal: getted: %v, orig: %v", p, perms)
	}

	perms = []entity.UserPermission{entity.UserPermissionBase}
	if err := permsRepo.UpdateUserPermissions(ctx, userUid, perms); err != nil {
		t.Fatalf("failed to update permissions: %v", err)
	}

	p, err = permsRepo.GetUserPermissions(ctx, userUid)
	if err != nil {
		t.Fatalf("failed to get user after update permissions: %v", err)
	}

	if !isPermsEqual(p, perms) {
		t.Fatalf("getted perms not equal after update: getted: %v, orig: %v", p, perms)
	}

	if err := permsRepo.DeleteUserPermissions(ctx, userUid); err != nil {
		t.Fatalf("failed to delete user permissions: %v", err)
	}

	_, err = permsRepo.GetUserPermissions(ctx, userUid)
	if err == nil {
		t.Fatalf("permissions not deleted")
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("uncorect err: %v", err)
	}

}

func isPermsEqual(p1, p2 []entity.UserPermission) bool {
	if len(p1) != len(p2) {
		return false
	}

	p1m := make(map[entity.UserPermission]struct{}, len(p1))
	for i := 0; i < len(p1); i++ {
		p1m[p1[i]] = struct{}{}
	}

	for i := 0; i < len(p2); i++ {
		if _, ok := p1m[p2[i]]; !ok {
			return false
		}
	}
	return true
}
