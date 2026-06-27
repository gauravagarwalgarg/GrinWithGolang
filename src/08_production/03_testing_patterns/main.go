/*
Module 8: Production - Testing Patterns

Demonstrates:
  - Table-driven tests (the Go testing idiom)
  - Subtests with t.Run for granular reporting
  - Test helpers that call t.Helper()
  - Testable function design (dependency injection)
  - Golden file pattern concept
  - Benchmark structure with testing.B

Key insight: Go tests live alongside code (*_test.go), run with `go test`,
and use table-driven patterns to test many cases concisely. No assertion
library needed the stdlib testing package is sufficient.

Run: go run main.go (demonstrates the testable functions)
Test: go test -v -run TestSlugify (if this were in a _test.go file)
*/
package main

import (
	"fmt"
	"strings"
	"unicode"
)

// --- Testable functions (designed for easy testing) ---

// Slugify converts a title to a URL-friendly slug.
// Design: pure function, no side effects → trivially testable.
func Slugify(title string) string {
	var result strings.Builder
	prevDash := false

	for _, r := range strings.ToLower(title) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			result.WriteRune(r)
			prevDash = false
		case r == ' ' || r == '-' || r == '_':
			if !prevDash && result.Len() > 0 {
				result.WriteRune('-')
				prevDash = true
			}
		}
	}

	s := result.String()
	return strings.TrimRight(s, "-")
}

// --- What the test file would look like (commented for illustration) ---
/*
// FILE: main_test.go
package main

import "testing"

func TestSlugify(t *testing.T) {
    // Table-driven test: each case is independent
    tests := []struct {
        name  string  // subtest name for clear failure messages
        input string
        want  string
    }{
        {"simple", "Hello World", "hello-world"},
        {"multiple spaces", "Hello   World", "hello-world"},
        {"special chars", "Go is #1!", "go-is-1"},
        {"leading trailing", "  Hello  ", "hello"},
        {"unicode", "Café Résumé", "café-résumé"},
        {"empty", "", ""},
        {"only special", "!@#$%", ""},
        {"hyphens", "already-slugged", "already-slugged"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Slugify(tt.input)
            if got != tt.want {
                t.Errorf("Slugify(%q) = %q, want %q", tt.input, got, tt.want)
            }
        })
    }
}

// Helper function pattern: marks itself so failures report caller's line.
func assertEqual(t *testing.T, got, want string) {
    t.Helper() // marks this as helper → error points to caller
    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}

// Benchmark: measures performance with testing.B
func BenchmarkSlugify(b *testing.B) {
    input := "Hello World This Is A Test Title For Benchmarking"
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        Slugify(input)
    }
}

// Golden file pattern: compare output against saved "golden" files.
// Update golden files with: go test -update-golden
//
// func TestOutput_Golden(t *testing.T) {
//     got := generateReport()
//     golden := filepath.Join("testdata", t.Name()+".golden")
//     if *updateGolden {
//         os.WriteFile(golden, []byte(got), 0644)
//     }
//     want, _ := os.ReadFile(golden)
//     if got != string(want) {
//         t.Errorf("output mismatch, run with -update-golden")
//     }
// }
*/

func main() {
	// Demonstrate the testable function
	examples := []string{
		"Hello World",
		"Go is Awesome!",
		"  Multiple   Spaces  ",
		"special-chars_and things!",
		"",
	}

	fmt.Println("=== Slugify Examples ===")
	for _, input := range examples {
		fmt.Printf("  %-30q → %q\n", input, Slugify(input))
	}

	fmt.Println("\n=== Testing Patterns Summary ===")
	fmt.Println("1. Table-driven: []struct{name, input, want} + t.Run()")
	fmt.Println("2. Helpers: t.Helper() marks utility functions")
	fmt.Println("3. Subtests: t.Run(name, func(t *testing.T){})")
	fmt.Println("4. Parallel: t.Parallel() for concurrent subtests")
	fmt.Println("5. Golden files: compare against saved expected output")
	fmt.Println("6. Run: go test -v -run TestSlugify -count=1")
}
