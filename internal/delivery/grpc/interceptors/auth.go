package deliveryGrpcInterceptors

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/entity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *Provider) UnaryAuthInterceptor() grpc.UnaryServerInterceptor {
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
			contextWithUserInfo(ctx, tokenInfo, accessJwt),
			req,
		)
	}
}

var (
	withoutAuth = map[string]struct{}{
		"/auth.Auth/VerifyEmail":                             {},
		"/auth.Auth/Verify":                                  {},
		"/auth.Auth/VerifyAccess":                            {},
		"/auth.Auth/Login":                                   {},
		"/auth.Auth/Register":                                {},
		"/auth.Auth/Refresh":                                 {},
		"/auth.Auth/CreateAdmin":                             {},
		"/auth.Auth/HealthCheck":                             {},
		"/auth_internal_api.AuthInternalApi/ServiceRegister": {},
		"/auth_internal_api.AuthInternalApi/ServiceLogin":    {},
		"/auth_internal_api.AuthInternalApi/ServiceRefresh":  {},
	}
)

type userCtxKey struct{}

func contextWithUserInfo(ctx context.Context, tokenInfo entity.TokenInfo, tokenStr string) context.Context {
	return context.WithValue(
		ctx, userCtxKey{},
		entity.UserInfo{
			UserUid:     tokenInfo.UserUid,
			DeviceUid:   tokenInfo.DeviceUid,
			Roles:       tokenInfo.Roles,
			Token:       tokenStr,
			IsSystem:    tokenInfo.Type == entity.TokenTypeSystem,
			ServiceUid:  tokenInfo.ServiceUid,
			ServiceName: tokenInfo.ServiceName,
			Domain:      tokenInfo.Domain,
		},
	)
}

func UserInfoFromContext(ctx context.Context) entity.UserInfo {
	info, ok := ctx.Value(userCtxKey{}).(entity.UserInfo)
	if !ok {
		return entity.UserInfo{}
	}
	return info
}

const (
	AuthErrMsgTokenNotProvided = "token not provided"
	AuthErrMsgInvalidToken     = "invalid token"
	AuthErrMsgTokenExpired     = "token expired"
)
