/*
What:       context.Background, WithCancel, WithTimeout, WithValue cancellation propagation
Level:      Beginner
Analogy:    Context = a leash for goroutines. Pull it and they all stop.
C++ Angle:  Like a cancellation token that propagates through the call graph.
Interview:  "How do you cancel a goroutine tree?" → Pass context, check ctx.Done().
*/
package main

import (
	"context"
	"fmt"
	"time"
)

// ─── Simulate a slow database query ──────────────────────────────
func slowQuery(ctx context.Context, query string) (string, error) {
	select {
	case <-time.After(200 * time.Millisecond): // simulate latency
		return fmt.Sprintf("result for '%s'", query), nil
	case <-ctx.Done():
		return "", ctx.Err() // context.DeadlineExceeded or context.Canceled
	}
}

// ─── Worker that respects cancellation ───────────────────────────
func worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("  Worker %d: stopping (%v)\n", id, ctx.Err())
			return
		case <-time.After(50 * time.Millisecond):
			fmt.Printf("  Worker %d: tick\n", id)
		}
	}
}

// ─── Nested cancellation parent cancel stops children ──────────
func parentTask(ctx context.Context) {
	childCtx, childCancel := context.WithCancel(ctx)
	defer childCancel()

	go worker(childCtx, 100)
	go worker(childCtx, 101)

	select {
	case <-ctx.Done():
		fmt.Println("  Parent: context cancelled, children will stop too.")
	case <-time.After(150 * time.Millisecond):
		fmt.Println("  Parent: finished work, cancelling children.")
	}
}

func main() {
	fmt.Println("=== 1. context.WithTimeout ===")
	// Give the query 100ms it needs 200ms, so it will timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel() // always call cancel to release resources

	result, err := slowQuery(ctx, "SELECT * FROM users")
	if err != nil {
		fmt.Printf("  Query failed: %v\n", err)
	} else {
		fmt.Printf("  Query result: %s\n", result)
	}

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 2. context.WithCancel ===")
	ctx2, cancel2 := context.WithCancel(context.Background())

	go worker(ctx2, 1)
	go worker(ctx2, 2)

	time.Sleep(130 * time.Millisecond)
	fmt.Println("  Main: cancelling all workers...")
	cancel2() // all workers see ctx.Done()
	time.Sleep(50 * time.Millisecond)

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 3. Nested Cancellation (Parent → Children) ===")
	ctx3, cancel3 := context.WithCancel(context.Background())
	go parentTask(ctx3)

	time.Sleep(100 * time.Millisecond)
	cancel3() // cancels parent, which cascades to children
	time.Sleep(100 * time.Millisecond)

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 4. context.WithValue (Request-Scoped Data) ===")
	type contextKey string
	const reqIDKey contextKey = "requestID"

	ctx4 := context.WithValue(context.Background(), reqIDKey, "req-abc-123")

	processRequest := func(ctx context.Context) {
		reqID := ctx.Value(reqIDKey).(string)
		fmt.Printf("  Processing request: %s\n", reqID)
	}
	processRequest(ctx4)

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 5. Context Best Practices ===")
	fmt.Println("  • Always pass context as first parameter: func Foo(ctx context.Context, ...)")
	fmt.Println("  • Always call cancel() use defer cancel() immediately after creation")
	fmt.Println("  • Never store context in a struct pass it through the call chain")
	fmt.Println("  • Use WithValue sparingly only for request-scoped data (trace IDs)")
	fmt.Println("  • context.Background() at top level, context.TODO() for unfinished code")
	fmt.Println("  • Check ctx.Err() after any blocking operation")
}
