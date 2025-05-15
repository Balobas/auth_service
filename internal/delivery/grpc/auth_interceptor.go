package deliveryGrpc

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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
			return nil, status.Error(codes.Unauthenticated, AuthErrMsgTokenNotProvided)
		}
		accessJwtMd := md.Get("accessJwt")
		if len(accessJwtMd) == 0 {
			log.Printf("empty accessJwt in metadata")
			return nil, status.Error(codes.Unauthenticated, AuthErrMsgTokenNotProvided)
		}

		accessJwt := accessJwtMd[0]
		if len(accessJwt) == 0 {
			log.Printf("empty accessJwt[0] in metadata")
			return nil, status.Error(codes.Unauthenticated, AuthErrMsgTokenNotProvided)
		}

		tokenInfo, err := s.ucAuth.VerifyAuth(ctx, accessJwt)
		if err != nil {
			log.Printf("failed to verify token: %v", err)
			return nil, status.Error(codes.Unauthenticated, err.Error())
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
		"/auth.Auth/CreateAdmin": {},
		"/auth.Auth/HealthCheck": {},
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

const (
	AuthErrMsgTokenNotProvided = "token not provided"
	AuthErrMsgInvalidToken     = "invalid token"
	AuthErrMsgTokenExpired     = "token expired"
)
