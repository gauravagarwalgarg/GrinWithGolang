/*
Module 8: Production - Benchmarking Patterns

Demonstrates:
  - How testing.B benchmarks work (structure shown in comments)
  - Manual timing with time.Now for demonstration
  - Comparing implementations: string concat vs strings.Builder
  - Memory allocation awareness
  - Sub-benchmarks and benchmark flags

Key insight: String concatenation in a loop creates O(n²) allocations
because strings are immutable. strings.Builder amortizes to O(n).
Benchmarks prove this empirically. Always benchmark before optimizing.

Run: go run main.go
Benchmark: go test -bench=. -benchmem (if in _test.go)
*/
package main

import (
	"fmt"
	"strings"
	"time"
)

// --- Implementations to compare ---

// concatStrings builds a string using + operator (naive).
func concatStrings(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += "x"
	}
	return s
}

// builderStrings builds a string using strings.Builder (optimized).
func builderStrings(n int) string {
	var b strings.Builder
	b.Grow(n) // pre-allocate: single allocation
	for i := 0; i < n; i++ {
		b.WriteByte('x')
	}
	return b.String()
}

// --- Benchmark harness (manual, since we need main) ---

func benchmark(name string, iterations int, fn func()) time.Duration {
	start := time.Now()
	for i := 0; i < iterations; i++ {
		fn()
	}
	elapsed := time.Since(start)
	perOp := elapsed / time.Duration(iterations)
	fmt.Printf("  %-25s %8d iterations  %12v/op\n", name, iterations, perOp)
	return elapsed
}

/* --- What real benchmarks look like (in a _test.go file) ---

func BenchmarkConcatStrings(b *testing.B) {
    for i := 0; i < b.N; i++ {
        concatStrings(1000)
    }
}

func BenchmarkBuilderStrings(b *testing.B) {
    for i := 0; i < b.N; i++ {
        builderStrings(1000)
    }
}

// Sub-benchmarks for different sizes:
func BenchmarkStringBuild(b *testing.B) {
    sizes := []int{10, 100, 1000, 10000}
    for _, size := range sizes {
        b.Run(fmt.Sprintf("concat/%d", size), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                concatStrings(size)
            }
        })
        b.Run(fmt.Sprintf("builder/%d", size), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                builderStrings(size)
            }
        })
    }
}

// Run: go test -bench=BenchmarkStringBuild -benchmem
// Flags: -count=5 (run 5 times for stability)
//        -benchtime=3s (run each for 3 seconds)
//        -cpuprofile=cpu.out (generate CPU profile)
*/

func main() {
	fmt.Println("=== String Building Benchmark ===")
	fmt.Println()

	sizes := []int{100, 1000, 10000}
	for _, size := range sizes {
		fmt.Printf("Size: %d characters\n", size)
		benchmark("concat (+)", 1000, func() { concatStrings(size) })
		benchmark("strings.Builder", 1000, func() { builderStrings(size) })
		fmt.Println()
	}

	fmt.Println("=== Key Takeaways ===")
	fmt.Println("• strings.Builder is O(n), concat is O(n²)")
	fmt.Println("• Builder.Grow(n) pre-allocates → single alloc")
	fmt.Println("• Always use -benchmem to see allocations")
	fmt.Println("• Real benchmarks: go test -bench=. -benchmem -count=5")
}
