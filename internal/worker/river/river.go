package riverWorkers

import (
	"context"
	"time"

	DBclient "github.com/balobas/sport_city_common/clients/database"
	outboxEntity "github.com/balobas/sport_city_common/entity/outbox"
	riverCommon "github.com/balobas/sport_city_common/worker/river"
	riverOutboxPublisher "github.com/balobas/sport_city_common/worker/river/outbox_publisher"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

var driver riverdriver.Driver[pgx.Tx] = &riverpgxv5.Driver{}

type Config interface {
	Queues() map[string]int
	QueueNames() []string
	MaxAttempts() int
	MaxWorkers() int
	NextRetry() time.Duration
	JobTimeout() time.Duration
	FetchCooldown() time.Duration
	FetchPollInterval() time.Duration
}

type Client struct {
	*riverCommon.Client
}

func NewClient(cfg Config, dbClient DBclient.ClientDB) (*Client, error) {
	c, err := riverCommon.NewClient(cfg, dbClient)
	if err != nil {
		return nil, err
	}
	return &Client{Client: c}, nil
}

func (c *Client) BuildWorkers(
	ctx context.Context,
	outboxPublisher *riverOutboxPublisher.Worker,
) {
	if outboxPublisher != nil {
		riverCommon.AddWorker(c.Client, river.WorkFunc(outboxPublisher.Work))
	}
}

func (c *Client) CreateTaskSendOutboxMessage(ctx context.Context, message outboxEntity.Message) error {
	args := riverOutboxPublisher.Args{
		Message: message,
	}
	opts := &river.InsertOpts{
		Queue: args.Kind(),
	}

	return c.InsertRiver(ctx, args, opts)
}
