/*
What this teaches:
    Exponential backoff with jitter for retrying transient failures. Configurable
    max retries, initial/max delay, and wrapping any function with retry logic.
    Essential for robust network clients.

Beginner analogy:
    "Like redialing a busy phone number: wait 1 second, then 2, then 4... plus a
     random wobble (jitter) so everyone isn't redialing at the exact same moment."

C++ comparison:
    "Similar to AWS SDK's retry strategies or gRPC's backoff policies. In C++ you'd
     use a template with std::function; in Go, closures and functional configuration
     make this cleaner."

Interview relevance:
    Interviewers ask: Why jitter? (thundering herd avoidance). Why cap the backoff?
    (prevent absurd delays). How to make it context-aware? (respect deadlines).
    This covers all three.
*/

package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// --- Retry configuration using functional options ---

type RetryConfig struct {
	MaxRetries   int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
	JitterFactor float64 // 0.0 = no jitter, 1.0 = full jitter
}

type RetryOption func(*RetryConfig)

func WithMaxRetries(n int) RetryOption {
	return func(c *RetryConfig) { c.MaxRetries = n }
}

func WithInitialDelay(d time.Duration) RetryOption {
	return func(c *RetryConfig) { c.InitialDelay = d }
}

func WithMaxDelay(d time.Duration) RetryOption {
	return func(c *RetryConfig) { c.MaxDelay = d }
}

func WithMultiplier(m float64) RetryOption {
	return func(c *RetryConfig) { c.Multiplier = m }
}

func WithJitter(j float64) RetryOption {
	return func(c *RetryConfig) { c.JitterFactor = j }
}

func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:   5,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     10 * time.Second,
		Multiplier:   2.0,
		JitterFactor: 0.5,
	}
}

// --- Core retry function ---

func Retry(ctx context.Context, fn func() error, opts ...RetryOption) error {
	cfg := DefaultRetryConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	var lastErr error
	delay := cfg.InitialDelay

	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		if attempt > 0 {
			// Apply jitter: delay ± (jitterFactor * delay * random)
			jitter := time.Duration(float64(delay) * cfg.JitterFactor * (rand.Float64()*2 - 1))
			actualDelay := delay + jitter
			if actualDelay < 0 {
				actualDelay = 0
			}
			fmt.Printf("  [retry] Attempt %d/%d waiting %v\n", attempt, cfg.MaxRetries, actualDelay)

			select {
			case <-time.After(actualDelay):
			case <-ctx.Done():
				return fmt.Errorf("retry cancelled: %w", ctx.Err())
			}

			// Exponential increase with cap
			delay = time.Duration(float64(delay) * cfg.Multiplier)
			if delay > cfg.MaxDelay {
				delay = cfg.MaxDelay
			}
		}

		lastErr = fn()
		if lastErr == nil {
			if attempt > 0 {
				fmt.Printf("  [retry] Succeeded on attempt %d\n", attempt+1)
			}
			return nil
		}

		// Check if error is retryable
		if !isRetryable(lastErr) {
			return fmt.Errorf("non-retryable error: %w", lastErr)
		}
	}

	return fmt.Errorf("max retries (%d) exceeded: %w", cfg.MaxRetries, lastErr)
}

// --- Retryable error detection ---

type retryableError struct {
	err error
}

func (e *retryableError) Error() string { return e.err.Error() }
func (e *retryableError) Unwrap() error { return e.err }
func (e *retryableError) IsRetryable() bool { return true }

func NewRetryableError(msg string) error {
	return &retryableError{err: errors.New(msg)}
}

func isRetryable(err error) bool {
	type retryable interface{ IsRetryable() bool }
	var r retryable
	if errors.As(err, &r) {
		return r.IsRetryable()
	}
	return true // Default: assume retryable for unknown errors
}

// --- Compute backoff for display ---

func computeBackoffSchedule(cfg *RetryConfig) {
	fmt.Println("  Delay schedule (no jitter):")
	delay := cfg.InitialDelay
	for i := 1; i <= cfg.MaxRetries; i++ {
		fmt.Printf("    Attempt %d: %v\n", i, delay)
		delay = time.Duration(math.Min(
			float64(delay)*cfg.Multiplier,
			float64(cfg.MaxDelay),
		))
	}
}

// --- Simulated flaky service ---

func flakyService(failCount *int, maxFails int) func() error {
	return func() error {
		*failCount++
		if *failCount <= maxFails {
			return NewRetryableError(fmt.Sprintf("connection refused (attempt %d)", *failCount))
		}
		return nil
	}
}

func main() {
	fmt.Println("=== Retry with Exponential Backoff ===")
	rand.Seed(time.Now().UnixNano())

	// 1. Show backoff schedule
	fmt.Println("\n--- Backoff Schedule ---")
	computeBackoffSchedule(DefaultRetryConfig())

	// 2. Successful retry after transient failures
	fmt.Println("\n--- Retry: Success after 3 failures ---")
	ctx := context.Background()
	failCount := 0
	err := Retry(ctx, flakyService(&failCount, 3),
		WithMaxRetries(5),
		WithInitialDelay(50*time.Millisecond),
		WithJitter(0.3),
	)
	fmt.Printf("  Result: err=%v\n", err)

	// 3. Exhausted retries
	fmt.Println("\n--- Retry: Max retries exhausted ---")
	failCount = 0
	err = Retry(ctx, flakyService(&failCount, 100),
		WithMaxRetries(3),
		WithInitialDelay(20*time.Millisecond),
	)
	fmt.Printf("  Result: %v\n", err)

	// 4. Context cancellation
	fmt.Println("\n--- Retry: Context cancelled ---")
	cancelCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	failCount = 0
	err = Retry(cancelCtx, flakyService(&failCount, 100),
		WithInitialDelay(60*time.Millisecond),
	)
	fmt.Printf("  Result: %v\n", err)

	// Key takeaways
	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. Exponential backoff: delay *= multiplier each attempt")
	fmt.Println("2. Jitter prevents thundering herd (all clients retrying simultaneously)")
	fmt.Println("3. MaxDelay caps prevent absurd wait times")
	fmt.Println("4. Context-aware: respects deadlines and cancellation")
	fmt.Println("5. Error classification: only retry retryable errors")
}
