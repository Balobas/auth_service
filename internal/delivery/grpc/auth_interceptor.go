package deliveryGrpc

import (
	"context"
	"errors"

	"github.com/balobas/auth_service/internal/entity"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

func (s *AuthServerGrpc) UnaryAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

		if _, allowWithoutAuth := withoutAuth[info.FullMethod]; allowWithoutAuth {
			return handler(ctx, req)
		}

		accessJwt, err := accessJwtFromContext(ctx)
		if err != nil {
			return nil, err
		}

		tokenInfo, err := s.ucAuth.VerifyAuth(ctx, accessJwt)
		if err != nil {
			return nil, err
		}

		return handler(
			contextWithUserInfo(ctx, tokenInfo),
			req,
		)
	}
}

var (
	withoutAuth = map[string]struct{}{
		"/auth.Auth/VerifyEmail": {},
		"/auth.Auth/Verify":      {},
		"/auth.Auth/Login":       {},
		"/auth.Auth/Register":    {},
		"/auth.Auth/Refresh":     {},
	}
)

func accessJwtFromContext(ctx context.Context) (string, error) {
	token, ok := (ctx.Value("accessJwt")).(string)
	if !ok {
		return "", errors.New("invalid access jwt")
	}

	if len(token) == 0 {
		return "", errors.New("invalid access jwt")
	}

	return token, nil
}

type userCtxKey struct{}

func contextWithUserInfo(ctx context.Context, tokenInfo entity.TokenInfo) context.Context {
	return context.WithValue(
		ctx, userCtxKey{},
		UserInfo{
			UserUid: tokenInfo.UserUid,
			Role:    entity.UserRole(tokenInfo.Role),
		},
	)
}

func userInfoFromContext(ctx context.Context) UserInfo {
	info, ok := ctx.Value(userCtxKey{}).(UserInfo)
	if !ok {
		return UserInfo{}
	}
	return info
}

type UserInfo struct {
	UserUid uuid.UUID
	Role    entity.UserRole
}
