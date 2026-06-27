/*
What this teaches:
    Pointers in Go & (address-of) and * (dereference), passing by value vs
    by pointer, nil pointers, and why Go has NO pointer arithmetic.

Beginner analogy:
    "A pointer is the address on an envelope, not the house itself. You can
     share the address so multiple people find the same house, but you can't
     wander to the neighbor's house by adding 1 to the address (no arithmetic)."

C++ comparison:
    "No pointer arithmetic, no void*, no manual free. GC handles deallocation.
     Go pointers are safe they point to valid memory or nil. You still need
     to check for nil, but you'll never have dangling pointers or use-after-free."

Interview relevance:
    When to pass by pointer vs value? What escapes to the heap? Method receivers
    (pointer vs value) covered in structs but rooted in this concept.
*/

package main

import "fmt"

// Pass by value the function gets a COPY
func doubleValue(n int) {
	n *= 2 // modifies the copy, not the original
}

// Pass by pointer the function gets the ADDRESS
func doublePointer(n *int) {
	*n *= 2 // dereference and modify the original
}

// Returning a pointer Go allocates on heap if needed
func newInt(val int) *int {
	x := val   // x lives on the heap (escapes this function)
	return &x  // perfectly safe in Go GC manages lifetime
}

// Nil pointer safety
func safeLength(s *string) int {
	if s == nil {
		return 0
	}
	return len(*s)
}

// Struct example pointer avoids copying large data
type Config struct {
	Host    string
	Port    int
	Debug   bool
}

func enableDebug(cfg *Config) {
	cfg.Debug = true // Go auto-dereferences: cfg.Debug == (*cfg).Debug
}

func main() {
	// 1. Basics: & and *
	fmt.Println("--- Basics: & and * ---")
	x := 42
	p := &x // p is a *int pointing to x
	fmt.Printf("x = %d | &x = %p | p = %p | *p = %d\n", x, &x, p, *p)

	*p = 100 // modify x through the pointer
	fmt.Printf("After *p = 100: x = %d\n", x)

	// 2. Pass by value vs pointer
	fmt.Println("\n--- Value vs Pointer ---")
	n := 10
	doubleValue(n)
	fmt.Printf("After doubleValue(n): n = %d (unchanged copy)\n", n)

	doublePointer(&n)
	fmt.Printf("After doublePointer(&n): n = %d (modified via pointer)\n", n)

	// 3. Returning pointers safe in Go
	fmt.Println("\n--- Returning Pointers ---")
	ptr := newInt(7)
	fmt.Printf("newInt(7) returned *int at %p with value %d\n", ptr, *ptr)

	// 4. Nil pointers
	fmt.Println("\n--- Nil Pointers ---")
	var nilPtr *int
	fmt.Printf("nilPtr = %v (zero value of a pointer is nil)\n", nilPtr)
	if nilPtr == nil {
		fmt.Println("✓ Always check for nil before dereferencing!")
	}

	msg := "hello"
	fmt.Printf("safeLength(&msg) = %d\n", safeLength(&msg))
	fmt.Printf("safeLength(nil) = %d\n", safeLength(nil))

	// 5. Pointers to structs
	fmt.Println("\n--- Pointers to Structs ---")
	cfg := Config{Host: "localhost", Port: 8080, Debug: false}
	fmt.Printf("Before: %+v\n", cfg)
	enableDebug(&cfg)
	fmt.Printf("After:  %+v\n", cfg)

	// 6. new() built-in allocates zeroed memory, returns pointer
	fmt.Println("\n--- new() Built-in ---")
	numPtr := new(int) // allocates an int, sets to zero value
	fmt.Printf("new(int) → %p, value = %d\n", numPtr, *numPtr)
	*numPtr = 55
	fmt.Printf("After assignment: value = %d\n", *numPtr)

	// 7. No pointer arithmetic!
	fmt.Println("\n--- No Pointer Arithmetic ---")
	fmt.Println("In C++: *(ptr + 1) walks to next memory slot")
	fmt.Println("In Go:  COMPILE ERROR pointer arithmetic is forbidden")
	fmt.Println("This eliminates buffer overflows and out-of-bounds access")

	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. & gets address, * dereferences same as C++ syntax")
	fmt.Println("2. No arithmetic pointers are safe handles, not offsets")
	fmt.Println("3. Returning &localVar is safe Go escapes it to the heap")
	fmt.Println("4. Pass pointer for mutation or large structs; value for small/read-only")
	fmt.Println("5. Zero value of any pointer is nil always check before use")
}
