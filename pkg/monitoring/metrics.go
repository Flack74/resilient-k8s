package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all the Prometheus metrics for the application
type Metrics struct {
	ExperimentsCreated    prometheus.Counter
	ExperimentsExecuted   prometheus.Counter
	ExperimentsSucceeded  prometheus.Counter
	ExperimentsFailed     prometheus.Counter
	ExperimentDuration    prometheus.Histogram
	ActiveExperiments     prometheus.Gauge
	TargetsAffected       prometheus.Counter
	APIRequestsTotal      *prometheus.CounterVec
	APIRequestDuration    *prometheus.HistogramVec
}

// NewMetrics creates and registers Prometheus metrics
func NewMetrics() *Metrics {
	m := &Metrics{
		ExperimentsCreated: promauto.NewCounter(prometheus.CounterOpts{
			Name: "chaos_experiments_created_total",
			Help: "The total number of chaos experiments created",
		}),
		ExperimentsExecuted: promauto.NewCounter(prometheus.CounterOpts{
			Name: "chaos_experiments_executed_total",
			Help: "The total number of chaos experiments executed",
		}),
		ExperimentsSucceeded: promauto.NewCounter(prometheus.CounterOpts{
			Name: "chaos_experiments_succeeded_total",
			Help: "The total number of chaos experiments that succeeded",
		}),
		ExperimentsFailed: promauto.NewCounter(prometheus.CounterOpts{
			Name: "chaos_experiments_failed_total",
			Help: "The total number of chaos experiments that failed",
		}),
		ExperimentDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "chaos_experiment_duration_seconds",
			Help:    "The duration of chaos experiments in seconds",
			Buckets: prometheus.DefBuckets,
		}),
		ActiveExperiments: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "chaos_experiments_active",
			Help: "The number of currently active chaos experiments",
		}),
		TargetsAffected: promauto.NewCounter(prometheus.CounterOpts{
			Name: "chaos_targets_affected_total",
			Help: "The total number of targets affected by chaos experiments",
		}),
		APIRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "chaos_api_requests_total",
				Help: "The total number of API requests",
			},
			[]string{"method", "endpoint", "status"},
		),
		APIRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "chaos_api_request_duration_seconds",
				Help:    "The duration of API requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),
	}

	return m
}