package useCaseUsers

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
)

func (uc *UseCaseUsers) GetAdmins(ctx context.Context, token string) ([]entity.User, error) {
	// var userToRegister entity.User

	// if token == "" {
	// 	superAdmin, err := uc.createSuperAdmin(ctx, user)
	// 	if err != nil {
	// 		return uuid.UUID{}, errors.WithStack(err)
	// 	}

	// 	userToRegister = superAdmin
	// } else {
	// 	admin, err := uc.createAdmin(ctx, user, token)
	// 	if err != nil {
	// 		return uuid.UUID{}, errors.WithStack(err)
	// 	}

	// 	userToRegister = admin
	// }

	// log.Printf("CreateAdmin: registering user\n")

	// uid, err := uc.Register(ctx, userToRegister, password)
	// if err != nil {
	// 	return uuid.UUID{}, errors.WithStack(err)
	// }

	return []entity.User{}, nil
}
