/*
What this teaches:
    Struct definition, methods (value vs pointer receiver), struct embedding
    (composition over inheritance), factory functions, and struct tags.

Beginner analogy:
    "A struct is a form with fields to fill in like a registration form with
     Name, Email, Age. Methods are actions the form can perform on itself."

C++ comparison:
    "Structs replace classes. No constructors use factory functions (NewXxx).
     Embedding ≈ composition, NOT inheritance. No virtual dispatch on structs."

Interview relevance:
    Value vs pointer receivers, when to use each, struct embedding vs interfaces,
    zero-value usefulness, and the NewXxx factory pattern.
*/

package main

import (
	"fmt"
	"math"
)

// Basic struct definition
type Point struct {
	X, Y float64
}

// Value receiver does NOT modify the original
func (p Point) Distance(other Point) float64 {
	dx := p.X - other.X
	dy := p.Y - other.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// Stringer interface controls fmt output
func (p Point) String() string {
	return fmt.Sprintf("(%0.1f, %0.1f)", p.X, p.Y)
}

// Pointer receiver CAN modify the original
func (p *Point) Scale(factor float64) {
	p.X *= factor
	p.Y *= factor
}

// Struct with unexported fields + factory function
type User struct {
	Name  string
	Email string
	age   int // unexported only this package can set it
}

// Factory function replaces constructors
func NewUser(name, email string, age int) *User {
	return &User{
		Name:  name,
		Email: email,
		age:   age,
	}
}

func (u User) Summary() string {
	return fmt.Sprintf("%s <%s> (age %d)", u.Name, u.Email, u.age)
}

// --- Embedding (Composition) ---
type Engine struct {
	Horsepower int
	Type       string
}

func (e Engine) Start() string {
	return fmt.Sprintf("%s engine (%d HP) started!", e.Type, e.Horsepower)
}

type Car struct {
	Engine // embedded Car "has-a" Engine (promotes methods)
	Brand  string
	Model  string
}

func (c Car) Info() string {
	return fmt.Sprintf("%s %s with %d HP %s engine",
		c.Brand, c.Model, c.Horsepower, c.Type) // promoted fields
}

// --- Method sets and interfaces ---
type Shape interface {
	Area() float64
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func printArea(s Shape) {
	fmt.Printf("  Area = %.2f\n", s.Area())
}

func main() {
	// 1. Creating structs
	fmt.Println("--- Creating Structs ---")
	p1 := Point{X: 3, Y: 4}
	p2 := Point{} // zero value: {0, 0}
	fmt.Printf("p1 = %v | p2 = %v\n", p1, p2)

	// 2. Methods
	fmt.Println("\n--- Value Receiver ---")
	origin := Point{0, 0}
	fmt.Printf("Distance from %v to %v = %.2f\n", p1, origin, p1.Distance(origin))

	fmt.Println("\n--- Pointer Receiver ---")
	fmt.Printf("Before Scale: p1 = %v\n", p1)
	p1.Scale(2) // Go auto-takes address: (&p1).Scale(2)
	fmt.Printf("After Scale(2): p1 = %v\n", p1)

	// 3. Factory function
	fmt.Println("\n--- Factory Function ---")
	user := NewUser("Alice", "alice@go.dev", 30)
	fmt.Println(user.Summary())

	// 4. Struct embedding
	fmt.Println("\n--- Embedding (Composition) ---")
	car := Car{
		Engine: Engine{Horsepower: 200, Type: "V6"},
		Brand:  "Toyota",
		Model:  "Camry",
	}
	fmt.Println(car.Info())
	fmt.Println(car.Start()) // promoted from Engine

	// 5. Structs satisfying interfaces
	fmt.Println("\n--- Structs as Interfaces ---")
	shapes := []Shape{
		Rectangle{Width: 5, Height: 3},
		Circle{Radius: 4},
	}
	for _, s := range shapes {
		fmt.Printf("  %T:", s)
		printArea(s)
	}

	// 6. Anonymous structs (quick one-off use)
	fmt.Println("\n--- Anonymous Struct ---")
	config := struct {
		Host string
		Port int
	}{Host: "localhost", Port: 8080}
	fmt.Printf("config = %+v\n", config)

	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. Value receiver = read-only, pointer receiver = mutating")
	fmt.Println("2. Use NewXxx factories instead of constructors")
	fmt.Println("3. Embedding promotes fields + methods (composition, NOT inheritance)")
	fmt.Println("4. Structs implicitly satisfy interfaces no 'implements' keyword")
	fmt.Println("5. Zero-value structs should be usable (design for it)")
}
