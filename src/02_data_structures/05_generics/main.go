/*
What this teaches:
    Type parameters (generics), constraints (any, comparable, custom),
    type sets using union (~int | ~float64), and building generic data
    structures like Stack[T].

Beginner analogy:
    "Generics = one recipe that works for any ingredient. A Stack recipe works
     for a stack of ints, strings, or structs you don't rewrite it each time."

C++ comparison:
    "Similar to templates but no monomorphization uses GCShape stenciling.
     Fewer binary bloat. Constraints ≈ C++20 concepts but enforced at
     instantiation. No SFINAE, no specialization."

Interview relevance:
    How do Go generics differ from C++ templates? What is comparable? When to
    use generics vs interfaces? What are type sets? Performance implications?
*/

package main

import (
	"cmp"
	"fmt"
	"strings"
)

// 1. Basic generic function works for any ordered type
func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// 2. Custom constraint using type sets
type Number interface {
	~int | ~int32 | ~int64 | ~float32 | ~float64
}

func Sum[T Number](nums []T) T {
	var total T
	for _, n := range nums {
		total += n
	}
	return total
}

// 3. Generic data structure: Stack[T]
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	last := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return last, true
}

func (s *Stack[T]) Peek() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

func (s *Stack[T]) Len() int {
	return len(s.items)
}

// 4. Generic Map function (transform a slice)
func Map[T any, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

// 5. Generic Filter function
func Filter[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// 6. Using comparable constraint (required for map keys / == operator)
func Contains[T comparable](slice []T, target T) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

func main() {
	// Basic generic function
	fmt.Println("--- Generic Min ---")
	fmt.Printf("Min(3, 7) = %d\n", Min(3, 7))
	fmt.Printf("Min(3.14, 2.71) = %.2f\n", Min(3.14, 2.71))
	fmt.Printf("Min(\"apple\", \"banana\") = %q\n", Min("apple", "banana"))

	// Custom Number constraint
	fmt.Println("\n--- Custom Constraint (Number) ---")
	ints := []int{1, 2, 3, 4, 5}
	floats := []float64{1.1, 2.2, 3.3}
	fmt.Printf("Sum(ints) = %d\n", Sum(ints))
	fmt.Printf("Sum(floats) = %.1f\n", Sum(floats))

	// Generic Stack
	fmt.Println("\n--- Stack[int] ---")
	var intStack Stack[int]
	intStack.Push(10)
	intStack.Push(20)
	intStack.Push(30)
	fmt.Printf("Len: %d\n", intStack.Len())
	if val, ok := intStack.Pop(); ok {
		fmt.Printf("Pop: %d\n", val)
	}
	if val, ok := intStack.Peek(); ok {
		fmt.Printf("Peek: %d\n", val)
	}

	fmt.Println("\n--- Stack[string] ---")
	var strStack Stack[string]
	strStack.Push("hello")
	strStack.Push("world")
	if val, ok := strStack.Pop(); ok {
		fmt.Printf("Pop: %q\n", val)
	}

	// Generic Map (transform)
	fmt.Println("\n--- Generic Map ---")
	words := []string{"hello", "world", "go"}
	upper := Map(words, strings.ToUpper)
	fmt.Printf("Map(ToUpper): %v\n", upper)

	lengths := Map(words, func(s string) int { return len(s) })
	fmt.Printf("Map(len): %v\n", lengths)

	// Generic Filter
	fmt.Println("\n--- Generic Filter ---")
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	fmt.Printf("Filter(even): %v\n", evens)

	// Contains with comparable
	fmt.Println("\n--- Contains (comparable) ---")
	fmt.Printf("Contains([1,2,3], 2) = %t\n", Contains([]int{1, 2, 3}, 2))
	fmt.Printf("Contains([\"a\",\"b\"], \"c\") = %t\n", Contains([]string{"a", "b"}, "c"))

	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. [T constraint] after func/type name declares type parameters")
	fmt.Println("2. 'any' = no constraint. 'comparable' = supports == and !=")
	fmt.Println("3. Custom constraints use interface with type sets (~int | ~float64)")
	fmt.Println("4. ~ means underlying type includes user-defined types based on int")
	fmt.Println("5. Generics reduce code duplication without sacrificing type safety")
}
