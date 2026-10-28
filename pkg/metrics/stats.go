package metrics

import (
	"context"
	"log"
	"time"
)

// ProductStats — агрегаты для продуктовых гейджей (считаются по БД).
type ProductStats struct {
	UsersTotal    int64
	UsersEmail    int64
	UsersTelegram int64
	DAU           int64
	MAU           int64
	NewUsersToday int64
}

// StatsSource отдаёт свежие продуктовые агрегаты (реализуется репозиторием).
type StatsSource interface {
	CollectProductStats(ctx context.Context) (ProductStats, error)
}

// StatsUpdater периодически пересчитывает продуктовые метрики и обновляет гейджи.
// DAU/MAU — уникальные пользователи, поэтому это агрегаты по БД, а не счётчики.
type StatsUpdater struct {
	src      StatsSource
	interval time.Duration
}

func NewStatsUpdater(src StatsSource, interval time.Duration) *StatsUpdater {
	if interval <= 0 {
		interval = time.Minute
	}
	return &StatsUpdater{src: src, interval: interval}
}

func (u *StatsUpdater) Run(ctx context.Context) {
	t := time.NewTicker(u.interval)
	defer t.Stop()

	u.update(ctx) // сразу на старте, чтобы гейджи не были пустыми
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			u.update(ctx)
		}
	}
}

func (u *StatsUpdater) update(ctx context.Context) {
	s, err := u.src.CollectProductStats(ctx)
	if err != nil {
		log.Printf("metrics: failed to collect product stats: %v", err)
		return
	}
	UsersTotal.Set(float64(s.UsersTotal))
	UsersByMethod.WithLabelValues("email").Set(float64(s.UsersEmail))
	UsersByMethod.WithLabelValues("telegram").Set(float64(s.UsersTelegram))
	DAU.Set(float64(s.DAU))
	MAU.Set(float64(s.MAU))
	NewUsersToday.Set(float64(s.NewUsersToday))
}
