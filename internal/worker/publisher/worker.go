package workerPublisher

import (
	"context"
	"log"
	"time"
)

type Worker struct {
	cfg              Config
	outboxRepository OutboxRepository
	publisher        Publisher
}

func New(
	cfg Config,
	outboxRepository OutboxRepository,
	publisher Publisher,
) *Worker {
	return &Worker{
		cfg:              cfg,
		outboxRepository: outboxRepository,
		publisher:        publisher,
	}
}

func (w *Worker) Run(ctx context.Context) {
	log.Printf("start publisher worker\n")
	msgsBatchSize := w.cfg.MqPublishMessagesBatchSize()
	timer := time.NewTimer(w.cfg.MqPublishMessagesInterval())

	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			log.Printf("stop publisher worker. ctx done %v\n", ctx.Err())
			return
		case <-timer.C:

			msgs, err := w.outboxRepository.GetReadyMessagesForPublish(ctx, msgsBatchSize)
			if err != nil {
				log.Printf("error get ready messages for publish %v\n", err)
				timer.Reset(w.cfg.MqPublishMessagesInterval())
				break
			}

			for _, msg := range msgs {
				if err := w.publisher.Publish(msg.SubjectName, msg.Payload); err != nil {
					log.Printf("failed to publish message %s into %s: %v", msg.Uid, msg.SubjectName, err)

					msg.UpdatedAt = time.Now().UTC()
					msg.LastErrorMessage = err.Error()
				} else {
					msg.SendAt = time.Now().UTC()
					log.Printf("successfuly send message %s into %s", msg.Uid, msg.SubjectName)
				}

				if err := w.outboxRepository.UpdateMessage(ctx, msg); err != nil {
					log.Printf("failed to update message %s: %v", msg.Uid, err)
				}
			}

			timer.Reset(w.cfg.MqPublishMessagesInterval())
		}
	}
}
