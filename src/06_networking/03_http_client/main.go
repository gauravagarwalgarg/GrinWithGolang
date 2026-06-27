/*
Module 6: Networking - HTTP Client Best Practices

Demonstrates:
  - Custom http.Transport with connection pooling and timeouts
  - Reusable http.Client (never use http.DefaultClient in production)
  - Retry with exponential backoff
  - Context cancellation for request deadlines
  - Safe JSON response reading with size limits
  - Proper resource cleanup (resp.Body.Close)

Key insight: http.DefaultClient has no timeout a production footgun.
Always create a client with explicit timeouts and transport settings.

Run: go run main.go
*/
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// newHTTPClient creates a production-ready HTTP client with pooling and timeouts.
func newHTTPClient() *http.Client {
	transport := &http.Transport{
		MaxIdleConns:        100,              // connection pool size
		MaxIdleConnsPerHost: 10,               // per-host pool
		IdleConnTimeout:     90 * time.Second, // reclaim idle connections
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // overall request timeout
	}
}

// doWithRetry executes a request with exponential backoff retry.
func doWithRetry(ctx context.Context, client *http.Client, req *http.Request, maxRetries int) (*http.Response, error) {
	var lastErr error
	backoff := 100 * time.Millisecond

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			log.Printf("retry attempt %d/%d after %v", attempt, maxRetries, backoff)
			select {
			case <-time.After(backoff):
				backoff *= 2 // exponential backoff
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		// Retry on 5xx server errors
		if resp.StatusCode >= 500 {
			resp.Body.Close()
			lastErr = fmt.Errorf("server error: %d", resp.StatusCode)
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("all %d retries exhausted: %w", maxRetries, lastErr)
}

// readJSON safely reads and decodes a JSON response with size limit.
func readJSON(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	// Limit response body to 1MB to prevent memory exhaustion
	limited := io.LimitReader(resp.Body, 1<<20)
	return json.NewDecoder(limited).Decode(target)
}

func main() {
	client := newHTTPClient()

	// Create a context with timeout for the entire operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://httpbin.org/json", nil)
	if err != nil {
		log.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "GrinWithGolang/1.0")

	// Execute with retry
	resp, err := doWithRetry(ctx, client, req, 3)
	if err != nil {
		log.Fatalf("request failed: %v", err)
	}

	// Decode response safely
	var result map[string]interface{}
	if err := readJSON(resp, &result); err != nil {
		log.Fatalf("failed to decode response: %v", err)
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Response: %+v\n", result)

	// Demonstrate context cancellation
	shortCtx, shortCancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer shortCancel()

	req2, _ := http.NewRequestWithContext(shortCtx, http.MethodGet,
		"https://httpbin.org/delay/5", nil)
	_, err = client.Do(req2)
	if err != nil {
		fmt.Printf("\nContext cancellation demo: %v\n", err)
	}
}
