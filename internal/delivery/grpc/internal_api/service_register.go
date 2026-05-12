package deliveryGrpcAuthInternalApi

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_internal_api"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AuthInternalApiServerGrpc) ServiceRegister(ctx context.Context, req *auth_internal_api.ServiceRegisterRequest) (*emptypb.Empty, error) {
	log.Printf("AuthInternalApiServerGrpc.ServiceRegister: name %s, domain %s, uid %s", req.GetName(), req.GetDomain(), req.GetUid())

	uid, err := uuid.FromString(req.GetUid())
	if err != nil {
		return nil, err
	}

	err = s.ucServices.Register(
		ctx, entity.RegisterServiceRequest{
			Uid:      uid,
			Name:     req.GetName(),
			Domain:   req.GetDomain(),
			Password: req.GetPassword(),
		},
	)
	if err != nil {
		log.Printf("AuthInternalApiServerGrpc.ServiceRegister: failed to register service with uid %s: %v", req.GetUid(), err)
		return nil, err
	}

	return nil, nil
}
