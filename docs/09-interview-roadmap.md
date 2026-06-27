# Go Interview Preparation Roadmap

## 4-Week Study Plan

### Week 1: Foundations & Data Structures

| Day | Topic | Practice |
|-----|-------|----------|
| 1-2 | Go syntax, types, slices, maps | Implement custom hashmap |
| 3-4 | Structs, interfaces, embedding | Design shape hierarchy |
| 5-6 | Error handling, defer, panic/recover | Build retry library |
| 7 | Generics (type parameters) | Generic stack, queue, set |

Key questions:
- How do slices differ from arrays? What's the backing array?
- Explain interface satisfaction. Value vs pointer receivers.
- When does a variable escape to the heap?

### Week 2: Concurrency & Context

| Day | Topic | Practice |
|-----|-------|----------|
| 1-2 | Goroutines, channels, select | Fan-out/fan-in pipeline |
| 3-4 | sync package (Mutex, WaitGroup, Once) | Thread-safe cache |
| 5-6 | Context (cancellation, timeout, values) | HTTP middleware chain |
| 7 | Race conditions, deadlock detection | Fix buggy concurrent code |

Key questions:
- Explain the GMP scheduler model.
- Buffered vs unbuffered channels when to use each?
- How does context cancellation propagate?
- What happens when you send to a closed channel?

### Week 3: Patterns & Networking

| Day | Topic | Practice |
|-----|-------|----------|
| 1-2 | Design patterns (factory, strategy, observer) | Plugin system |
| 3-4 | net/http server, middleware, routing | REST API with middleware |
| 5-6 | HTTP client, connection pooling, timeouts | API client with retry |
| 7 | Testing (table-driven, mocks, benchmarks) | 90% coverage on a package |

Key questions:
- How does http.Handler interface enable middleware?
- Explain graceful shutdown of an HTTP server.
- How would you implement rate limiting?
- Table-driven tests why are they idiomatic Go?

### Week 4: System Design & DSA in Go

| Day | Topic | Practice |
|-----|-------|----------|
| 1-2 | System design with Go services | Design URL shortener |
| 3-4 | Distributed patterns (circuit breaker, saga) | Implement token bucket |
| 5-6 | DSA in Go (trees, graphs, DP) | LeetCode medium in Go |
| 7 | Mock interviews, review weak areas | Timed practice |

Key questions:
- Design a distributed cache with Go.
- How would you handle 10K concurrent WebSocket connections?
- Implement worker pool with configurable concurrency.
- Design an event-driven microservice architecture.

## Common Go Interview Questions

### Language Mechanics
1. What is the zero value of a slice, map, channel, interface?
2. How does Go's garbage collector work? (tri-color, concurrent)
3. What's the difference between `new(T)` and `&T{}`?
4. Explain method sets and why `*T` satisfies interfaces `T` doesn't.

### Concurrency
5. Implement a semaphore using channels.
6. How do you prevent goroutine leaks?
7. Explain happens-before relationships in Go.
8. sync.Pool use cases and GC interaction.

### Production
9. How do you profile a Go application? (pprof, trace)
10. Explain escape analysis and its impact on performance.
11. How do you handle graceful shutdown in Kubernetes?
12. Structured logging vs fmt.Printf why?

## Company-Specific Focus

| Company | Focus Areas |
|---------|-------------|
| **Google** | Concurrency, system design, algorithms, protocol buffers |
| **Uber** | Distributed systems, rate limiting, service mesh |
| **Cloudflare** | Networking, performance, edge computing, Workers |
| **Datadog** | Observability, metrics pipelines, high throughput |
| **Stripe** | API design, idempotency, financial correctness |
| **Docker/K8s** | Container runtime, networking, scheduler internals |

## Resources

- "The Go Programming Language" (Donovan & Kernighan)
- "Concurrency in Go" (Katherine Cox-Buday)
- Go blog: https://go.dev/blog
- Effective Go: https://go.dev/doc/effective_go
- Go by Example: https://gobyexample.com
