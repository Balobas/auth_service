package converterGrpc

import (
	"time"

	"github.com/balobas/auth_service/internal/entity"
	"github.com/balobas/auth_service/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromCreateUserRequestToUserEntity(userProto *auth_v1.CreateRequest) entity.User {
	return entity.User{
		Name:            userProto.Name,
		Email:           userProto.Email,
		Password:        userProto.Password,
		ConfirmPassword: userProto.PasswordConfirm,
		Role:            userProto.Role.String(),
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}
}

func FromUserEntityToGetResponse(user entity.User) *auth_v1.GetResponse {
	return &auth_v1.GetResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Role:      auth_v1.Role(auth_v1.Role_value[user.Role]),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func FromUpdateRequestToUserEntity(updateReq *auth_v1.UpdateRequest) entity.User {
	return entity.User{
		Id:    updateReq.Id,
		Name:  updateReq.GetName().GetValue(),
		Email: updateReq.GetEmail().GetValue(),
	}
}
