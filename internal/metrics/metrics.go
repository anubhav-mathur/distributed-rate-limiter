package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_limiter_requests_total",
			Help: "Total number of requests received",
		},
		[]string{"user", "status"},
	)

	RedisLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "redis_latency_seconds",
			Help:    "Latency of Redis operations",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func Init() {
	prometheus.MustRegister(RequestsTotal)
	prometheus.MustRegister(RedisLatency)
}
