/*
What:       Stack vs heap allocation, escape analysis, what causes escapes
Level:      Beginner
Analogy:    Stack = your desk (fast, cleaned automatically). Heap = a storage locker (slow, GC cleans it).
C++ Angle:  Go decides stack/heap at compile time via escape analysis no manual new/malloc.
Interview:  "How do you check if a variable escapes?" → go build -gcflags="-m"
Run:        go build -gcflags="-m" main.go  ← shows escape decisions
*/
package main

import (
	"fmt"
)

// ─── Case 1: Does NOT escape stays on stack ────────────────────
func sumOnStack(a, b int) int {
	result := a + b // result lives and dies with this function call
	return result   // value is copied to caller no pointer needed
}

// ─── Case 2: ESCAPES returned pointer forces heap allocation ───
func createOnHeap(name string) *string {
	s := "Hello, " + name // s must outlive this function
	return &s             // &s escapes to heap because caller holds the pointer
}

// ─── Case 3: ESCAPES interface assignment ──────────────────────
func interfaceEscape() {
	x := 42
	// fmt.Println takes interface{} compiler can't prove x won't escape
	fmt.Println("  Interface escape:", x) // x escapes to heap
}

// ─── Case 4: ESCAPES closure captures local variable ───────────
func closureEscape() func() int {
	counter := 0
	return func() int {
		counter++ // counter must live beyond closureEscape's stack frame
		return counter
	}
}

// ─── Case 5: Does NOT escape pointer used only within function ─
func pointerNoEscape() {
	x := 100
	p := &x // p doesn't escape used only here
	*p += 1
	_ = *p
}

// ─── Case 6: ESCAPES slice too large for stack ─────────────────
func largeSliceEscape() []byte {
	// Slices larger than a threshold are heap-allocated
	buf := make([]byte, 64*1024) // 64KB too big for stack
	buf[0] = 1
	return buf
}

// ─── Case 7: Does NOT escape small, non-returned slice ─────────
func smallSliceNoEscape() {
	buf := make([]byte, 128) // small, doesn't leave the function
	buf[0] = 'A'
	_ = buf
}

func main() {
	fmt.Println("=== Escape Analysis Demo ===")
	fmt.Println("  Run: go build -gcflags=\"-m\" main.go")
	fmt.Println()

	// Case 1: Stack allocation (no escape)
	fmt.Println("=== 1. Value return stays on stack ===")
	result := sumOnStack(3, 4)
	fmt.Printf("  sum: %d (allocated on stack)\n", result)

	// Case 2: Pointer return escapes to heap
	fmt.Println("\n=== 2. Pointer return escapes to heap ===")
	greeting := createOnHeap("World")
	fmt.Printf("  greeting: %s (heap allocated)\n", *greeting)

	// Case 3: Interface assignment escapes
	fmt.Println("\n=== 3. Interface assignment escapes ===")
	interfaceEscape()

	// Case 4: Closure capture escapes
	fmt.Println("\n=== 4. Closure capture variable escapes ===")
	inc := closureEscape()
	fmt.Printf("  counter: %d, %d, %d\n", inc(), inc(), inc())

	// Case 5: Pointer stays local no escape
	fmt.Println("\n=== 5. Local pointer no escape ===")
	pointerNoEscape()
	fmt.Println("  (pointer used only within function, stays on stack)")

	// Case 6: Large allocation escapes
	fmt.Println("\n=== 6. Large slice escapes ===")
	big := largeSliceEscape()
	fmt.Printf("  big slice len=%d (heap allocated)\n", len(big))

	// Case 7: Small slice no escape
	fmt.Println("\n=== 7. Small slice no escape ===")
	smallSliceNoEscape()
	fmt.Println("  (small slice stays on stack)")

	// ─────────────────────────────────────────────
	fmt.Println("\n=== Escape Analysis Cheat Sheet ===")
	fmt.Println("  ESCAPES (heap):")
	fmt.Println("    • Returning a pointer to a local variable")
	fmt.Println("    • Assigning to an interface (fmt.Println, error, etc.)")
	fmt.Println("    • Closure capturing a mutable local")
	fmt.Println("    • Slice/map too large for stack")
	fmt.Println("    • Sending pointer to a channel")
	fmt.Println("  STAYS ON STACK:")
	fmt.Println("    • Value types returned by copy")
	fmt.Println("    • Pointers that don't leave the function")
	fmt.Println("    • Small slices not returned or shared")
	fmt.Println("  CHECK: go build -gcflags=\"-m=2\" for detailed decisions")
}
