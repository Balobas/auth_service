package deliveryGrpcAuth

import (
	"context"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
	uuid "github.com/satori/go.uuid"
)

func deviceFromProto(d *auth_v1.Device) entity.LoginDeviceData {
	return entity.LoginDeviceData{
		Uid:      uuid.FromStringOrNil(d.GetUid()),
		Name:     d.GetName(),
		Agent:    d.GetAgent(),
		Language: d.GetLanguage(),
	}
}

func (s *AuthServerGrpc) AuthWithTelegram(ctx context.Context, req *auth_v1.TelegramAuthRequest) (*auth_v1.TelegramAuthResponse, error) {
	access, refresh, created, err := s.ucAuth.AuthWithTelegram(ctx, entity.TelegramAuthParams{
		InitData:        req.GetInitData(),
		Device:          deviceFromProto(req.GetDevice()),
		CreateIfMissing: req.GetCreateIfMissing(),
	})
	if err != nil {
		return nil, err
	}
	return &auth_v1.TelegramAuthResponse{
		AccessJwt:  access,
		RefreshJwt: refresh,
		Created:    created,
	}, nil
}

func (s *AuthServerGrpc) LinkTelegram(ctx context.Context, req *auth_v1.TelegramLinkRequest) (*auth_v1.JwtResponse, error) {
	access, refresh, err := s.ucAuth.LinkTelegram(ctx, entity.TelegramLinkParams{
		InitData: req.GetInitData(),
		Device:   deviceFromProto(req.GetDevice()),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, err
	}
	return &auth_v1.JwtResponse{AccessJwt: access, RefreshJwt: refresh}, nil
}
