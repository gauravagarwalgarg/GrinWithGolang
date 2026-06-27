/*
What this teaches:
    Variables (var, :=, const), fundamental types (int, float64, string, bool,
    byte, rune), type inference, and Go's zero-value guarantee.

Beginner analogy:
    "Every box has a default label even if empty Go never leaves memory
     uninitialized. An int box starts at 0, a string box starts at '', a bool
     box starts at false."

C++ comparison:
    "No uninitialized memory every variable has a zero value. No implicit
     conversions between numeric types. := replaces auto with even less ceremony."

Interview relevance:
    Asked frequently: What is the zero value of a slice? a map? a pointer?
    Also: difference between var and :=, why Go has no implicit type conversion,
    and how constants are untyped until used.
*/

package main

import "fmt"

func main() {
	// 1. var declaration explicit type
	fmt.Println("--- var with explicit type ---")
	var x int = 42
	var greeting string = "Hello"
	fmt.Printf("x = %d (type %T)\n", x, x)
	fmt.Printf("greeting = %q (type %T)\n", greeting, greeting)

	// 2. var with type inference
	fmt.Println("\n--- var with inference ---")
	var pi = 3.14159 // inferred as float64
	var flag = true   // inferred as bool
	fmt.Printf("pi = %f (type %T)\n", pi, pi)
	fmt.Printf("flag = %t (type %T)\n", flag, flag)

	// 3. Short declaration := (most common inside functions)
	fmt.Println("\n--- Short declaration := ---")
	name := "Gopher"    // inferred string
	age := 15           // inferred int
	height := 1.75      // inferred float64
	fmt.Printf("name=%s age=%d height=%.2f\n", name, age, height)

	// 4. Zero values Go's safety net
	fmt.Println("\n--- Zero Values ---")
	var zeroInt int
	var zeroFloat float64
	var zeroString string
	var zeroBool bool
	var zeroPointer *int
	fmt.Printf("int: %d | float64: %f | string: %q | bool: %t | pointer: %v\n",
		zeroInt, zeroFloat, zeroString, zeroBool, zeroPointer)

	// 5. Constants immutable, untyped until used
	fmt.Println("\n--- Constants ---")
	const daysInWeek = 7           // untyped integer constant
	const appName string = "GrinGo" // typed constant
	const (
		statusOK    = 200
		statusNotFound = 404
	)
	fmt.Printf("Days: %d | App: %s | OK: %d | NotFound: %d\n",
		daysInWeek, appName, statusOK, statusNotFound)

	// 6. Fundamental types
	fmt.Println("\n--- Fundamental Types ---")
	var b byte = 'A'        // byte = uint8
	var r rune = '🚀'       // rune = int32 (Unicode code point)
	var i int = 100         // platform-sized (64-bit on modern systems)
	var f float64 = 2.718   // IEEE-754 double
	var s string = "Go 🎉"  // UTF-8 encoded
	fmt.Printf("byte: %d (%c) | rune: %d (%c) | int: %d | float64: %.3f | string: %s\n",
		b, b, r, r, i, f, s)

	// 7. No implicit conversion must be explicit
	fmt.Println("\n--- No Implicit Conversion ---")
	var small int32 = 10
	var big int64 = int64(small) // explicit cast required
	fmt.Printf("int32(%d) → int64(%d)\n", small, big)

	// 8. Multiple assignment
	fmt.Println("\n--- Multiple Assignment ---")
	a, b2, c := 1, "two", 3.0
	fmt.Printf("a=%v b2=%v c=%v\n", a, b2, c)

	// Swap without temp
	a, _ = 99, 0 // blank identifier _ discards a value
	fmt.Printf("a after swap: %d\n", a)

	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. Every variable has a zero value no garbage memory")
	fmt.Println("2. := is shorthand inside functions; var for package level")
	fmt.Println("3. No implicit type conversion safety over convenience")
	fmt.Println("4. Constants can be untyped they adapt to context")
}
