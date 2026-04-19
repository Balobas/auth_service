package deliveryGrpc

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
	uuid "github.com/satori/go.uuid"
)

func (s *AuthServerGrpc) Login(ctx context.Context, req *auth_v1.LoginRequest) (*auth_v1.JwtResponse, error) {
	log.Printf("authServerGrpc.Login: email %s, device uid %s", req.GetEmail(), req.GetDevice().GetUid())
	accessjwt, refreshJwt, err := s.ucAuth.Login(
		ctx, entity.LoginParams{
			Email:    req.GetEmail(),
			Password: req.GetPassword(),
			Device: entity.UserDevice{
				Uid:      uuid.FromStringOrNil(req.GetDevice().GetUid()),
				Name:     req.GetDevice().GetName(),
				Agent:    req.GetDevice().GetAgent(),
				Language: req.GetDevice().GetLanguage(),
			},
		},
	)
	if err != nil {
		log.Printf("authServerGrpc.Login: failed to login user with email %s: %v", req.GetEmail(), err)
		return nil, err
	}

	return &auth_v1.JwtResponse{
		AccessJwt:  accessjwt,
		RefreshJwt: refreshJwt,
	}, nil
}
