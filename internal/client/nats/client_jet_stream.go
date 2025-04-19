package natsClient

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/balobas/auth_service/internal/client"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pkg/errors"
)

type NatsClientJetStream struct {
	cfg           Config
	conn          *nats.Conn
	js            jetstream.JetStream
	consumersCtxs []jetstream.ConsumeContext

	wg        *sync.WaitGroup
	connected chan struct{}
}

func NewJs(ctx context.Context, cfg Config) (client.MqClient, error) {
	connectedChan := make(chan struct{})

	conn, err := nats.Connect(
		cfg.NatsUrl(), nats.Name(cfg.NatsClientName()),
		nats.ReconnectHandler(func(c *nats.Conn) {
			select {
			case connectedChan <- struct{}{}:
			default:
			}
			log.Printf("nats has been recconected")
		}),
		nats.ErrorHandler(func(c *nats.Conn, s *nats.Subscription, err error) {
			log.Printf("nats error handler: error occured: sub %s : %v", s.Subject, err)
		}),
		nats.DisconnectHandler(func(c *nats.Conn) {
			log.Printf("nats disconnect")
		}),
		nats.ClosedHandler(func(c *nats.Conn) {
			log.Printf("nats closed")
		}),
		nats.MaxReconnects(-1),
		nats.ConnectHandler(func(c *nats.Conn) {

			select {
			case connectedChan <- struct{}{}:
			default:
			}
			log.Printf("nats successfully connected to %s", c.ConnectedAddr())
		}),
		nats.RetryOnFailedConnect(true),
	)
	if err != nil {
		log.Printf("failed to connect to nats (url: %s): %v", cfg.NatsUrl(), err)
		return nil, err
	}

	js, err := jetstream.New(conn)
	if err != nil {
		conn.Close()
		log.Printf("failed to init nats jet stream: %v", err)
		return nil, errors.WithStack(err)
	}

	return &NatsClientJetStream{
		cfg:       cfg,
		conn:      conn,
		js:        js,
		wg:        &sync.WaitGroup{},
		connected: connectedChan,
	}, nil
}

func (nc *NatsClientJetStream) Publish(ctx context.Context, subj string, data []byte) error {
	ack, err := nc.js.Publish(ctx, subj, data)
	if err != nil {
		log.Printf("failed to publish message into subject %s: %v", subj, err)
		return errors.WithStack(err)
	}
	log.Printf("ack info: stream %s, domain %s, duplicate %t, sequence %d", ack.Stream, ack.Domain, ack.Duplicate, ack.Sequence)
	return nil
}

func (nc *NatsClientJetStream) Subscribe(ctx context.Context, handlersStreams map[string]map[string]client.MqMsgHandler) error {
	failedStreams := map[string]map[string]client.MqMsgHandler{}

	defer nc.resubscribeOnFailedStreams(ctx, failedStreams)

	for streamName, handlers := range handlersStreams {

		stream, err := nc.js.Stream(ctx, streamName)
		if err != nil {
			log.Printf("failed to get stream %s: %v", streamName, err)
			failedStreams[streamName] = handlers
			continue
		}

		for subject, handler := range handlers {

			consumerName := fmt.Sprintf("%s_%s_consumer", nc.cfg.ServiceName(), strings.Split(subject, ".")[1])
			consumer, err := stream.Consumer(ctx, consumerName)
			if err != nil {
				log.Printf("failed to get consumer %s on stream %s subject %s: %v", consumerName, streamName, subject, err)
				addSubjectToFailedStreams(streamName, subject, handler, failedStreams)
				continue
			}

			consumerCtx, err := consumer.Consume(convertToNatsJsMsgHandler(ctx, handler))
			if err != nil {
				log.Printf("failed to init consumer %s on stream %s subject %s: %v", consumerName, streamName, subject, err)
				addSubjectToFailedStreams(streamName, subject, handler, failedStreams)
				continue
			}
			log.Printf("successfully init consumer %s on stream %s subject %s", consumerName, streamName, subject)

			nc.consumersCtxs = append(nc.consumersCtxs, consumerCtx)
		}
	}

	return nil
}

func addSubjectToFailedStreams(
	streamName string,
	subject string,
	handler client.MqMsgHandler,
	failedStreams map[string]map[string]client.MqMsgHandler,
) {
	failedStream, ok := failedStreams[streamName]
	if !ok {
		failedStreams[streamName] = map[string]client.MqMsgHandler{
			subject: handler,
		}
	} else {
		failedStream[subject] = handler
	}
}

func (nc *NatsClientJetStream) resubscribeOnFailedStreams(ctx context.Context, failedStreams map[string]map[string]client.MqMsgHandler) {
	if len(failedStreams) == 0 {
		return
	}

	nc.wg.Add(1)
	go func() {
		defer nc.wg.Done()

		select {
		case <-ctx.Done():
			log.Printf("stop attempts to init failed consumers")
			return
		case <-nc.connected:
			select {
			case <-ctx.Done():
				log.Printf("stop attempts to init failed consumers")
				return
			default:
			}

			nc.Subscribe(ctx, failedStreams)
		case <-time.After(1 * time.Minute):
			select {
			case <-ctx.Done():
				log.Printf("stop attempts to init failed consumers")
				return
			default:
			}

			nc.Subscribe(ctx, failedStreams)
		}
	}()
}

func (nc *NatsClientJetStream) Close(ctx context.Context) error {
	for _, consumerCtx := range nc.consumersCtxs {
		consumerCtx.Stop()
	}
	nc.conn.Close()
	nc.wg.Wait()
	close(nc.connected)
	log.Printf("nats js client closed successfully")
	return nil
}

func convertToNatsJsMsgHandler(ctx context.Context, handler client.MqMsgHandler) jetstream.MessageHandler {
	return func(msg jetstream.Msg) {
		if err := handler(ctx, msg.Data()); err != nil {
			log.Printf("failed to handle message %s: %v", msg.Data(), err)
			nackJsMsgWithLog(msg)
			return
		}

		if err := msg.Ack(); err != nil {
			log.Printf("failed to ack message %s: %v", msg.Data(), err)
		}
	}
}

func nackJsMsgWithLog(msg jetstream.Msg) {
	if err := msg.Nak(); err != nil {
		log.Printf("failed to nack msg %s: %v", msg.Data(), err)
	}
}
