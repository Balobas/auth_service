// Package metrics — Prometheus-метрики auth_service: технические (RED) по gRPC и
// продуктовые гейджи (users/DAU/MAU). Экспонируются на /metrics.
package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	// ── Технические (RED) ──────────────────────────────────────────────────
	grpcRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_requests_total",
		Help: "Total number of gRPC requests by method and status code.",
	}, []string{"method", "code"})

	grpcDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "grpc_request_duration_seconds",
		Help:    "gRPC request latency in seconds by method.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method"})

	grpcInflight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "grpc_requests_in_flight",
		Help: "Number of gRPC requests currently being served.",
	})

	// ── Продуктовые ────────────────────────────────────────────────────────
	// Значения выставляет периодический джоб (см. StatsUpdater): это агрегаты
	// по БД, а не per-request счётчики (DAU/MAU — уникальные пользователи).
	UsersTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "auth_users_total",
		Help: "Total number of registered users.",
	})

	UsersByMethod = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "auth_users_by_method",
		Help: "Registered users by auth method (email/telegram).",
	}, []string{"method"})

	DAU = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "auth_dau",
		Help: "Daily active users — distinct users active today (UTC).",
	})

	MAU = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "auth_mau",
		Help: "Monthly active users — distinct users active in the last 30 days.",
	})

	NewUsersToday = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "auth_new_users_today",
		Help: "Users registered today (UTC).",
	})
)

// UnaryServerInterceptor считает количество, длительность и in-flight gRPC-запросов.
// Ставится ПЕРВЫМ в цепочке (внешним), чтобы видеть итоговый status-код после
// errors-post-интерцептора.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		grpcInflight.Inc()
		start := time.Now()

		resp, err := handler(ctx, req)

		grpcInflight.Dec()
		grpcDuration.WithLabelValues(info.FullMethod).Observe(time.Since(start).Seconds())
		grpcRequests.WithLabelValues(info.FullMethod, status.Code(err).String()).Inc()

		return resp, err
	}
}

// Handler — http.Handler для эндпоинта /metrics.
func Handler() http.Handler {
	return promhttp.Handler()
}
