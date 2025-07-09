package metrics

import (
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Collectors struct {
	Reloaded            *prometheus.CounterVec
	ReloadedByNamespace *prometheus.CounterVec
	QueueSize           *prometheus.GaugeVec
	Errors              *prometheus.CounterVec
	Requeues            *prometheus.CounterVec
	Dropped             *prometheus.CounterVec
	Added               *prometheus.CounterVec
	Updated             *prometheus.CounterVec
	Deleted             *prometheus.CounterVec
	ActionTime          *prometheus.GaugeVec
}

func NewCollectors() Collectors {
	reloaded := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "reloader",
			Name:      "reload_executed_total",
			Help:      "Counter of reloads executed by Reloader.",
		},
		[]string{
			"success",
		},
	)

	//set 0 as default value
	reloaded.With(prometheus.Labels{"success": "true"}).Add(0)
	reloaded.With(prometheus.Labels{"success": "false"}).Add(0)

	reloaded_by_namespace := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "reloader",
			Name:      "reload_executed_total_by_namespace",
			Help:      "Counter of reloads executed by Reloader by namespace.",
		},
		[]string{
			"success",
			"namespace",
		},
	)

	queueSize := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "reloader",
			Name:      "queue_size",
			Help:      "Gauge for the size of the work queue.",
		},
		[]string{
			"resource",
		},
	)

	errors := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "reloader",
			Name:      "errors_total",
			Help:      "Counter of errors encountered by Reloader.",
		},
		[]string{
			"error_type",
		},
	)
	requeues := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "reloader",
			Name:      "requeues_total",
			Help:      "Counter of requeues encountered by Reloader.",
		},
		[]string{
			"resource",
		},
	)

	dropped := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "reloader",
			Name:      "dropped_total",
			Help:      "Counter of dropped events by Reloader.",
		},
		[]string{
			"resource",
		},
	)

	added := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "reloader",
			Name:      "added_total",
			Help:      "Counter of resources added to the queue by Reloader.",
		},
		[]string{
			"resource",
		},
	)
	updated := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "reloader",
			Name:      "updated_total",
			Help:      "Counter of resources updated in the queue by Reloader.",
		},
		[]string{
			"resource",
		},
	)
	deleted := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "reloader",
			Name:      "deleted_total",
			Help:      "Counter of resources deleted from the queue by Reloader.",
		},
		[]string{
			"resource",
		},
	)
	actionTime := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "reloader",
			Name:      "action_time_seconds",
			Help:      "Gauge for the time taken to perform actions in seconds.",
		},
		[]string{
			"resource",
		},
	)

	return Collectors{
		Reloaded:            reloaded,
		ReloadedByNamespace: reloaded_by_namespace,
		QueueSize:           queueSize,
		Errors:              errors,
		Requeues:            requeues,
		Dropped:             dropped,
		Added:               added,
		Updated:             updated,
		Deleted:             deleted,
		ActionTime:          actionTime,
	}
}

func SetupPrometheusEndpoint() Collectors {
	collectors := NewCollectors()
	prometheus.MustRegister(collectors.Reloaded)

	if os.Getenv("METRICS_COUNT_BY_NAMESPACE") == "enabled" {
		prometheus.MustRegister(collectors.ReloadedByNamespace)
	}

	prometheus.MustRegister(collectors.QueueSize)
	prometheus.MustRegister(collectors.Errors)
	prometheus.MustRegister(collectors.Requeues)
	prometheus.MustRegister(collectors.Dropped)
	prometheus.MustRegister(collectors.Added)
	prometheus.MustRegister(collectors.Updated)
	prometheus.MustRegister(collectors.Deleted)
	prometheus.MustRegister(collectors.ActionTime)

	http.Handle("/metrics", promhttp.Handler())

	return collectors
}
