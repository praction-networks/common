package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics defines Prometheus metrics with support for labeled metrics.
type Metrics struct {
	PublishedEvents   *prometheus.CounterVec
	ProcessedMessages *prometheus.CounterVec
	FailedMessages    *prometheus.CounterVec
	Duration          *prometheus.HistogramVec // Added duration metric
}

// NewMetrics initializes and returns a new Metrics instance with labeled counters and histograms.
func NewMetrics() *Metrics {
	return &Metrics{
		PublishedEvents: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "nats_published_events_total",
				Help: "Total number of published events, labeled by stream and event",
			},
			[]string{"stream", "event"}, // Labels: stream and event
		),
		ProcessedMessages: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "nats_processed_messages_total",
				Help: "Total number of processed messages, labeled by stream and event",
			},
			[]string{"stream", "event"}, // Labels: stream and event
		),
		FailedMessages: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "nats_failed_messages_total",
				Help: "Total number of failed messages, labeled by stream and event",
			},
			[]string{"stream", "event"}, // Labels: stream and event
		),
		Duration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "nats_message_processing_duration_seconds",
				Help:    "Histogram of message processing durations in seconds, labeled by stream and event",
				Buckets: prometheus.DefBuckets, // Default histogram buckets
			},
			[]string{"stream", "event"}, // Labels: stream and event
		),
	}
}

// Register registers the metrics with the Prometheus default registry.
func (m *Metrics) Register() {
	prometheus.MustRegister(m.PublishedEvents)
	prometheus.MustRegister(m.ProcessedMessages)
	prometheus.MustRegister(m.FailedMessages)
	prometheus.MustRegister(m.Duration) // Register the duration metric
}
