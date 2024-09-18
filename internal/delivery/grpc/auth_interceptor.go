package deliveryGrpc

import (
	"context"
	"errors"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (s *AuthServerGrpc) UnaryAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

		if _, allowWithoutAuth := withoutAuth[info.FullMethod]; allowWithoutAuth {
			log.Printf("method %s allowed without auth", info.FullMethod)
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Printf("empty metadata")
			return nil, errors.New("token not provided")
		}
		accessJwtMd := md.Get("accessJwt")
		if len(accessJwtMd) == 0 {
			return nil, errors.New("empty token")
		}

		accessJwt := accessJwtMd[0]
		if len(accessJwt) == 0 {
			return nil, errors.New("invalid access jwt")
		}

		tokenInfo, err := s.ucAuth.VerifyAuth(ctx, accessJwt)
		if err != nil {
			log.Printf("failed to verify token: %v", err)
			return nil, err
		}

		log.Printf("user %s successfully verified", tokenInfo.UserUid)

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
