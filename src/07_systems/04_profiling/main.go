/*
What:       CPU and memory profiling with runtime/pprof, writing profiles, finding hot spots
Level:      Beginner
Analogy:    Profiling = a fitness tracker for your code. Shows where it spends energy.
C++ Angle:  Like perf but built into the Go toolchain. go tool pprof.
Interview:  "How do you find performance bottlenecks in Go?" → CPU profile + go tool pprof.
Run:        go run main.go && go tool pprof cpu.prof
*/
package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

// ─── CPU-intensive work: compute primes ──────────────────────────
func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func countPrimes(max int) int {
	count := 0
	for i := 2; i <= max; i++ {
		if isPrime(i) {
			count++
		}
	}
	return count
}

// ─── Memory-intensive work: allocations ──────────────────────────
func allocateStrings(n int) []string {
	result := make([]string, 0, n)
	for i := 0; i < n; i++ {
		result = append(result, fmt.Sprintf("string-%d-padding-data", i))
	}
	return result
}

func main() {
	fmt.Println("=== Go Profiling Demo ===")
	fmt.Println("  This generates cpu.prof and mem.prof for analysis.")
	fmt.Println()

	// ─── CPU Profile ─────────────────────────────────────────────
	fmt.Println("=== 1. CPU Profiling ===")
	cpuFile, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not create CPU profile: %v\n", err)
		os.Exit(1)
	}
	defer cpuFile.Close()

	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		fmt.Fprintf(os.Stderr, "could not start CPU profile: %v\n", err)
		os.Exit(1)
	}

	// Do CPU-intensive work while profiling
	start := time.Now()
	primes := countPrimes(500000)
	elapsed := time.Since(start)
	fmt.Printf("  Found %d primes below 500,000 in %v\n", primes, elapsed)

	pprof.StopCPUProfile()
	fmt.Println("  CPU profile written to: cpu.prof")

	// ─── Memory Profile ──────────────────────────────────────────
	fmt.Println("\n=== 2. Memory Profiling ===")
	// Do memory-intensive work
	strings := allocateStrings(100000)
	fmt.Printf("  Allocated %d strings\n", len(strings))

	// Force GC to get accurate heap profile
	runtime.GC()

	memFile, err := os.Create("mem.prof")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not create memory profile: %v\n", err)
		os.Exit(1)
	}
	defer memFile.Close()

	if err := pprof.WriteHeapProfile(memFile); err != nil {
		fmt.Fprintf(os.Stderr, "could not write memory profile: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("  Memory profile written to: mem.prof")

	// ─── How to Analyze ──────────────────────────────────────────
	fmt.Println("\n=== 3. Analyzing Profiles ===")
	fmt.Println("  CPU profile (find slow functions):")
	fmt.Println("    go tool pprof cpu.prof")
	fmt.Println("    (pprof) top10        ← hottest functions")
	fmt.Println("    (pprof) list isPrime  ← line-level cost")
	fmt.Println("    (pprof) web           ← open call graph in browser")
	fmt.Println()
	fmt.Println("  Memory profile (find allocations):")
	fmt.Println("    go tool pprof mem.prof")
	fmt.Println("    go tool pprof -alloc_space mem.prof  ← total allocations")
	fmt.Println("    go tool pprof -inuse_space mem.prof  ← current heap")
	fmt.Println()
	fmt.Println("  HTTP server profiling (production):")
	fmt.Println("    import _ \"net/http/pprof\"")
	fmt.Println("    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30")
	fmt.Println()
	fmt.Println("  Benchmarks with profiling:")
	fmt.Println("    go test -bench=. -cpuprofile=cpu.prof -memprofile=mem.prof")

	// ─── Cleanup note ────────────────────────────────────────────
	fmt.Println("\n=== 4. Quick Tips ===")
	fmt.Println("  • Profile in production with net/http/pprof (low overhead)")
	fmt.Println("  • Use -benchmem with go test to see allocations per op")
	fmt.Println("  • Look for: high flat% (self time) and high cum% (including callees)")
	fmt.Println("  • Common wins: reduce allocations, avoid copying large structs")

	// Clean up generated files note
	fmt.Println("\n  Generated files: cpu.prof, mem.prof")
	fmt.Println("  Delete when done: rm cpu.prof mem.prof")
}
