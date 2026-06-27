/*
What this teaches:
    Functions in Go multiple returns, named returns, variadic arguments,
    first-class functions (functions as values), and closures.

Beginner analogy:
    "A function is a vending machine you put in coins (parameters), press a
     button, and get a drink (return value). Some machines give you change too
     (multiple returns)."

C++ comparison:
    "Multiple returns replace out-params and std::tuple. Closures capture by
     reference (pointer to stack/heap-escaped variable). No function overloading."

Interview relevance:
    Multiple returns are idiomatic for error handling. Closures appear in
    goroutine patterns. Named returns + bare return are a common code-review
    discussion. Variadic args power fmt.Println itself.
*/

package main

import (
	"fmt"
	"strings"
)

// 1. Basic function two params, one return
func add(a, b int) int {
	return a + b
}

// 2. Multiple return values the Go idiom for error handling
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

// 3. Named return values documents what's coming back
func swap(x, y string) (first, second string) {
	first = y
	second = x
	return // bare return uses named values
}

// 4. Variadic function accepts any number of ints
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// 5. Functions as first-class values
func applyOp(a, b int, op func(int, int) int) int {
	return op(a, b)
}

// 6. Function returning a function (closure / factory)
func makeGreeter(prefix string) func(string) string {
	// prefix is "captured" by the returned closure
	return func(name string) string {
		return prefix + " " + name + "!"
	}
}

// 7. Closure with mutable state counter
func newCounter() func() int {
	count := 0 // lives on heap because closure escapes
	return func() int {
		count++
		return count
	}
}

func main() {
	fmt.Println("--- Basic Function ---")
	fmt.Printf("add(3, 5) = %d\n", add(3, 5))

	fmt.Println("\n--- Multiple Returns ---")
	result, err := divide(10, 3)
	fmt.Printf("divide(10, 3) = %.4f, err = %v\n", result, err)
	_, err = divide(10, 0)
	fmt.Printf("divide(10, 0) err = %v\n", err)

	fmt.Println("\n--- Named Returns ---")
	a, b := swap("hello", "world")
	fmt.Printf("swap(\"hello\", \"world\") = %q, %q\n", a, b)

	fmt.Println("\n--- Variadic Function ---")
	fmt.Printf("sum(1,2,3,4,5) = %d\n", sum(1, 2, 3, 4, 5))
	nums := []int{10, 20, 30}
	fmt.Printf("sum(slice...) = %d\n", sum(nums...)) // spread a slice

	fmt.Println("\n--- First-Class Functions ---")
	multiply := func(a, b int) int { return a * b }
	fmt.Printf("applyOp(4, 5, multiply) = %d\n", applyOp(4, 5, multiply))

	fmt.Println("\n--- Closures ---")
	hello := makeGreeter("Hello")
	hey := makeGreeter("Hey")
	fmt.Println(hello("Gopher"))
	fmt.Println(hey("World"))

	fmt.Println("\n--- Closure with State ---")
	counter := newCounter()
	fmt.Printf("counter() = %d\n", counter())
	fmt.Printf("counter() = %d\n", counter())
	fmt.Printf("counter() = %d\n", counter())

	fmt.Println("\n--- Practical: strings.Map with func ---")
	shout := strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' {
			return r - 32
		}
		return r
	}, "hello, go!")
	fmt.Println("shout:", shout)

	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. Multiple returns are idiomatic especially (result, error)")
	fmt.Println("2. Closures capture variables by reference (pointer)")
	fmt.Println("3. No overloading use variadic or option structs instead")
	fmt.Println("4. Functions are values pass them, return them, store them")
}
