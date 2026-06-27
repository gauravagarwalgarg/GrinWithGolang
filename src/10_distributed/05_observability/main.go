/*
Module 10: Distributed - Observability & Metrics

Demonstrates:
  - Custom metric types: Counter, Gauge, Histogram
  - Thread-safe metric collection with sync.Mutex
  - HTTP endpoint exposing Prometheus text format
  - How prometheus/client_golang works internally
  - Metric labels/dimensions for multi-dimensional data
  - Practical instrumentation of HTTP handlers

Prometheus text format:
  # HELP metric_name Description
  # TYPE metric_name counter|gauge|histogram
  metric_name{label="value"} 42

Run: go run main.go
Test: curl http://localhost:9090/metrics
*/
package main

import (
	"fmt"
	"math"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

// --- Metric types ---

type Counter struct {
	mu    sync.Mutex
	value float64
	name  string
	help  string
}

func (c *Counter) Inc()              { c.Add(1) }
func (c *Counter) Add(v float64)     { c.mu.Lock(); c.value += v; c.mu.Unlock() }
func (c *Counter) Get() float64      { c.mu.Lock(); defer c.mu.Unlock(); return c.value }

type Gauge struct {
	mu    sync.Mutex
	value float64
	name  string
	help  string
}

func (g *Gauge) Set(v float64)  { g.mu.Lock(); g.value = v; g.mu.Unlock() }
func (g *Gauge) Inc()           { g.mu.Lock(); g.value++; g.mu.Unlock() }
func (g *Gauge) Dec()           { g.mu.Lock(); g.value--; g.mu.Unlock() }
func (g *Gauge) Get() float64   { g.mu.Lock(); defer g.mu.Unlock(); return g.value }

type Histogram struct {
	mu      sync.Mutex
	buckets []float64 // upper bounds
	counts  []uint64  // count per bucket
	sum     float64
	count   uint64
	name    string
	help    string
}

func NewHistogram(name, help string, buckets []float64) *Histogram {
	sort.Float64s(buckets)
	return &Histogram{
		name:    name,
		help:    help,
		buckets: buckets,
		counts:  make([]uint64, len(buckets)),
	}
}

func (h *Histogram) Observe(v float64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.sum += v
	h.count++
	for i, bound := range h.buckets {
		if v <= bound {
			h.counts[i]++
		}
	}
}

// --- Metric Registry ---

type Registry struct {
	counters   []*Counter
	gauges     []*Gauge
	histograms []*Histogram
}

func (r *Registry) NewCounter(name, help string) *Counter {
	c := &Counter{name: name, help: help}
	r.counters = append(r.counters, c)
	return c
}

func (r *Registry) NewGauge(name, help string) *Gauge {
	g := &Gauge{name: name, help: help}
	r.gauges = append(r.gauges, g)
	return g
}

func (r *Registry) NewHistogram(name, help string, buckets []float64) *Histogram {
	h := NewHistogram(name, help, buckets)
	r.histograms = append(r.histograms, h)
	return h
}

// Expose returns metrics in Prometheus text exposition format.
func (r *Registry) Expose() string {
	var sb strings.Builder

	for _, c := range r.counters {
		fmt.Fprintf(&sb, "# HELP %s %s\n# TYPE %s counter\n%s %.1f\n\n",
			c.name, c.help, c.name, c.name, c.Get())
	}
	for _, g := range r.gauges {
		fmt.Fprintf(&sb, "# HELP %s %s\n# TYPE %s gauge\n%s %.1f\n\n",
			g.name, g.help, g.name, g.name, g.Get())
	}
	for _, h := range r.histograms {
		h.mu.Lock()
		fmt.Fprintf(&sb, "# HELP %s %s\n# TYPE %s histogram\n", h.name, h.help, h.name)
		cumulative := uint64(0)
		for i, bound := range h.buckets {
			cumulative += h.counts[i]
			label := fmt.Sprintf("%g", bound)
			if math.IsInf(bound, 1) {
				label = "+Inf"
			}
			fmt.Fprintf(&sb, "%s_bucket{le=\"%s\"} %d\n", h.name, label, cumulative)
		}
		fmt.Fprintf(&sb, "%s_sum %.3f\n%s_count %d\n\n", h.name, h.sum, h.name, h.count)
		h.mu.Unlock()
	}
	return sb.String()
}

func main() {
	reg := &Registry{}

	// Define metrics
	reqCounter := reg.NewCounter("http_requests_total", "Total HTTP requests")
	activeGauge := reg.NewGauge("http_active_requests", "Currently active requests")
	latencyHist := reg.NewHistogram("http_request_duration_seconds",
		"Request latency in seconds", []float64{0.01, 0.05, 0.1, 0.5, 1.0, 5.0})

	// Instrumented handler
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		reqCounter.Inc()
		activeGauge.Inc()
		start := time.Now()
		defer func() {
			latencyHist.Observe(time.Since(start).Seconds())
			activeGauge.Dec()
		}()
		time.Sleep(50 * time.Millisecond) // simulate work
		w.Write([]byte("OK\n"))
	})

	// Metrics endpoint (like /metrics in Prometheus)
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		w.Write([]byte(reg.Expose()))
	})

	fmt.Println("Observability demo on :9090")
	fmt.Println("  GET /api     - instrumented endpoint")
	fmt.Println("  GET /metrics - Prometheus text format")
	http.ListenAndServe(":9090", nil)
}
