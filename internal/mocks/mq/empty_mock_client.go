package mqMock

import (
	"context"

	mqClient "github.com/balobas/sport_city_common/clients/mq"
)

type EmptyMqClientMock struct {
}

func NewEmptyMqClientMock() mqClient.MqClient {
	return &EmptyMqClientMock{}
}

func (m *EmptyMqClientMock) Publish(ctx context.Context, subj string, data []byte) error {
	return nil
}

func (m *EmptyMqClientMock) Subscribe(ctx context.Context, handlers map[string]map[string]mqClient.MqMsgHandler) error {
	return nil
}

func (m *EmptyMqClientMock) SubscribeV2(ctx context.Context, handlers map[string]map[string]mqClient.MqMsgHandler) error {
	return nil
}

func (m *EmptyMqClientMock) Close(ctx context.Context) error {
	return nil
}
