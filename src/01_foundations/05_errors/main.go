/*
What this teaches:
    Error handling in Go the error interface, errors.New, fmt.Errorf with %w
    for wrapping, errors.Is/As for inspection, and sentinel error patterns.

Beginner analogy:
    "Errors are just values like a report card you check after every exam.
     You don't throw it across the room (exceptions); you look at it calmly
     and decide what to do."

C++ comparison:
    "No exceptions. Errors are explicit return values. No hidden control flow,
     no stack unwinding costs. Every function that can fail returns (T, error).
     This makes error paths visible and greppable."

Interview relevance:
    You'll be asked: Why no exceptions? How to wrap/unwrap errors? Difference
    between errors.Is and ==? When to use sentinel vs custom error types?
*/

package main

import (
	"errors"
	"fmt"
	"strconv"
)

// Sentinel errors package-level, reusable error values
var (
	ErrNotFound   = errors.New("not found")
	ErrOutOfRange = errors.New("index out of range")
)

// Custom error type implements the error interface
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on %q: %s", e.Field, e.Message)
}

// Functions that return errors
func findUser(id int) (string, error) {
	users := map[int]string{1: "Alice", 2: "Bob"}
	name, ok := users[id]
	if !ok {
		return "", fmt.Errorf("findUser(%d): %w", id, ErrNotFound)
	}
	return name, nil
}

func validateAge(age string) (int, error) {
	n, err := strconv.Atoi(age)
	if err != nil {
		return 0, fmt.Errorf("validateAge: %w", err)
	}
	if n < 0 || n > 150 {
		return 0, &ValidationError{Field: "age", Message: "must be 0-150"}
	}
	return n, nil
}

func main() {
	// 1. Basic error handling pattern
	fmt.Println("--- Basic Error Handling ---")
	name, err := findUser(1)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Found:", name)
	}

	// 2. Handling a missing user
	fmt.Println("\n--- Sentinel Error ---")
	_, err = findUser(99)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// 3. errors.Is unwraps the chain to match sentinel
	fmt.Println("\n--- errors.Is (unwrap matching) ---")
	if errors.Is(err, ErrNotFound) {
		fmt.Println("✓ Confirmed: the root cause is ErrNotFound")
	}

	// 4. fmt.Errorf with %w wrapping adds context
	fmt.Println("\n--- Error Wrapping with %w ---")
	_, err = validateAge("abc")
	fmt.Println("Wrapped error:", err)

	// Unwrap to inspect the original strconv error
	var numErr *strconv.NumError
	if errors.As(err, &numErr) {
		fmt.Printf("✓ errors.As found NumError: Func=%s Num=%s\n",
			numErr.Func, numErr.Num)
	}

	// 5. errors.As match custom error types
	fmt.Println("\n--- errors.As (custom type) ---")
	_, err = validateAge("-5")
	var valErr *ValidationError
	if errors.As(err, &valErr) {
		fmt.Printf("✓ ValidationError: field=%q msg=%s\n",
			valErr.Field, valErr.Message)
	}

	// 6. Multiple error checks in sequence (idiomatic Go)
	fmt.Println("\n--- Idiomatic Sequential Checks ---")
	ids := []int{1, 2, 3, 99}
	for _, id := range ids {
		user, err := findUser(id)
		if err != nil {
			fmt.Printf("  id=%d → ERROR: %v\n", id, err)
			continue
		}
		fmt.Printf("  id=%d → %s\n", id, user)
	}

	// 7. Creating simple errors
	fmt.Println("\n--- Creating Errors ---")
	e1 := errors.New("something went wrong")  // static message
	e2 := fmt.Errorf("failed at step %d", 3) // formatted message
	fmt.Println("errors.New:", e1)
	fmt.Println("fmt.Errorf:", e2)

	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. Errors are values check them explicitly with if err != nil")
	fmt.Println("2. Wrap with %w to add context while preserving the original")
	fmt.Println("3. errors.Is checks the chain for a specific value (sentinel)")
	fmt.Println("4. errors.As checks the chain for a specific type")
	fmt.Println("5. No hidden control flow error paths are always visible")
}
