/*
What this teaches:
    The Functional Options pattern THE most important Go pattern for library
    design. Avoids constructor explosion, provides zero-value defaults, and allows
    optional configuration without breaking API compatibility.

Beginner analogy:
    "Ordering coffee: 'Give me a coffee' (defaults). 'With oat milk, extra shot,
     no sugar' (options). You don't need a different counter for every combination."

C++ comparison:
    "In C++ you'd use builder pattern or parameter structs with optional<T>. Go's
     functional options are closures that mutate a config cleaner than builders
     and more extensible than config structs."

Interview relevance:
    This pattern appears in production libraries (gRPC, Zap, fx). Interviewers ask
    you to design a configurable server/client functional options is the expected
    answer for senior Go roles.
*/

package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

// --- The Server we want to configure ---

type Server struct {
	host         string
	port         int
	timeout      time.Duration
	maxConns     int
	tls          bool
	logger       *log.Logger
}

// Option is a function that configures the server
type Option func(*Server)

// --- Option constructors (the "With" pattern) ---

func WithPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

func WithHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}

func WithTimeout(d time.Duration) Option {
	return func(s *Server) {
		s.timeout = d
	}
}

func WithMaxConns(n int) Option {
	return func(s *Server) {
		s.maxConns = n
	}
}

func WithTLS(enabled bool) Option {
	return func(s *Server) {
		s.tls = enabled
	}
}

func WithLogger(l *log.Logger) Option {
	return func(s *Server) {
		s.logger = l
	}
}

// --- Constructor with sensible defaults ---

func NewServer(opts ...Option) *Server {
	// Step 1: Set defaults
	s := &Server{
		host:     "localhost",
		port:     8080,
		timeout:  30 * time.Second,
		maxConns: 100,
		tls:      false,
		logger:   log.New(os.Stdout, "[server] ", log.LstdFlags),
	}

	// Step 2: Apply each option
	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Start() {
	protocol := "http"
	if s.tls {
		protocol = "https"
	}
	s.logger.Printf("Starting on %s://%s:%d", protocol, s.host, s.port)
	s.logger.Printf("Timeout: %v | MaxConns: %d | TLS: %v", s.timeout, s.maxConns, s.tls)
}

// --- Another example: HTTP Client ---

type HTTPClient struct {
	baseURL    string
	retries    int
	timeout    time.Duration
	headers    map[string]string
}

type ClientOption func(*HTTPClient)

func WithBaseURL(url string) ClientOption {
	return func(c *HTTPClient) { c.baseURL = url }
}

func WithRetries(n int) ClientOption {
	return func(c *HTTPClient) { c.retries = n }
}

func WithClientTimeout(d time.Duration) ClientOption {
	return func(c *HTTPClient) { c.timeout = d }
}

func WithHeader(key, value string) ClientOption {
	return func(c *HTTPClient) { c.headers[key] = value }
}

func NewHTTPClient(opts ...ClientOption) *HTTPClient {
	c := &HTTPClient{
		baseURL: "http://localhost",
		retries: 3,
		timeout: 10 * time.Second,
		headers: make(map[string]string),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func main() {
	fmt.Println("=== Functional Options Pattern ===")

	// 1. Default server zero config needed
	fmt.Println("\n--- Default Server ---")
	s1 := NewServer()
	s1.Start()

	// 2. Customized server only set what you need
	fmt.Println("\n--- Custom Server ---")
	s2 := NewServer(
		WithPort(9443),
		WithHost("0.0.0.0"),
		WithTLS(true),
		WithTimeout(60*time.Second),
		WithMaxConns(1000),
	)
	s2.Start()

	// 3. HTTP client with options
	fmt.Println("\n--- HTTP Client ---")
	client := NewHTTPClient(
		WithBaseURL("https://api.example.com"),
		WithRetries(5),
		WithHeader("Authorization", "Bearer token123"),
		WithHeader("Content-Type", "application/json"),
	)
	fmt.Printf("  Base URL: %s\n", client.baseURL)
	fmt.Printf("  Retries: %d\n", client.retries)
	fmt.Printf("  Headers: %v\n", client.headers)

	// 4. Key takeaways
	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. type Option func(*T) closure that mutates config")
	fmt.Println("2. Constructor: NewT(opts ...Option) applies defaults then options")
	fmt.Println("3. Adding new options never breaks existing callers")
	fmt.Println("4. Self-documenting: WithPort(9443) reads like English")
	fmt.Println("5. Used in: gRPC, Zap, go-kit, uber-fx, and most Go libraries")
}
