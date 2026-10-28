// Package useCaseActivity — трекер дневной активности пользователей для DAU/MAU.
package useCaseActivity

import (
	"context"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Repo interface {
	RecordActivity(ctx context.Context, userUid uuid.UUID) error
}

// Tracker дедуплицирует отметки активности в памяти (одна запись в БД на юзера
// в день на инстанс) и пишет их в фоне, чтобы не влиять на латентность запроса.
// Набор активных сбрасывается на границе UTC-суток.
type Tracker struct {
	repo Repo

	mu   sync.Mutex
	day  string
	seen map[string]struct{}
}

func New(repo Repo) *Tracker {
	return &Tracker{repo: repo, seen: make(map[string]struct{})}
}

// Touch отмечает пользователя активным сегодня. Неблокирующе и best-effort:
// при ошибке записи метка снимается, чтобы следующий запрос повторил попытку.
func (t *Tracker) Touch(userUid uuid.UUID) {
	if t == nil {
		return
	}
	day := time.Now().UTC().Format("2006-01-02")
	key := userUid.String()

	t.mu.Lock()
	if day != t.day {
		t.day = day
		t.seen = make(map[string]struct{})
	}
	if _, ok := t.seen[key]; ok {
		t.mu.Unlock()
		return
	}
	t.seen[key] = struct{}{}
	t.mu.Unlock()

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := t.repo.RecordActivity(ctx, userUid); err != nil {
			t.mu.Lock()
			if t.day == day {
				delete(t.seen, key)
			}
			t.mu.Unlock()
		}
	}()
}
