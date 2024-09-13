package deliveryGrpc

import (
	"context"
	"errors"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *AuthServerGrpc) GetUser(ctx context.Context, req *auth_v1.GetUserRequest) (*auth_v1.GetUserResponse, error) {
	log.Printf("auth.GetUser\n")
	uid, email := req.GetUid(), req.GetEmail()

	if len(uid) == 0 && len(email) == 0 {
		log.Printf("auth.GetUser empty request\n")
		return nil, errors.New("empty request")
	}

	var (
		user    entity.User
		err     error
		isFound bool
	)
	if len(uid) == 0 {
		user, isFound, err = s.ucUsers.GetUserByEmail(ctx, email)
	} else {

		UID, uidErr := uuid.FromString(uid)
		if uidErr != nil {
			log.Printf("auth.GetUser invalid uid %v\n", uidErr)
			return nil, uidErr
		}

		user, isFound, err = s.ucUsers.GetUserByUid(ctx, UID)
	}

	if err != nil {
		log.Printf("auth.GetUser error %v\n", err)
		return nil, err
	}

	if !isFound {
		log.Printf("auth.GetUser not found user\n")
		return nil, errors.New("user not found")
	}

	return &auth_v1.GetUserResponse{
		Uid:         user.Uid.String(),
		Email:       user.Email,
		Role:        auth_v1.Role(auth_v1.Role_value[string(user.Role)]),
		Permissions: user.PermissionsStrings(),
		CreatedAt:   timestamppb.New(user.CreatedAt),
		UpdatedAt:   timestamppb.New(user.UpdatedAt),
	}, nil
}
