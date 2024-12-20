package deliveryGrpc

import (
	"context"
	"errors"
	"log"

	"github.com/balobas/auth_service/pkg/auth_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *AuthServerGrpc) GetAdmins(ctx context.Context, _ *emptypb.Empty) (*auth_v1.AdminsResponse, error) {
	log.Printf("auth.GetAdmins\n")

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Printf("failed to get token context\n")
		return nil, errors.New("failed to get token context")
	}

	accessJwtSlice := md.Get("accessJwt")
	if len(accessJwtSlice) == 0 {
		log.Printf("failed to get token \n")
		return nil, errors.New("failed to get token")
	}

	users, err := s.ucUsers.GetAdmins(ctx, accessJwtSlice[0])
	if err != nil {
		log.Printf("failed to get admin users\n")
		return nil, err
	}

	respUsers := make([]*auth_v1.GetUserResponse, len(users))
	for ind, user := range users {
		respUsers[ind] = &auth_v1.GetUserResponse{
			Uid:         user.Uid.String(),
			Email:       user.Email,
			Role:        1,
			Permissions: user.PermissionsStrings(),
			CreatedAt:   timestamppb.New(user.CreatedAt),
			UpdatedAt:   timestamppb.New(user.UpdatedAt),
		}
	}

	return &auth_v1.AdminsResponse{
		Users: respUsers,
	}, nil
}
