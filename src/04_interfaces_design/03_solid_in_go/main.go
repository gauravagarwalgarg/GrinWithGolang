/*
What this teaches:
    All 5 SOLID principles demonstrated in idiomatic Go. Shows how Go achieves
    these design principles without class hierarchies, using interfaces, embedding,
    and package-level separation.

Beginner analogy:
    "SOLID is like building with LEGO each brick does one thing (SRP), new pieces
     snap on without modifying old ones (OCP), any brick fits if it has the right
     studs (LSP), connectors are minimal (ISP), and you build on pegs not glue (DIP)."

C++ comparison:
    "Same principles, radically different mechanics. C++ uses abstract classes,
     virtual dispatch, and templates. Go uses implicit interfaces, embedding, and
     package boundaries. The spirit is identical; the ceremony is minimal."

Interview relevance:
    Senior Go interviews ask how SOLID applies without class hierarchies. Knowing
    that Go's implicit interfaces inherently encourage ISP and DIP, and that package
    boundaries enforce SRP, demonstrates design maturity.
*/

package main

import (
	"fmt"
	"math"
)

// ============================================================
// S Single Responsibility Principle
// Each type has one reason to change.
// ============================================================

// BAD: User handles both data AND persistence
// GOOD: Separate concerns into distinct types

type User struct {
	ID    int
	Name  string
	Email string
}

// UserValidator only validates
type UserValidator struct{}

func (v *UserValidator) Validate(u User) error {
	if u.Name == "" {
		return fmt.Errorf("name is required")
	}
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	return nil
}

// UserRepository only handles persistence
type UserRepository interface {
	Save(u User) error
}

type InMemoryUserRepo struct {
	users []User
}

func (r *InMemoryUserRepo) Save(u User) error {
	r.users = append(r.users, u)
	fmt.Printf("  [SRP] Saved user: %s\n", u.Name)
	return nil
}

// ============================================================
// O Open/Closed Principle
// Open for extension, closed for modification.
// ============================================================

type Shape interface {
	Area() float64
	Name() string
}

type Circle struct{ Radius float64 }

func (c Circle) Area() float64 { return math.Pi * c.Radius * c.Radius }
func (c Circle) Name() string  { return "Circle" }

type Rectangle struct{ Width, Height float64 }

func (r Rectangle) Area() float64 { return r.Width * r.Height }
func (r Rectangle) Name() string  { return "Rectangle" }

// Adding Triangle requires NO changes to existing code
type Triangle struct{ Base, Height float64 }

func (t Triangle) Area() float64 { return 0.5 * t.Base * t.Height }
func (t Triangle) Name() string  { return "Triangle" }

// This function never changes new shapes just implement Shape
func TotalArea(shapes []Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

// ============================================================
// L Liskov Substitution Principle
// Subtypes must honor the contract of their parent type.
// ============================================================

type Bird interface {
	Move() string
}

type Sparrow struct{}

func (s Sparrow) Move() string { return "flies through the air" }

type Penguin struct{}

func (p Penguin) Move() string { return "waddles on land" } // Still moves honors contract

// NOT: func (p Penguin) Fly() would violate LSP if Bird required Fly()

// ============================================================
// I Interface Segregation Principle
// Clients should not depend on methods they don't use.
// ============================================================

// BAD: One fat interface
// type Worker interface { Work(); Eat(); Sleep(); Code(); Deploy() }

// GOOD: Small, focused interfaces (like stdlib io.Reader, io.Writer)
type Printer interface {
	Print(doc string)
}

type Scanner interface {
	Scan() string
}

// Only implement what you need
type SimplePrinter struct{}

func (p SimplePrinter) Print(doc string) {
	fmt.Printf("  [ISP] Printing: %s\n", doc)
}

// A multifunction device composes both
type MultiFunctionDevice struct {
	SimplePrinter
}

func (m MultiFunctionDevice) Scan() string {
	return "scanned document"
}

// ============================================================
// D Dependency Inversion Principle
// High-level modules depend on abstractions, not details.
// ============================================================

type DataStore interface {
	Get(key string) (string, bool)
}

type CacheStore struct {
	data map[string]string
}

func (c *CacheStore) Get(key string) (string, bool) {
	v, ok := c.data[key]
	return v, ok
}

// High-level: depends on DataStore interface, not CacheStore concrete type
type AppService struct {
	store DataStore
}

func NewAppService(ds DataStore) *AppService {
	return &AppService{store: ds}
}

func (a *AppService) Lookup(key string) string {
	if v, ok := a.store.Get(key); ok {
		return v
	}
	return "<not found>"
}

func main() {
	fmt.Println("=== SOLID Principles in Go ===")

	// S SRP
	fmt.Println("\n--- S: Single Responsibility ---")
	validator := &UserValidator{}
	repo := &InMemoryUserRepo{}
	u := User{ID: 1, Name: "Alice", Email: "alice@go.dev"}
	if err := validator.Validate(u); err == nil {
		repo.Save(u)
	}

	// O OCP
	fmt.Println("\n--- O: Open/Closed ---")
	shapes := []Shape{
		Circle{Radius: 5},
		Rectangle{Width: 3, Height: 4},
		Triangle{Base: 6, Height: 3}, // Added without modifying TotalArea
	}
	for _, s := range shapes {
		fmt.Printf("  %s area: %.2f\n", s.Name(), s.Area())
	}
	fmt.Printf("  Total area: %.2f\n", TotalArea(shapes))

	// L LSP
	fmt.Println("\n--- L: Liskov Substitution ---")
	birds := []Bird{Sparrow{}, Penguin{}}
	for _, b := range birds {
		fmt.Printf("  Bird %T: %s\n", b, b.Move())
	}

	// I ISP
	fmt.Println("\n--- I: Interface Segregation ---")
	var p Printer = SimplePrinter{}
	p.Print("report.pdf")
	mfd := MultiFunctionDevice{}
	mfd.Print("invoice.pdf")
	fmt.Printf("  Scanned: %s\n", mfd.Scan())

	// D DIP
	fmt.Println("\n--- D: Dependency Inversion ---")
	cache := &CacheStore{data: map[string]string{"host": "localhost", "port": "8080"}}
	app := NewAppService(cache)
	fmt.Printf("  host = %s\n", app.Lookup("host"))
	fmt.Printf("  missing = %s\n", app.Lookup("db"))

	fmt.Println("\n--- Summary ---")
	fmt.Println("Go achieves SOLID through:")
	fmt.Println("  S → Package boundaries + single-purpose types")
	fmt.Println("  O → Interfaces allow extension without modification")
	fmt.Println("  L → Implicit interfaces enforce behavioral contracts")
	fmt.Println("  I → Small interfaces (io.Reader = 1 method)")
	fmt.Println("  D → Constructor injection of interface dependencies")
}
