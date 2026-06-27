/*
What:       GC behavior, runtime.GC(), ReadMemStats, GOGC/GOMEMLIMIT, allocation pressure
Level:      Beginner
Analogy:    GC = janitor who cleans when rooms get messy. More trash = more cleaning.
C++ Angle:  Concurrent tri-color mark-and-sweep. Sub-ms pauses. No RAII needed but understand allocation pressure.
Interview:  "How do you tune GC in Go?" → GOGC (growth ratio), GOMEMLIMIT (hard cap), reduce allocations.
Run:        GODEBUG=gctrace=1 go run main.go  ← shows GC events
*/
package main

import (
	"fmt"
	"runtime"
	"time"
)

func printMemStats(label string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("  [%s]\n", label)
	fmt.Printf("    Alloc      = %6d KB (currently allocated heap)\n", m.Alloc/1024)
	fmt.Printf("    TotalAlloc = %6d KB (cumulative allocated)\n", m.TotalAlloc/1024)
	fmt.Printf("    Sys        = %6d KB (OS memory obtained)\n", m.Sys/1024)
	fmt.Printf("    NumGC      = %6d    (GC cycles completed)\n", m.NumGC)
	fmt.Printf("    PauseTotal = %6v    (total GC pause time)\n", time.Duration(m.PauseTotalNs))
	fmt.Println()
}

// allocateGarbage creates short-lived allocations that pressure the GC.
func allocateGarbage(n int) {
	for i := 0; i < n; i++ {
		// Each iteration allocates a slice that becomes garbage immediately
		_ = make([]byte, 1024) // 1KB per allocation
	}
}

// allocateRetained keeps references alive grows heap.
func allocateRetained(n int) [][]byte {
	retained := make([][]byte, n)
	for i := range retained {
		retained[i] = make([]byte, 1024)
	}
	return retained
}

func main() {
	fmt.Println("=== 1. Baseline Memory Stats ===")
	printMemStats("startup")

	// ─────────────────────────────────────────────
	fmt.Println("=== 2. Allocation Pressure (Short-Lived Objects) ===")
	fmt.Println("  Creating 10,000 short-lived 1KB allocations...")
	allocateGarbage(10000)
	printMemStats("after garbage")

	// ─────────────────────────────────────────────
	fmt.Println("=== 3. Forcing GC ===")
	runtime.GC() // explicit GC rarely needed in production
	printMemStats("after runtime.GC()")

	// ─────────────────────────────────────────────
	fmt.Println("=== 4. Retained Allocations (Heap Growth) ===")
	fmt.Println("  Allocating 50,000 retained 1KB slices (~50MB)...")
	retained := allocateRetained(50000)
	printMemStats("after retained alloc")

	// ─────────────────────────────────────────────
	fmt.Println("=== 5. Releasing References ===")
	_ = retained      // keep compiler happy
	retained = nil    // release reference GC can now reclaim
	runtime.GC()      // trigger collection
	printMemStats("after release + GC")

	// ─────────────────────────────────────────────
	fmt.Println("=== 6. GC Tuning Environment Variables ===")
	fmt.Println("  GOGC=100 (default)")
	fmt.Println("    GC triggers when heap grows 100% since last GC.")
	fmt.Println("    GOGC=50  → GC more often, less memory used")
	fmt.Println("    GOGC=200 → GC less often, more memory, less CPU")
	fmt.Println("    GOGC=off → disable GC (for benchmarks only!)")
	fmt.Println()
	fmt.Println("  GOMEMLIMIT=1GiB (Go 1.19+)")
	fmt.Println("    Hard memory limit. GC works harder to stay under.")
	fmt.Println("    Better than GOGC for memory-constrained environments.")
	fmt.Println()
	fmt.Println("  GODEBUG=gctrace=1")
	fmt.Println("    Prints GC events to stderr:")
	fmt.Println("    gc 1 @0.012s 2%: 0.1+1.2+0.1 ms clock, ...")

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 7. Reducing Allocation Pressure ===")
	fmt.Println("  • Reuse buffers with sync.Pool")
	fmt.Println("  • Pre-allocate slices: make([]T, 0, expectedCap)")
	fmt.Println("  • Avoid fmt.Sprintf in hot paths (allocates strings)")
	fmt.Println("  • Use value types over pointers when possible")
	fmt.Println("  • Profile with: go tool pprof -alloc_space")
}
