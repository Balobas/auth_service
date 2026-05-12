package deliveryGrpcAuthInternalApi

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_internal_api"
	uuid "github.com/satori/go.uuid"
)

func (s *AuthInternalApiServerGrpc) ServiceLogin(ctx context.Context, req *auth_internal_api.ServiceLoginRequest) (*auth_internal_api.JwtResponse, error) {
	log.Printf("AuthInternalApiServerGrpc.ServiceLogin: service uid %s, device uid %s", req.GetUid(), req.GetDevice().GetUid())

	uid, err := uuid.FromString(req.GetUid())
	if err != nil {
		return nil, err
	}

	accessjwt, refreshJwt, err := s.ucAuth.ServiceLogin(
		ctx, entity.ServiceLoginParams{
			Uid:      uid,
			Password: req.GetPassword(),
			Device: entity.LoginDeviceData{
				Uid:      uuid.FromStringOrNil(req.GetDevice().GetUid()),
				Name:     req.GetDevice().GetName(),
				Agent:    req.GetDevice().GetAgent(),
				Language: req.GetDevice().GetLanguage(),
			},
		},
	)
	if err != nil {
		log.Printf("AuthInternalApiServerGrpc.ServiceLogin: failed to login service with uid %s: %v", req.GetUid(), err)
		return nil, err
	}

	return &auth_internal_api.JwtResponse{
		AccessJwt:  accessjwt,
		RefreshJwt: refreshJwt,
	}, nil
}
