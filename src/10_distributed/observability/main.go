// Observability: Prometheus metrics + Tempo distributed tracing + structured logging.
//
// LEARNING NOTES:
// - Three Pillars of Observability: Metrics, Traces, Logs
// - Prometheus: pull-based metrics (counters, histograms, gauges) scraped at /metrics
// - Tempo: distributed tracing backend (stores traces, correlates spans)
// - OpenTelemetry (OTEL): vendor-neutral instrumentation SDK for traces + metrics
// - Trace = request journey across services; Span = single operation within a trace
//
// For C++ devs: Think of traces as flame graphs that span multiple processes/machines.
// Prometheus is like perf counters exposed via HTTP.
//
// Run: go run main.go
// Dashboard: Grafana at localhost:3000 (provisioned via docker-compose)
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// --- Prometheus Metrics ---

var (
	// Counter: total requests (only goes up)
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// Histogram: request duration distribution (p50, p90, p99)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets, // .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10
		},
		[]string{"method", "endpoint"},
	)

	// Gauge: currently in-flight requests (goes up and down)
	httpInFlightRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_in_flight_requests",
			Help: "Number of HTTP requests currently in-flight",
		},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration, httpInFlightRequests)
}

// --- OpenTelemetry Tracing Setup ---

func initTracer() (*sdktrace.TracerProvider, error) {
	// Export traces to Tempo via OTLP HTTP (default port 4318)
	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("observability-demo"),
			semconv.ServiceVersion("1.0.0"),
			attribute.String("environment", "development"),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}

// --- HTTP Handlers with Instrumentation ---

func orderHandler(w http.ResponseWriter, r *http.Request) {
	// Start a trace span
	ctx := r.Context()
	tracer := otel.Tracer("order-service")
	ctx, span := tracer.Start(ctx, "handle-order",
		trace.WithAttributes(attribute.String("http.method", r.Method)),
	)
	defer span.End()

	// Track in-flight
	httpInFlightRequests.Inc()
	defer httpInFlightRequests.Dec()

	start := time.Now()

	// Simulate processing with sub-spans
	validateOrder(ctx, tracer)
	processPayment(ctx, tracer)

	duration := time.Since(start).Seconds()

	// Record metrics
	httpRequestsTotal.WithLabelValues(r.Method, "/order", "200").Inc()
	httpRequestDuration.WithLabelValues(r.Method, "/order").Observe(duration)

	span.SetAttributes(attribute.Float64("duration_ms", duration*1000))

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","trace_id":"%s"}`, span.SpanContext().TraceID())
}

func validateOrder(ctx context.Context, tracer trace.Tracer) {
	_, span := tracer.Start(ctx, "validate-order")
	defer span.End()

	// Simulate validation work
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
	span.SetAttributes(attribute.Bool("valid", true))
}

func processPayment(ctx context.Context, tracer trace.Tracer) {
	_, span := tracer.Start(ctx, "process-payment")
	defer span.End()

	// Simulate payment gateway call
	time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)
	span.SetAttributes(attribute.String("gateway", "stripe"))
}

func main() {
	// Initialize tracing
	tp, err := initTracer()
	if err != nil {
		log.Printf("Warning: tracing disabled (%v)", err)
	} else {
		defer tp.Shutdown(context.Background())
	}

	// Prometheus metrics endpoint (Prometheus scrapes this)
	http.Handle("/metrics", promhttp.Handler())

	// Application endpoint (instrumented)
	http.HandleFunc("/order", orderHandler)

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"healthy"}`)
	})

	fmt.Println("Server running on :8081")
	fmt.Println("  /metrics  → Prometheus scrapes here")
	fmt.Println("  /order    → Instrumented endpoint")
	fmt.Println("  /health   → Health check")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
