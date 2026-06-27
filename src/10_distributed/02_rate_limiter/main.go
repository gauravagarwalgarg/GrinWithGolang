/*
Module 10: Distributed - Rate Limiter

Demonstrates:
  - Token bucket algorithm: refill tokens at fixed rate, consume on request
  - Thread-safe implementation with sync.Mutex
  - HTTP middleware pattern for rate limiting
  - Comparison: Fixed Window vs Sliding Window vs Token Bucket
  - Configurable rate and burst capacity

Algorithm comparison:
  Fixed Window:  Simple counter reset each interval. Bursty at boundaries.
  Sliding Window: Weighted count across windows. Smoother but more complex.
  Token Bucket:  Tokens refill at rate R, bucket holds max B. Allows bursts
                 up to B while maintaining average rate R. Best for APIs.

Key insight: Token bucket is preferred because it allows controlled bursting
(good UX) while enforcing average rate (protects backend). Used by AWS,
Stripe, and most API gateways.

Run: go run main.go
*/
package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// TokenBucket implements the token bucket rate limiting algorithm.
type TokenBucket struct {
	mu         sync.Mutex
	tokens     float64
	maxTokens  float64
	refillRate float64   // tokens per second
	lastRefill time.Time
}

// NewTokenBucket creates a rate limiter with given rate and burst capacity.
func NewTokenBucket(rate float64, burst int) *TokenBucket {
	return &TokenBucket{
		tokens:     float64(burst),
		maxTokens:  float64(burst),
		refillRate: rate,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request is allowed (consumes one token).
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// Refill tokens based on elapsed time
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.tokens += elapsed * tb.refillRate
	if tb.tokens > tb.maxTokens {
		tb.tokens = tb.maxTokens
	}
	tb.lastRefill = now

	// Try to consume a token
	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}
	return false
}

// --- Per-client rate limiting ---

type RateLimiterMiddleware struct {
	mu       sync.Mutex
	limiters map[string]*TokenBucket
	rate     float64
	burst    int
}

func NewRateLimiterMiddleware(rate float64, burst int) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		limiters: make(map[string]*TokenBucket),
		rate:     rate,
		burst:    burst,
	}
}

func (rl *RateLimiterMiddleware) getLimiter(key string) *TokenBucket {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	limiter, exists := rl.limiters[key]
	if !exists {
		limiter = NewTokenBucket(rl.rate, rl.burst)
		rl.limiters[key] = limiter
	}
	return limiter
}

// Middleware returns an HTTP middleware that rate-limits by client IP.
func (rl *RateLimiterMiddleware) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		limiter := rl.getLimiter(clientIP)

		if !limiter.Allow() {
			w.Header().Set("Retry-After", "1")
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			log.Printf("[rate-limit] blocked request from %s", clientIP)
			return
		}
		next(w, r)
	}
}

func main() {
	// Demo 1: Direct token bucket usage
	fmt.Println("=== Token Bucket Demo ===")
	limiter := NewTokenBucket(5, 3) // 5 tokens/sec, burst of 3

	for i := 0; i < 7; i++ {
		allowed := limiter.Allow()
		fmt.Printf("  Request %d: allowed=%v\n", i+1, allowed)
	}

	// Wait for refill
	fmt.Println("\n  (waiting 1 second for refill...)")
	time.Sleep(1 * time.Second)

	for i := 0; i < 3; i++ {
		allowed := limiter.Allow()
		fmt.Printf("  Request %d: allowed=%v\n", i+8, allowed)
	}

	// Demo 2: HTTP middleware
	fmt.Println("\n=== HTTP Rate Limiter Middleware ===")
	rl := NewRateLimiterMiddleware(2, 5) // 2 req/sec, burst 5

	handler := rl.Middleware(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK\n"))
	})

	http.HandleFunc("/api", handler)
	log.Println("Rate-limited server on :8080 (2 req/sec, burst 5)")
	log.Println("Test: for i in {1..10}; do curl -s -o /dev/null -w '%{http_code}\\n' localhost:8080/api; done")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
