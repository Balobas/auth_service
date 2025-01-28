package useCaseUsers

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func (uc *UseCaseUsers) getAdminUser(ctx context.Context, email string) (entity.User, bool, error) {
	creatingUser, isFound, err := uc.usersRepo.GetByEmail(ctx, email)
	if err != nil {
		log.Printf("CreateAdmin: failed to get user\n")
		return entity.User{}, false, fmt.Errorf("failed to get user")
	}

	return creatingUser, isFound, nil
}

func (uc *UseCaseUsers) createSuperAdmin(ctx context.Context, user entity.User) (entity.User, error) {
	superAdminMail := os.Getenv("SUPER_ADMIN_EMAIL")

	_, isFound, err := uc.getAdminUser(ctx, superAdminMail)
	if err != nil {
		log.Printf("CreateAdmin: failed to get super user\n")
		return entity.User{}, fmt.Errorf("failed to get super user")
	}

	if isFound {
		log.Printf("CreateAdmin: super admin already exists \n")
		return entity.User{}, fmt.Errorf("super admin already exists")
	}

	if user.Email != superAdminMail {
		log.Printf("CreateAdmin: super admin invalid email \n")
		return entity.User{}, fmt.Errorf("super admin invalid email")
	}

	superUser := entity.User{
		Email: superAdminMail,
		Role:  entity.UserRoleAdmin,
	}

	return superUser, nil
}

func (uc *UseCaseUsers) createAdmin(ctx context.Context, user entity.User, token string) (entity.User, error) {
	tokenInfo, err := uc.jwtManager.ParseToken(token)
	if err != nil {
		log.Printf("CreateAdmin: failed to parse token\n")
		return entity.User{}, errors.WithStack(err)
	}

	creatingUser, isFound, err := uc.getAdminUser(ctx, tokenInfo.Email)
	if err != nil {
		log.Printf("CreateAdmin: failed to get user\n")
		return entity.User{}, fmt.Errorf("failed to get user")
	}

	if !isFound {
		log.Printf("CreateAdmin: failed to find user\n")
		return entity.User{}, fmt.Errorf("failed to find user")
	}

	if creatingUser.Role != entity.UserRoleAdmin || creatingUser.Role != entity.UserRole(tokenInfo.Role) {
		log.Printf("role from token is invalid: cannot create admin with this role")
		return entity.User{}, fmt.Errorf("role is invalid")
	}

	user.Role = entity.UserRoleAdmin

	return user, nil
}

func (uc *UseCaseUsers) CreateAdmin(ctx context.Context, user entity.User, password string, token string) (uuid.UUID, error) {
	var userToRegister entity.User

	if token == "" {
		superAdmin, err := uc.createSuperAdmin(ctx, user)
		if err != nil {
			return uuid.UUID{}, errors.WithStack(err)
		}

		userToRegister = superAdmin
	} else {
		admin, err := uc.createAdmin(ctx, user, token)
		if err != nil {
			return uuid.UUID{}, errors.WithStack(err)
		}

		userToRegister = admin
	}

	log.Printf("CreateAdmin: registering user\n")

	uid, err := uc.Register(ctx, userToRegister, password)
	if err != nil {
		return uuid.UUID{}, errors.WithStack(err)
	}

	return uid, nil
}
