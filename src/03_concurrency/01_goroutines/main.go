/*
What:       Launching goroutines, WaitGroup synchronization, closures & loop variable pitfall
Level:      Beginner
Analogy:    A goroutine is a lightweight worker bee you can have millions.
C++ Angle:  Not OS threads. 2KB initial stack, grows dynamically. No pthread_create overhead.
Interview:  "How are goroutines scheduled?" → M:N scheduler, multiplexes onto OS threads.
*/
package main

import (
	"fmt"
	"sync"
	"time"
)

func namedWorker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("  Worker %d started\n", id)
	time.Sleep(10 * time.Millisecond) // simulate work
	fmt.Printf("  Worker %d done\n", id)
}

func main() {
	fmt.Println("=== 1. Basic Goroutine Launch ===")
	// The 'go' keyword spawns a goroutine a lightweight concurrent function.
	go func() {
		fmt.Println("  Hello from anonymous goroutine!")
	}()
	time.Sleep(50 * time.Millisecond) // crude wait; use WaitGroup in real code

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 2. WaitGroup for Synchronization ===")
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go namedWorker(i, &wg)
	}
	wg.Wait() // blocks until all workers call Done()
	fmt.Println("  All workers finished.")

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 3. Loop Variable Pitfall (Pre Go 1.22) ===")
	// BUG: In Go < 1.22, the loop variable is shared across iterations.
	// All goroutines may print the LAST value of i.
	fmt.Println("  Buggy version (may print all 5s):")
	var wg2 sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg2.Add(1)
		go func() {
			// In Go 1.22+, each iteration gets its own copy of i (fixed).
			// In older Go, this captures a shared variable.
			defer wg2.Done()
			fmt.Printf("    i=%d\n", i) // Go 1.22: correct; older: may print 5
		}()
	}
	wg2.Wait()

	// FIX for pre-1.22: pass as argument
	fmt.Println("  Fixed version (parameter copy):")
	var wg3 sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg3.Add(1)
		go func(val int) {
			defer wg3.Done()
			fmt.Printf("    val=%d\n", val)
		}(i) // copy i into val
	}
	wg3.Wait()

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 4. 1000 Goroutines They're Cheap! ===")
	// Each goroutine starts with ~2KB stack. 1000 goroutines ≈ 2MB.
	// An OS thread typically uses 1-8MB per thread.
	const numGoroutines = 1000
	var wg4 sync.WaitGroup
	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg4.Add(1)
		go func(id int) {
			defer wg4.Done()
			// Simulate trivial work
			_ = id * id
		}(i)
	}
	wg4.Wait()
	elapsed := time.Since(start)

	fmt.Printf("  Launched and joined %d goroutines in %v\n", numGoroutines, elapsed)
	fmt.Println("  Compare: pthread_create for 1000 threads would take 10-100x longer")

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 5. Goroutine vs Thread Summary ===")
	fmt.Println("  ┌──────────────┬───────────────┬───────────────────┐")
	fmt.Println("  │ Property     │ Goroutine     │ OS Thread         │")
	fmt.Println("  ├──────────────┼───────────────┼───────────────────┤")
	fmt.Println("  │ Stack size   │ 2KB (grows)   │ 1-8MB (fixed)     │")
	fmt.Println("  │ Create cost  │ ~300ns        │ ~30μs             │")
	fmt.Println("  │ Scheduling   │ Go runtime    │ OS kernel         │")
	fmt.Println("  │ Typical max  │ 100K-1M       │ ~10K              │")
	fmt.Println("  └──────────────┴───────────────┴───────────────────┘")
}
