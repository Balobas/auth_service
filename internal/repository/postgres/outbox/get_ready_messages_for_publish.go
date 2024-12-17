package outboxRepository

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/balobas/auth_service/internal/entity"
	pgEntity "github.com/balobas/auth_service/internal/repository/postgres/pg_entity"
	"github.com/pkg/errors"
)

func (r *OutboxRepository) GetReadyMessagesForPublish(ctx context.Context, batchSize int64) ([]entity.MqMessage, error) {
	log.Printf("outboxRepository.GetReadyMessagesForPublish: batch size %d", batchSize)

	msgRow := pgEntity.NewMqMessageRow()

	sql, args, err := sq.Select(
		msgRow.Columns()...,
	).From(
		msgRow.Table(),
	).PlaceholderFormat(
		sq.Dollar,
	).Where(
		msgRow.ConditionSendAtIsNull(),
	).OrderBy("created_at").Limit(uint64(batchSize)).ToSql()
	if err != nil {
		log.Printf("failed to build sql query for GetReadyMessagesForPublish: %v", err)
		return nil, errors.WithStack(err)
	}

	log.Println(sql)

	rows, err := r.DB().Query(ctx, sql, args...)
	if err != nil {
		log.Printf("failed to get ready messages for publish: %v", err)
		return nil, errors.WithStack(err)
	}

	msgRows := pgEntity.NewMqMessageRows()
	if err := msgRows.ScanAll(rows); err != nil {
		log.Printf("failed to scan ready messages for publish: %v", err)
		return nil, errors.WithStack(err)
	}

	return msgRows.ToEntity(), nil
}
