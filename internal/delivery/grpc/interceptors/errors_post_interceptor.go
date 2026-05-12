package deliveryGrpcInterceptors

import (
	"context"

	serviceErrors "github.com/balobas/auth_service/pkg/service_errors"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Provider) UnaryErrorsPostInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

		res, err := handler(ctx, req)
		return res, grpcErrorWithStatus(err)
	}
}

func grpcErrorWithStatus(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, serviceErrors.ErrBadRequest):
		// TODO: подумать над расширением, часть ошибок можно сделать FailedPrecondition
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, serviceErrors.ErrNotFound):

		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, serviceErrors.ErrIdempotentOperation),
		errors.Is(err, serviceErrors.ErrAlreadyExists):

		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, serviceErrors.ErrNotAllowedByPermissions):

		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, serviceErrors.ErrInvalidToken),
		errors.Is(err, serviceErrors.ErrTokenNotProvided),
		errors.Is(err, serviceErrors.ErrTokenExpired):

		return status.Error(codes.Unauthenticated, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
