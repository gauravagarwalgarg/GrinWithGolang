/*
What this teaches:
    Common Go interview pitfalls and patterns: goroutine leak detection, slice
    append gotchas, nil interface != nil, map concurrent write panic, and range
    pointer aliasing. Each shows the bug then the fix.

Beginner analogy:
    "These are the landmines every Go developer steps on once. Learn them here so
     you don't learn them in production at 2 AM. Each example shows the broken code,
     explains WHY it breaks, and shows the correct version."

C++ comparison:
    "C++ has its own pitfalls (dangling references, UB, iterator invalidation). Go's
     pitfalls are different but equally subtle: interface nil wrapping, goroutine
     leaks (no RAII destructor), and channel misuse."

Interview relevance:
    These are the exact gotchas interviewers use to separate junior from senior Go
    developers. Understanding WHY each bug occurs demonstrates deep language knowledge.
*/

package main

import (
	"fmt"
	"sync"
	"time"
)

// nilDemoError is package-level so it can have methods (needed for pitfall 3)
type nilDemoError struct{ Msg string }

func (e *nilDemoError) Error() string { return e.Msg }

func main() {
	fmt.Println("=== Go Interview Pitfalls & Patterns ===")

	pitfall1_goroutineLeak()
	pitfall2_sliceAppend()
	pitfall3_nilInterface()
	pitfall4_mapConcurrentWrite()
	pitfall5_rangePointerAlias()
}

// --- Pitfall 1: Goroutine Leak ---

func pitfall1_goroutineLeak() {
	fmt.Println("\n--- Pitfall 1: Goroutine Leak ---")

	// BUG: goroutine blocked forever on unbuffered channel nobody reads
	fmt.Println("  BUG: Channel never read → goroutine leaks forever")
	fmt.Println(`  ch := make(chan int)
  go func() { ch <- expensiveCompute() }() // LEAKED if we never read ch`)

	// FIX 1: Always ensure channels are consumed or use context cancellation
	fmt.Println("\n  FIX: Use buffered channel or context cancellation")
	done := make(chan int, 1) // Buffered: writer won't block even if reader is gone
	go func() {
		done <- 42
	}()
	val := <-done
	fmt.Printf("  Result: %d (no leak channel was read)\n", val)

	// FIX 2: Context-based cancellation for long-running goroutines
	fmt.Println("  FIX 2: Use context.WithCancel for long-running goroutines")
}

// --- Pitfall 2: Slice Append Gotcha ---

func pitfall2_sliceAppend() {
	fmt.Println("\n--- Pitfall 2: Slice Append Gotcha ---")

	// BUG: Slicing shares underlying array; append may or may not reallocate
	fmt.Println("  BUG: Slices share backing array")
	original := make([]int, 3, 5) // len=3, cap=5
	original[0], original[1], original[2] = 1, 2, 3

	sliceA := original[:2]       // Shares backing array, cap=5
	sliceA = append(sliceA, 99)  // Writes into original[2]!

	fmt.Printf("  original: %v (CORRUPTED! original[2] is now 99)\n", original)

	// FIX: Use full slice expression to limit capacity
	fmt.Println("\n  FIX: Full slice expression limits capacity")
	original2 := []int{1, 2, 3, 4, 5}
	safe := original2[:2:2] // len=2, cap=2 forces new allocation on append
	safe = append(safe, 99)
	fmt.Printf("  original2: %v (unchanged)\n", original2)
	fmt.Printf("  safe:      %v (independent copy)\n", safe)
}

// --- Pitfall 3: nil Interface != nil ---

func pitfall3_nilInterface() {
	fmt.Println("\n--- Pitfall 3: nil Interface != nil ---")

	// BUG: An interface holding a nil pointer is NOT nil
	fmt.Println("  BUG: Interface with nil concrete value is NOT nil")

	// Using a package-level error type (see nilDemoError above)
	var ptr *nilDemoError = nil    // nil pointer
	var iface error = ptr          // interface wraps (type=*nilDemoError, value=nil)
	fmt.Printf("  ptr == nil:   %v\n", ptr == nil)    // true
	fmt.Printf("  iface == nil: %v\n", iface == nil)  // FALSE!
	fmt.Printf("  Why: interface has type info (*nilDemoError) even though value is nil\n")

	// FIX: Return the interface type directly, don't assign typed nil
	fmt.Println("\n  FIX: Return error (interface) directly")
	fmt.Println(`  func getError() error {
    var e *MyError = nil
    if e == nil { return nil }  // Return untyped nil
    return e
  }`)
}

// --- Pitfall 4: Map Concurrent Write Panic ---

func pitfall4_mapConcurrentWrite() {
	fmt.Println("\n--- Pitfall 4: Map Concurrent Write Panic ---")

	fmt.Println("  BUG: Concurrent map writes cause runtime panic")
	fmt.Println(`  m := map[string]int{}
  go func() { m["a"] = 1 }()
  go func() { m["b"] = 2 }()  // PANIC: concurrent map writes`)

	// FIX: Use sync.Mutex or sync.Map
	fmt.Println("\n  FIX: Protect with sync.Mutex")
	var mu sync.Mutex
	m := make(map[string]int)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			mu.Lock()
			m[fmt.Sprintf("key%d", n)] = n
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	fmt.Printf("  Safe map has %d entries (no panic)\n", len(m))
}

// --- Pitfall 5: Range Pointer Aliasing ---

func pitfall5_rangePointerAlias() {
	fmt.Println("\n--- Pitfall 5: Range Variable Pointer Aliasing ---")

	// BUG: Taking address of range variable all pointers point to same memory
	type User struct{ Name string }
	users := []User{{"Alice"}, {"Bob"}, {"Charlie"}}

	// BUG version
	fmt.Println("  BUG: All pointers alias the loop variable")
	var bugPtrs []*User
	for _, u := range users {
		bugPtrs = append(bugPtrs, &u) // All point to same 'u'!
	}
	fmt.Printf("  bugPtrs: ")
	for _, p := range bugPtrs {
		fmt.Printf("%s ", p.Name) // All print "Charlie"!
	}
	fmt.Println(" (all 'Charlie'!)")

	// FIX 1: Use index
	fmt.Println("\n  FIX 1: Use index to get stable pointer")
	var fixPtrs []*User
	for i := range users {
		fixPtrs = append(fixPtrs, &users[i])
	}
	fmt.Printf("  fixPtrs: ")
	for _, p := range fixPtrs {
		fmt.Printf("%s ", p.Name)
	}
	fmt.Println()

	// FIX 2: Local copy (pre-Go 1.22)
	fmt.Println("\n  FIX 2: Shadow the variable")
	var fixPtrs2 []*User
	for _, u := range users {
		u := u // Shadow creates new variable each iteration
		fixPtrs2 = append(fixPtrs2, &u)
	}
	fmt.Printf("  fixPtrs2: ")
	for _, p := range fixPtrs2 {
		fmt.Printf("%s ", p.Name)
	}
	fmt.Println()

	// Note about Go 1.22+
	fmt.Println("\n  NOTE: Go 1.22+ changes loop variable semantics (per-iteration scope)")

	// Brief pause to let any goroutines finish
	time.Sleep(10 * time.Millisecond)

	fmt.Println("\n--- Summary of Pitfalls ---")
	fmt.Println("1. Goroutine leak: always ensure channels are drained or use context")
	fmt.Println("2. Slice append: use s[:n:n] to prevent shared-array corruption")
	fmt.Println("3. Nil interface: interface{type, nil} != nil; return untyped nil")
	fmt.Println("4. Map panic: maps are not goroutine-safe; use mutex or sync.Map")
	fmt.Println("5. Range alias: &v captures loop variable; use index or shadow")
}
