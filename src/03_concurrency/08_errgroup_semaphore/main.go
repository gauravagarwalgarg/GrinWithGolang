/*
What:       Semaphore pattern (bounded concurrency) and first-error collection no external deps
Level:      Beginner
Analogy:    Semaphore = a bouncer at a club. Only N people inside at once.
C++ Angle:  Like std::counting_semaphore (C++20) implemented via a buffered channel.
Interview:  "Implement errgroup without golang.org/x/sync" → channel semaphore + sync.Once for error.
*/
package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ─── Channel-based semaphore ─────────────────────────────────────
type Semaphore chan struct{}

func NewSemaphore(n int) Semaphore { return make(Semaphore, n) }
func (s Semaphore) Acquire()       { s <- struct{}{} }
func (s Semaphore) Release()       { <-s }

// ─── ErrGroup-like: run N tasks, bounded concurrency, first error ─
func runBounded(tasks []func() error, maxConcurrency int) error {
	sem := NewSemaphore(maxConcurrency)
	var wg sync.WaitGroup
	var once sync.Once
	var firstErr error

	for _, task := range tasks {
		wg.Add(1)
		sem.Acquire() // blocks if maxConcurrency goroutines are running
		go func(t func() error) {
			defer wg.Done()
			defer sem.Release()
			if err := t(); err != nil {
				once.Do(func() { firstErr = err }) // capture first error only
			}
		}(task)
	}

	wg.Wait()
	return firstErr
}

func main() {
	fmt.Println("=== 1. Semaphore Bounded Concurrency ===")
	sem := NewSemaphore(3) // allow max 3 concurrent operations
	var wg sync.WaitGroup

	for i := 1; i <= 8; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sem.Acquire()
			defer sem.Release()
			fmt.Printf("  Task %d: running (max 3 concurrent)\n", id)
			time.Sleep(time.Duration(rand.Intn(50)+20) * time.Millisecond)
			fmt.Printf("  Task %d: done\n", id)
		}(i)
	}
	wg.Wait()

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 2. ErrGroup Pattern Collect First Error ===")
	tasks := make([]func() error, 10)
	for i := range tasks {
		id := i
		tasks[i] = func() error {
			time.Sleep(time.Duration(rand.Intn(30)+10) * time.Millisecond)
			if id == 4 || id == 7 { // simulate failures
				return fmt.Errorf("task %d failed", id)
			}
			fmt.Printf("  Task %d: success\n", id)
			return nil
		}
	}

	err := runBounded(tasks, 3)
	if err != nil {
		fmt.Printf("  First error: %v\n", err)
	} else {
		fmt.Println("  All tasks succeeded")
	}

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 3. Collecting ALL Errors ===")
	var mu sync.Mutex
	var allErrors []error
	var wg2 sync.WaitGroup
	sem2 := NewSemaphore(4)

	for i := 0; i < 6; i++ {
		wg2.Add(1)
		sem2.Acquire()
		go func(id int) {
			defer wg2.Done()
			defer sem2.Release()
			time.Sleep(10 * time.Millisecond)
			if id%2 == 0 {
				mu.Lock()
				allErrors = append(allErrors, fmt.Errorf("task %d failed", id))
				mu.Unlock()
			}
		}(i)
	}
	wg2.Wait()

	combined := errors.Join(allErrors...)
	if combined != nil {
		fmt.Printf("  All errors: %v\n", combined)
	}

	fmt.Println("\n=== Key Pattern ===")
	fmt.Println("  Buffered channel as semaphore: make(chan struct{}, N)")
	fmt.Println("  Acquire = send (blocks when full), Release = receive")
}
