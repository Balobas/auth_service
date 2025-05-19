package useCaseUsers

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/balobas/auth_service/internal/entity"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

//TODO: отрефакторить, тут не надо проверять права и роль вызывающего метод юзера
// перенести проверки прав и ролей в delivery слой

func (uc *UseCaseUsers) getAdminUser(ctx context.Context, email string) (entity.User, bool, error) {
	log.Printf("usecaseUsers.getAdminUser: email %s", email)

	creatingUser, isFound, err := uc.usersRepo.GetByEmail(ctx, email)
	if err != nil {
		log.Printf("usecaseUsers.getAdminUser: failed to get user by email %s: %v", email, email)
		return entity.User{}, false, errors.Wrap(err, "failed to get user")
	}

	return creatingUser, isFound, nil
}

func (uc *UseCaseUsers) createSuperAdmin(ctx context.Context, user entity.User) (entity.User, error) {
	log.Printf("usecaseUsers.createSuperAdmin: user uid %s", user.Uid)

	superAdminMail := os.Getenv("SUPER_ADMIN_EMAIL")

	_, isFound, err := uc.getAdminUser(ctx, superAdminMail)
	if err != nil {
		log.Printf("usecaseUsers.createSuperAdmin: failed to get super user by email %s: %v", superAdminMail, err)
		return entity.User{}, fmt.Errorf("failed to get super user")
	}

	if isFound {
		log.Printf("usecaseUsers.createSuperAdmin: super admin already exists \n")
		return entity.User{}, errors.Wrap(serviceErrors.ErrAlreadyExists, "super admin already exists")
	}

	if user.Email != superAdminMail {
		log.Printf("usecaseUsers.createSuperAdmin: super admin invalid email \n")
		return entity.User{}, errors.Wrap(serviceErrors.ErrBadRequest, "super admin invalid email")
	}

	superUser := entity.User{
		Email: superAdminMail,
		Role:  entity.UserRoleAdmin,
	}

	return superUser, nil
}

func (uc *UseCaseUsers) createAdmin(ctx context.Context, user entity.User, token string) (entity.User, error) {
	log.Printf("usecaseUsers.createAdmin: user uid %s", user.Uid)

	tokenInfo, err := uc.jwtManager.ParseToken(token)
	if err != nil {
		log.Printf("usecaseUsers.createAdmin: failed to parse token: %v", err)
		return entity.User{}, errors.WithStack(err)
	}

	creatingUser, isFound, err := uc.getAdminUser(ctx, tokenInfo.Email)
	if err != nil {
		log.Printf("usecaseUsers.createAdmin: failed to get user by email %s: %v", tokenInfo.Email, err)
		return entity.User{}, fmt.Errorf("failed to get user")
	}

	if !isFound {
		log.Printf("usecaseUsers.createAdmin: failed to find user by email %s", tokenInfo.Email)
		return entity.User{}, errors.Wrap(serviceErrors.ErrNotFound, "caller user")
	}

	if creatingUser.Role != entity.UserRoleAdmin || creatingUser.Role != entity.UserRole(tokenInfo.Role) {
		log.Printf("usecaseUsers.createAdmin: role from token is invalid: cannot create admin with this role")
		return entity.User{}, errors.Wrap(serviceErrors.ErrNotAllowedByPermissions, "role is invalid")
	}

	user.Role = entity.UserRoleAdmin

	return user, nil
}

func (uc *UseCaseUsers) CreateAdmin(ctx context.Context, user entity.User, password string, token string) (uuid.UUID, error) {
	log.Printf("usecaseAuth.CreateAdmin: user %s", user.Uid)

	var userToRegister entity.User

	if token == "" {
		superAdmin, err := uc.createSuperAdmin(ctx, user)
		if err != nil {
			log.Printf("usecaseAuth.CreateAdmin: failed to create super admin %s: %v", user.Uid, err)
			return uuid.UUID{}, errors.WithStack(err)
		}

		userToRegister = superAdmin
	} else {
		admin, err := uc.createAdmin(ctx, user, token)
		if err != nil {
			log.Printf("usecaseAuth.CreateAdmin: failed to create admin %s: %v", user.Uid, err)
			return uuid.UUID{}, errors.WithStack(err)
		}

		userToRegister = admin
	}

	uid, err := uc.Register(ctx, userToRegister, password)
	if err != nil {
		log.Printf("usecaseAuth.CreateAdmin: failed to register admin %s: %v", user.Uid, err)
		return uuid.UUID{}, errors.WithStack(err)
	}

	return uid, nil
}
