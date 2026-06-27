/*
What:       sync.Mutex, sync.RWMutex, sync.Once, sync/atomic when mutex vs channel
Level:      Beginner
Analogy:    Mutex = bathroom lock. Only one person at a time.
C++ Angle:  Same semantics as std::mutex but with a race detector built into the toolchain.
Interview:  "When channels vs mutex?" → Channels for communication, mutex for state protection.
Run:        go run -race main.go  ← detects data races at runtime
*/
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ─── Unsafe counter (DATA RACE with -race flag) ──────────────────
type UnsafeCounter struct {
	count int
}

// ─── Mutex-protected counter ─────────────────────────────────────
type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// ─── RWMutex: many readers, one writer ───────────────────────────
type Config struct {
	mu      sync.RWMutex
	setting string
}

func (c *Config) Get() string {
	c.mu.RLock() // multiple readers allowed
	defer c.mu.RUnlock()
	return c.setting
}

func (c *Config) Set(val string) {
	c.mu.Lock() // exclusive writer
	defer c.mu.Unlock()
	c.setting = val
}

func main() {
	fmt.Println("=== 1. Data Race (run with: go run -race main.go) ===")
	unsafe := &UnsafeCounter{}
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			unsafe.count++ // DATA RACE! No synchronization.
		}()
	}
	wg.Wait()
	fmt.Printf("  Unsafe counter (likely wrong): %d\n", unsafe.count)

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 2. sync.Mutex Correct Counter ===")
	safe := &SafeCounter{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			safe.Increment()
		}()
	}
	wg.Wait()
	fmt.Printf("  Safe counter: %d (always 1000)\n", safe.Value())

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 3. sync.RWMutex Many Readers, One Writer ===")
	cfg := &Config{setting: "initial"}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("  Reader %d: %s\n", id, cfg.Get())
		}(i)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		cfg.Set("updated")
		fmt.Println("  Writer: set to 'updated'")
	}()
	wg.Wait()

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 4. sync.Once Exactly Once Initialization ===")
	var once sync.Once
	initDB := func() {
		fmt.Println("  Connecting to database... (happens only once)")
		time.Sleep(50 * time.Millisecond)
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			once.Do(initDB) // only first goroutine runs this
			fmt.Printf("  Goroutine %d: using DB connection\n", id)
		}(i)
	}
	wg.Wait()

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 5. sync/atomic Lock-Free Counter ===")
	var atomicCount int64
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&atomicCount, 1)
		}()
	}
	wg.Wait()
	fmt.Printf("  Atomic counter: %d\n", atomic.LoadInt64(&atomicCount))

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 6. When to Use What ===")
	fmt.Println("  • Mutex: protecting shared state (maps, counters, configs)")
	fmt.Println("  • RWMutex: read-heavy workloads (95% reads, 5% writes)")
	fmt.Println("  • Atomic: simple counters, flags, CAS operations")
	fmt.Println("  • Channels: communication between goroutines, signaling")
	fmt.Println("  Rule: 'Share memory by communicating, not communicate by sharing memory'")
}
