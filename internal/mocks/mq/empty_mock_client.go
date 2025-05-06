package mqMock

import (
	"context"

	"github.com/balobas/auth_service/internal/client"
)

type EmptyMqClientMock struct {
}

func NewEmptyMqClientMock() client.MqClient {
	return &EmptyMqClientMock{}
}

func (m *EmptyMqClientMock) Publish(ctx context.Context, subj string, data []byte) error {
	return nil
}

func (m *EmptyMqClientMock) Subscribe(ctx context.Context, handlers map[string]map[string]client.MqMsgHandler) error {
	return nil
}

func (m *EmptyMqClientMock) Close(ctx context.Context) error {
	return nil
}
