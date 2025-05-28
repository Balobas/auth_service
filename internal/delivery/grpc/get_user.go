package deliveryGrpc

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *AuthServerGrpc) GetUser(ctx context.Context, req *auth_v1.GetUserRequest) (*auth_v1.GetUserResponse, error) {
	log.Printf("authServerGrpc.GetUser\n")
	uid, email := uuid.FromStringOrNil(req.GetUid()), req.GetEmail()

	if uuid.Equal(uid, uuid.UUID{}) && len(email) == 0 {
		log.Printf("authServerGrpc.GetUser empty request")
		return nil, errors.Wrap(serviceErrors.ErrBadRequest, "empty request")
	}

	var (
		user entity.User
		err  error
	)
	if uuid.Equal(uid, uuid.UUID{}) {
		user, err = s.ucUsers.GetUserByEmail(ctx, email)
	} else {
		user, err = s.ucUsers.GetUserByUid(ctx, uid)
	}

	if err != nil {
		log.Printf("authServerGrpc.GetUser error %v\n", err)
		return nil, errors.WithStack(err)
	}

	return &auth_v1.GetUserResponse{
		Uid:       user.Uid.String(),
		Email:     user.Email,
		Roles:     user.Roles,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}
