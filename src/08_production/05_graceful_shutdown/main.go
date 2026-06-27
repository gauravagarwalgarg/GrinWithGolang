/*
Module 8: Production - Graceful Shutdown Pattern

Demonstrates:
  - Signal handling (SIGINT/SIGTERM) for container orchestration
  - Context cancellation cascade through all components
  - Draining in-flight HTTP requests
  - Closing database connections (simulated)
  - Ordered shutdown: stop accepting → drain → close resources
  - Production-ready pattern used by Go services at scale

Key insight: Kubernetes sends SIGTERM, then waits terminationGracePeriodSeconds
before SIGKILL. Your service must: stop accepting traffic, finish in-flight
requests, flush buffers, close connections all within the grace period.

Run: go run main.go
Test: kill -SIGTERM <pid> or Ctrl+C
*/
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// --- Simulated dependencies ---

type Database struct {
	mu   sync.Mutex
	open bool
}

func (db *Database) Connect() {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.open = true
	log.Println("[db] connected")
}

func (db *Database) Close() {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.open = false
	log.Println("[db] connection closed")
}

type MessageQueue struct {
	buffer []string
	mu     sync.Mutex
}

func (mq *MessageQueue) Flush() {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	log.Printf("[mq] flushing %d buffered messages", len(mq.buffer))
	mq.buffer = nil
}

func (mq *MessageQueue) Publish(msg string) {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	mq.buffer = append(mq.buffer, msg)
}

// --- Request tracking for drain ---

type RequestTracker struct {
	wg sync.WaitGroup
}

func (rt *RequestTracker) Track(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rt.wg.Add(1)
		defer rt.wg.Done()
		next(w, r)
	}
}

func (rt *RequestTracker) Wait(timeout time.Duration) bool {
	done := make(chan struct{})
	go func() {
		rt.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}

func main() {
	// Initialize dependencies
	db := &Database{}
	db.Connect()

	mq := &MessageQueue{}
	tracker := &RequestTracker{}

	// Simulate background work
	mq.Publish("startup event")

	// HTTP handler that simulates work
	mux := http.NewServeMux()
	mux.HandleFunc("/", tracker.Track(func(w http.ResponseWriter, r *http.Request) {
		mq.Publish(fmt.Sprintf("request: %s %s", r.Method, r.URL.Path))
		time.Sleep(100 * time.Millisecond) // simulate work
		w.Write([]byte("OK\n"))
	}))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start server
	go func() {
		log.Println("[server] listening on :8080")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("[server] error: %v", err)
		}
	}()

	// --- Graceful shutdown orchestration ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("[shutdown] received signal: %v", sig)

	// Phase 1: Stop accepting new requests
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	log.Println("[shutdown] phase 1: stop accepting connections")
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("[shutdown] server shutdown error: %v", err)
	}

	// Phase 2: Wait for in-flight requests to drain
	log.Println("[shutdown] phase 2: draining in-flight requests...")
	if drained := tracker.Wait(10 * time.Second); !drained {
		log.Println("[shutdown] WARNING: timed out waiting for requests to drain")
	} else {
		log.Println("[shutdown] all in-flight requests completed")
	}

	// Phase 3: Close resources in reverse-dependency order
	log.Println("[shutdown] phase 3: closing resources")
	mq.Flush()
	db.Close()

	log.Println("[shutdown] clean shutdown complete")
}
