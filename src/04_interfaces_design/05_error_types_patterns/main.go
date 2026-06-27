/*
What this teaches:
    Production error patterns in Go: custom error types, errors.Is/As, error
    wrapping chains, sentinel errors, opaque errors, and interface-based error
    behavior checking (e.g., IsTemporary).

Beginner analogy:
    "Errors in Go are like medical referrals each level wraps context ('patient
     had fever → caused by infection → caused by bacteria') and you can unwrap to
     find the root cause."

C++ comparison:
    "Go errors are values, not exceptions. No stack unwinding, no catch blocks.
     Think of it as returning std::expected<T, Error> everywhere. errors.Is is like
     dynamic_cast for error identity; errors.As is like catching by type."

Interview relevance:
    Interviewers ask: What's the difference between errors.Is and errors.As? How
    do you create error hierarchies without inheritance? How do you check error
    behavior (temporary/retryable) without coupling to concrete types?
*/

package main

import (
	"errors"
	"fmt"
	"net"
)

// --- Sentinel errors: package-level identity ---

var (
	ErrNotFound      = errors.New("not found")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrRateLimited   = errors.New("rate limited")
)

// --- Custom error type with context ---

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: field %q %s", e.Field, e.Message)
}

// --- Wrapping error type ---

type ServiceError struct {
	Op      string // Operation that failed
	Kind    string // Category: "database", "network", etc.
	Wrapped error  // Underlying cause
}

func (e *ServiceError) Error() string {
	if e.Wrapped != nil {
		return fmt.Sprintf("%s [%s]: %v", e.Op, e.Kind, e.Wrapped)
	}
	return fmt.Sprintf("%s [%s]", e.Op, e.Kind)
}

func (e *ServiceError) Unwrap() error {
	return e.Wrapped
}

// --- Interface-based error behavior (opaque errors) ---

type TemporaryError interface {
	IsTemporary() bool
}

type NetworkError struct {
	Host    string
	Retries int
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error connecting to %s (retries: %d)", e.Host, e.Retries)
}

func (e *NetworkError) IsTemporary() bool {
	return e.Retries < 3
}

// Check behavior, not type the opaque error pattern
func isTemporary(err error) bool {
	var te TemporaryError
	if errors.As(err, &te) {
		return te.IsTemporary()
	}
	return false
}

// --- Functions that produce wrapped errors ---

func fetchUser(id int) error {
	if id <= 0 {
		return &ServiceError{
			Op:      "fetchUser",
			Kind:    "validation",
			Wrapped: &ValidationError{Field: "id", Message: "must be positive"},
		}
	}
	if id == 404 {
		return &ServiceError{
			Op:      "fetchUser",
			Kind:    "database",
			Wrapped: ErrNotFound,
		}
	}
	return nil
}

func callExternalAPI() error {
	return &ServiceError{
		Op:   "callExternalAPI",
		Kind: "network",
		Wrapped: &NetworkError{
			Host:    "api.example.com",
			Retries: 1,
		},
	}
}

func main() {
	fmt.Println("=== Error Types & Patterns ===")

	// 1. Sentinel errors with errors.Is
	fmt.Println("\n--- Sentinel Errors (errors.Is) ---")
	err := fetchUser(404)
	if errors.Is(err, ErrNotFound) {
		fmt.Printf("  Correctly identified: %v\n", err)
	}

	// 2. Custom error types with errors.As
	fmt.Println("\n--- Custom Error Types (errors.As) ---")
	err = fetchUser(0)
	var ve *ValidationError
	if errors.As(err, &ve) {
		fmt.Printf("  Validation error field: %s, message: %s\n", ve.Field, ve.Message)
	}

	// 3. Error wrapping chain
	fmt.Println("\n--- Error Wrapping Chain ---")
	var se *ServiceError
	if errors.As(err, &se) {
		fmt.Printf("  Service error op: %s, kind: %s\n", se.Op, se.Kind)
		fmt.Printf("  Full chain: %v\n", err)
	}

	// 4. fmt.Errorf with %w for wrapping
	fmt.Println("\n--- fmt.Errorf %%w Wrapping ---")
	wrapped := fmt.Errorf("handler: %w", ErrRateLimited)
	fmt.Printf("  errors.Is(wrapped, ErrRateLimited) = %v\n", errors.Is(wrapped, ErrRateLimited))

	// 5. Interface-based behavior checking (opaque errors)
	fmt.Println("\n--- Opaque Error Behavior ---")
	apiErr := callExternalAPI()
	fmt.Printf("  Error: %v\n", apiErr)
	fmt.Printf("  Is temporary? %v\n", isTemporary(apiErr))

	// 6. Standard library pattern: net.Error
	fmt.Println("\n--- stdlib net.Error ---")
	var netErr net.Error
	fakeErr := &ServiceError{Op: "dial", Kind: "network", Wrapped: errors.New("timeout")}
	if errors.As(fakeErr, &netErr) {
		fmt.Println("  Is net.Error (won't match here demonstrating the pattern)")
	} else {
		fmt.Println("  Not a net.Error but we checked behavior, not type")
	}

	// 7. Key takeaways
	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. Sentinel errors: var ErrX = errors.New(...) for identity checks")
	fmt.Println("2. errors.Is: walks the Unwrap chain checking identity")
	fmt.Println("3. errors.As: walks the chain checking type (like type assertion)")
	fmt.Println("4. Implement Unwrap() to enable chain traversal")
	fmt.Println("5. Check behavior via interface (IsTemporary) not concrete type")
}
