/*
What this teaches:
    Interface declaration, implicit satisfaction (duck typing), the empty
    interface (any), type assertions, type switches, and the internal
    representation of interfaces (two-word: type pointer + data pointer).

Beginner analogy:
    "An interface is a job description anyone who can do the job qualifies.
     You don't need to 'apply' (no implements keyword). If you have the right
     skills (methods), you're hired."

C++ comparison:
    "Implicit satisfaction = duck typing with compile-time checking. No vtable
     pointer stored IN the object itself the interface value holds it. Like
     concepts (C++20) enforced at assignment, not declaration."

Interview relevance:
    nil interface vs interface holding nil, the two-word layout, why io.Reader
    is powerful, interface composition, and the any type.
*/

package main

import (
	"fmt"
	"math"
)

// 1. Interface declaration just method signatures
type Shape interface {
	Area() float64
	Perimeter() float64
}

// 2. Another interface single method (very Go-idiomatic)
type Stringer interface {
	String() string
}

// 3. Concrete types that implicitly satisfy Shape
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64      { return r.Width * r.Height }
func (r Rectangle) Perimeter() float64  { return 2 * (r.Width + r.Height) }
func (r Rectangle) String() string {
	return fmt.Sprintf("Rect(%.1f×%.1f)", r.Width, r.Height)
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64  { return 2 * math.Pi * c.Radius }
func (c Circle) String() string {
	return fmt.Sprintf("Circle(r=%.1f)", c.Radius)
}

// 4. Interface composition building bigger interfaces from smaller ones
type Describer interface {
	Shape
	Stringer
}

// 5. Function accepting an interface polymorphism
func printShape(s Shape) {
	fmt.Printf("  Area=%.2f Perimeter=%.2f\n", s.Area(), s.Perimeter())
}

func describeAll(items []Describer) {
	for _, item := range items {
		fmt.Printf("  %s → Area=%.2f\n", item.String(), item.Area())
	}
}

// 6. The empty interface: any (alias for interface{})
func printAnything(val any) {
	fmt.Printf("  Type: %-12T | Value: %v\n", val, val)
}

func main() {
	// Implicit satisfaction no "implements" keyword needed
	fmt.Println("--- Implicit Interface Satisfaction ---")
	var s Shape = Rectangle{Width: 5, Height: 3}
	printShape(s)
	s = Circle{Radius: 4}
	printShape(s)

	// Polymorphism with slices
	fmt.Println("\n--- Polymorphic Slice ---")
	shapes := []Shape{
		Rectangle{10, 5},
		Circle{7},
		Rectangle{3, 3},
	}
	for _, sh := range shapes {
		printShape(sh)
	}

	// Interface composition
	fmt.Println("\n--- Composed Interface (Describer) ---")
	items := []Describer{
		Rectangle{4, 6},
		Circle{2.5},
	}
	describeAll(items)

	// empty interface: any
	fmt.Println("\n--- Empty Interface (any) ---")
	printAnything(42)
	printAnything("hello")
	printAnything(3.14)
	printAnything([]int{1, 2, 3})

	// Type assertion extract concrete type from interface
	fmt.Println("\n--- Type Assertion ---")
	var val any = "Go is great"
	str, ok := val.(string) // comma-ok pattern
	fmt.Printf("val.(string) → %q, ok=%t\n", str, ok)

	num, ok := val.(int)
	fmt.Printf("val.(int) → %d, ok=%t (failed safely)\n", num, ok)

	// Type switch handle multiple types cleanly
	fmt.Println("\n--- Type Switch ---")
	things := []any{42, "hello", 3.14, true, Rectangle{1, 2}}
	for _, thing := range things {
		switch v := thing.(type) {
		case int:
			fmt.Printf("  int: %d\n", v)
		case string:
			fmt.Printf("  string: %q (len=%d)\n", v, len(v))
		case float64:
			fmt.Printf("  float64: %.2f\n", v)
		case Shape:
			fmt.Printf("  Shape: area=%.2f\n", v.Area())
		default:
			fmt.Printf("  unknown: %v\n", v)
		}
	}

	// nil interface pitfall
	fmt.Println("\n--- Nil Interface Pitfall ---")
	var sh Shape // nil interface (no type, no value)
	fmt.Printf("nil interface: %v, ==nil? %t\n", sh, sh == nil)

	var rp *Rectangle // nil pointer of concrete type
	sh = rp           // interface now holds (type=*Rectangle, value=nil)
	fmt.Printf("interface holding nil ptr: %v, ==nil? %t\n", sh, sh == nil)
	fmt.Println("⚠ An interface with a nil concrete value is NOT nil!")

	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. Interfaces are satisfied implicitly just implement the methods")
	fmt.Println("2. Small interfaces (1-2 methods) are idiomatic: io.Reader, error")
	fmt.Println("3. Type assertions + switches let you recover concrete types")
	fmt.Println("4. Interface value = {type pointer, data pointer} two words")
	fmt.Println("5. nil interface ≠ interface holding nil a classic gotcha")
}
