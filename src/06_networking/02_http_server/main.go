/*
Module 6: Networking - HTTP Server with net/http

Demonstrates:
  - Route registration with http.HandleFunc
  - Middleware pattern (logging, auth) via function composition
  - JSON responses with proper Content-Type headers
  - Structured error handling with HTTP status codes
  - Graceful shutdown with http.Server + os/signal
  - Context propagation through request lifecycle

Key insight: Go's http.Handler interface (single method: ServeHTTP) enables
composable middleware without frameworks. Compare to Express.js middleware
or C++ Crow/Beast handler chains.

Run: go run main.go
Test: curl http://localhost:8080/api/users
*/
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Middleware: logging - wraps handler to log request details
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("→ %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
		log.Printf("← %s %s completed in %v", r.Method, r.URL.Path, time.Since(start))
	}
}

// Middleware: auth - checks for API key in header
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "missing X-API-Key header",
			})
			return
		}
		// In production: validate against a store
		next(w, r)
	}
}

// writeJSON marshals data and writes it as a JSON response.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("json encode error: %v", err)
	}
}

// Handlers
func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}
	users := []map[string]string{
		{"id": "1", "name": "Alice", "role": "engineer"},
		{"id": "2", "name": "Bob", "role": "designer"},
	}
	writeJSON(w, http.StatusOK, users)
}

func handleProtected(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"message": "you have access to protected resource",
	})
}

func main() {
	mux := http.NewServeMux()

	// Public routes with logging
	mux.HandleFunc("/health", loggingMiddleware(handleHealth))
	mux.HandleFunc("/api/users", loggingMiddleware(handleUsers))

	// Protected route: logging + auth middleware stacked
	mux.HandleFunc("/api/admin", loggingMiddleware(authMiddleware(handleProtected)))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in background
	go func() {
		log.Println("HTTP server listening on :8080")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown on signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	log.Println("server stopped gracefully")
}
