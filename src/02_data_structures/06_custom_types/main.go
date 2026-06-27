/*
What this teaches:
    Type definitions vs type aliases, attaching methods to custom types,
    implementing multiple interfaces, embedding for composition, and the
    Stringer interface for custom string representation.

Beginner analogy:
    "Custom types are like creating your own LEGO brick shapes you start
     with a standard brick (int, string) but give it a new name and new
     abilities (methods)."

C++ comparison:
    "typedef + methods. Go interfaces compose like traits/concepts. A custom
     type based on int is a NEW type it doesn't inherit int's methods, but
     you can convert between them. No operator overloading."

Interview relevance:
    Type definitions vs aliases, method sets, satisfying multiple interfaces,
    the Stringer pattern, type-safe enums with iota, and when to use embedding.
*/

package main

import (
	"fmt"
	"strings"
)

// 1. Type definition creates a NEW type (not an alias)
type Celsius float64
type Fahrenheit float64

// Methods on custom types
func (c Celsius) ToFahrenheit() Fahrenheit {
	return Fahrenheit(c*9/5 + 32)
}

func (f Fahrenheit) ToCelsius() Celsius {
	return Celsius((f - 32) * 5 / 9)
}

// Stringer interface controls how fmt prints your type
func (c Celsius) String() string {
	return fmt.Sprintf("%.1f°C", float64(c))
}

func (f Fahrenheit) String() string {
	return fmt.Sprintf("%.1f°F", float64(f))
}

// 2. Type-safe enum with iota
type Direction int

const (
	North Direction = iota // 0
	East                   // 1
	South                  // 2
	West                   // 3
)

func (d Direction) String() string {
	names := [...]string{"North", "East", "South", "West"}
	if d < North || d > West {
		return "Unknown"
	}
	return names[d]
}

func (d Direction) Opposite() Direction {
	return (d + 2) % 4
}

// 3. Custom string type with methods
type Email string

func (e Email) Domain() string {
	parts := strings.Split(string(e), "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

func (e Email) IsValid() bool {
	return strings.Contains(string(e), "@") && strings.Contains(string(e), ".")
}

// 4. Implementing multiple interfaces
type Printable interface {
	Print()
}

type Saveable interface {
	Save() string
}

type Document struct {
	Title   string
	Content string
}

func (d Document) Print() {
	fmt.Printf("  [PRINT] %s: %s\n", d.Title, d.Content)
}

func (d Document) Save() string {
	return fmt.Sprintf("saved:%s", d.Title)
}

func (d Document) String() string {
	return fmt.Sprintf("Document(%q, %d chars)", d.Title, len(d.Content))
}

// 5. Embedding for composition
type Timestamp struct {
	CreatedAt string
	UpdatedAt string
}

func (t Timestamp) Age() string {
	return "created: " + t.CreatedAt
}

type Article struct {
	Timestamp // embedded promotes CreatedAt, UpdatedAt, Age()
	Title     string
	Body      string
}

// 6. Type alias vs type definition
type MyInt = int // ALIAS same type, no new methods allowed
// type NewInt int // DEFINITION new type, can have methods

func main() {
	// Custom numeric types with methods
	fmt.Println("--- Custom Types (Celsius/Fahrenheit) ---")
	boiling := Celsius(100)
	fmt.Printf("%v = %v\n", boiling, boiling.ToFahrenheit())
	body := Fahrenheit(98.6)
	fmt.Printf("%v = %v\n", body, body.ToCelsius())

	// Type safety can't accidentally mix them
	// var temp Celsius = Fahrenheit(72) // COMPILE ERROR: different types!
	fmt.Println("Can't assign Fahrenheit to Celsius type safety!")

	// Enum with iota
	fmt.Println("\n--- Enum with iota ---")
	dir := North
	fmt.Printf("Direction: %v (value: %d)\n", dir, dir)
	fmt.Printf("Opposite of %v: %v\n", dir, dir.Opposite())
	for d := North; d <= West; d++ {
		fmt.Printf("  %v = %d\n", d, d)
	}

	// Custom string type
	fmt.Println("\n--- Custom String Type (Email) ---")
	email := Email("gopher@golang.org")
	fmt.Printf("Email: %s\n", email)
	fmt.Printf("Domain: %s\n", email.Domain())
	fmt.Printf("Valid: %t\n", email.IsValid())
	bad := Email("invalid")
	fmt.Printf("Bad email valid: %t\n", bad.IsValid())

	// Multiple interfaces
	fmt.Println("\n--- Multiple Interfaces ---")
	doc := Document{Title: "README", Content: "Hello World"}
	fmt.Println(doc) // uses String()
	doc.Print()      // Printable interface
	fmt.Printf("  Save result: %s\n", doc.Save())

	// Prove it satisfies both interfaces
	var p Printable = doc
	var s Saveable = doc
	p.Print()
	fmt.Println("  Saveable:", s.Save())

	// Embedding for composition
	fmt.Println("\n--- Embedding (Composition) ---")
	article := Article{
		Timestamp: Timestamp{CreatedAt: "2024-01-15", UpdatedAt: "2024-06-01"},
		Title:     "Go Generics",
		Body:      "Generics arrived in Go 1.18...",
	}
	fmt.Printf("Title: %s\n", article.Title)
	fmt.Printf("Created: %s\n", article.CreatedAt) // promoted field
	fmt.Printf("Age: %s\n", article.Age())         // promoted method

	// Type alias vs definition
	fmt.Println("\n--- Alias vs Definition ---")
	var x MyInt = 42
	var y int = x // works same underlying type (alias)
	fmt.Printf("Alias MyInt: %d == int: %d\n", x, y)

	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. Type definitions create NEW types enable methods, prevent mixing")
	fmt.Println("2. Implement Stringer for human-readable fmt output")
	fmt.Println("3. iota + custom int type = type-safe enums")
	fmt.Println("4. Embedding promotes fields/methods composition over inheritance")
	fmt.Println("5. Type alias (=) shares identity; definition creates a new one")
}
