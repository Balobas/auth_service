// Package repositoryMetrics — доступ к данным для продуктовых метрик
// (запись активности для DAU/MAU и сбор агрегатов).
package repositoryMetrics

import (
	"context"

	"github.com/balobas/auth_service/pkg/metrics"
	DBclient "github.com/balobas/sport_city_common/clients/database"
	repositoryPostgres "github.com/balobas/sport_city_common/repository/postgres"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type MetricsRepository struct {
	*repositoryPostgres.BasePgRepository
}

func New(client DBclient.ClientDB) *MetricsRepository {
	return &MetricsRepository{repositoryPostgres.New(client)}
}

// RecordActivity отмечает пользователя активным сегодня (UTC). Идемпотентно:
// ON CONFLICT DO NOTHING — одна строка на (user, день).
func (r *MetricsRepository) RecordActivity(ctx context.Context, userUid uuid.UUID) error {
	const stmt = `
		INSERT INTO user_activity (user_uid, active_date)
		VALUES ($1, (now() at time zone 'utc')::date)
		ON CONFLICT DO NOTHING`
	if _, err := r.Exec(ctx, stmt, userUid.String()); err != nil {
		return errors.Wrap(err, "RecordActivity")
	}
	return nil
}

// CollectProductStats считает агрегаты одним запросом.
func (r *MetricsRepository) CollectProductStats(ctx context.Context) (metrics.ProductStats, error) {
	const stmt = `
		SELECT
			(SELECT count(*) FROM users)                                             AS users_total,
			(SELECT count(*) FROM users WHERE telegram_id IS NOT NULL)               AS users_tg,
			(SELECT count(*) FROM user_activity
				WHERE active_date = (now() at time zone 'utc')::date)                AS dau,
			(SELECT count(DISTINCT user_uid) FROM user_activity
				WHERE active_date > (now() at time zone 'utc')::date - 30)           AS mau,
			(SELECT count(*) FROM users
				WHERE created_at::date = (now() at time zone 'utc')::date)           AS new_today`

	var s metrics.ProductStats
	var usersTotal, usersTg, dau, mau, newToday int64
	if err := r.QueryRow(ctx, stmt).Scan(&usersTotal, &usersTg, &dau, &mau, &newToday); err != nil {
		return metrics.ProductStats{}, errors.Wrap(err, "CollectProductStats")
	}

	s.UsersTotal = usersTotal
	s.UsersTelegram = usersTg
	s.UsersEmail = usersTotal - usersTg // без привязанного Telegram
	s.DAU = dau
	s.MAU = mau
	s.NewUsersToday = newToday
	return s, nil
}
