/*
What this teaches:
    Circuit breaker pattern: Closed → Open → HalfOpen states. Tracks failures and
    successes, trips at a threshold, and enters cooldown before retrying. Protects
    against cascading failures in microservices.

Beginner analogy:
    "Like an electrical circuit breaker: too many failures (overcurrent) trips the
     breaker open. After a cooldown period, it goes half-open (testing one request).
     If that succeeds, it resets to closed. If it fails, it opens again."

C++ comparison:
    "Same pattern as Netflix Hystrix or resilience4j, but implemented with Go's
     sync.Mutex and time.Now() instead of atomic CAS loops. Channel-based
     alternatives exist but mutex is simpler for state machines."

Interview relevance:
    Microservice interviews frequently ask about resilience patterns. Demonstrating
    a working circuit breaker with configurable thresholds shows systems thinking
    beyond happy-path coding.
*/

package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// --- Circuit Breaker States ---

type State int

const (
	Closed   State = iota // Normal operation
	Open                  // Rejecting requests
	HalfOpen              // Testing one request
)

func (s State) String() string {
	switch s {
	case Closed:
		return "CLOSED"
	case Open:
		return "OPEN"
	case HalfOpen:
		return "HALF-OPEN"
	}
	return "UNKNOWN"
}

// --- Circuit Breaker ---

type CircuitBreaker struct {
	mu             sync.Mutex
	state          State
	failures       int
	successes      int
	threshold      int           // Failures to trip
	resetTimeout   time.Duration // Cooldown before half-open
	halfOpenMax    int           // Successes to reset from half-open
	lastFailure    time.Time
}

func NewCircuitBreaker(threshold int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:        Closed,
		threshold:    threshold,
		resetTimeout: resetTimeout,
		halfOpenMax:  2, // 2 successes in half-open → close
	}
}

var ErrCircuitOpen = errors.New("circuit breaker is open")

func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mu.Lock()

	// Check if we should transition from Open → HalfOpen
	if cb.state == Open {
		if time.Since(cb.lastFailure) > cb.resetTimeout {
			cb.state = HalfOpen
			cb.successes = 0
			fmt.Printf("  [CB] Transitioning: OPEN → HALF-OPEN\n")
		} else {
			cb.mu.Unlock()
			return ErrCircuitOpen
		}
	}

	state := cb.state
	cb.mu.Unlock()

	// Execute the function
	err := fn()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failures++
		cb.lastFailure = time.Now()
		if state == HalfOpen || cb.failures >= cb.threshold {
			cb.state = Open
			fmt.Printf("  [CB] Transitioning: %s → OPEN (failures: %d)\n", state, cb.failures)
		}
		return err
	}

	// Success path
	if state == HalfOpen {
		cb.successes++
		if cb.successes >= cb.halfOpenMax {
			cb.state = Closed
			cb.failures = 0
			fmt.Printf("  [CB] Transitioning: HALF-OPEN → CLOSED (recovered!)\n")
		}
	} else {
		cb.failures = 0 // Reset consecutive failure count on success
	}
	return nil
}

func (cb *CircuitBreaker) State() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// --- Simulated external service ---

type ExternalService struct {
	callCount int
	failUntil int // Fail for the first N calls
}

func (s *ExternalService) Call() error {
	s.callCount++
	if s.callCount <= s.failUntil {
		return fmt.Errorf("service unavailable (call %d)", s.callCount)
	}
	return nil
}

func main() {
	fmt.Println("=== Circuit Breaker Pattern ===")

	cb := NewCircuitBreaker(3, 500*time.Millisecond)
	svc := &ExternalService{failUntil: 5} // First 5 calls fail

	// Phase 1: Accumulate failures → trip open
	fmt.Println("\n--- Phase 1: Building up failures ---")
	for i := 1; i <= 5; i++ {
		err := cb.Execute(svc.Call)
		fmt.Printf("  Call %d: state=%s err=%v\n", i, cb.State(), err)
	}

	// Phase 2: Requests rejected while open
	fmt.Println("\n--- Phase 2: Circuit OPEN fast-fail ---")
	for i := 6; i <= 8; i++ {
		err := cb.Execute(svc.Call)
		fmt.Printf("  Call %d: state=%s err=%v\n", i, cb.State(), err)
	}

	// Phase 3: Wait for cooldown → half-open → recover
	fmt.Println("\n--- Phase 3: Cooldown + Recovery ---")
	time.Sleep(600 * time.Millisecond)
	for i := 9; i <= 12; i++ {
		err := cb.Execute(svc.Call)
		fmt.Printf("  Call %d: state=%s err=%v\n", i, cb.State(), err)
	}

	// Key takeaways
	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. CLOSED: normal operation, tracking consecutive failures")
	fmt.Println("2. OPEN: reject immediately (fail-fast) no load on failing service")
	fmt.Println("3. HALF-OPEN: test with limited requests after cooldown")
	fmt.Println("4. Prevents cascading failures across microservice boundaries")
	fmt.Println("5. Configurable: threshold, cooldown, half-open success count")
}
