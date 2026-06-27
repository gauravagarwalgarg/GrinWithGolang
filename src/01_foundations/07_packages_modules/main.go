/*
What this teaches:
    Package organization, visibility rules (uppercase = exported), Go modules
    (go.mod), and how to structure a multi-package project. We simulate a
    "mathutil" package inline with explanatory comments.

Beginner analogy:
    "Packages are like departments in a company Accounting, Engineering, HR.
     Each has its own responsibilities. Uppercase names are 'public-facing staff'
     (exported), lowercase names are 'internal workers' (unexported)."

C++ comparison:
    "Uppercase = public, lowercase = private. No header/impl split the
     compiler reads the source directly. go.mod replaces CMakeLists.txt for
     dependency management. No #include guards needed."

Interview relevance:
    Common questions: What makes a name exported? How does go.mod work? What's
    the difference between a module and a package? How do you organize a large
    Go project?
*/

package main

import (
	"fmt"
	"math"
)

// --- Simulating a "mathutil" package inline ---
// In a real project, this would live in a separate directory:
//   mymodule/mathutil/mathutil.go
//
// With go.mod:
//   module github.com/yourname/mymodule
//   go 1.22

// Exported function starts with uppercase
// In a real package, other packages can call mathutil.Add(...)
func Add(a, b float64) float64 {
	return a + b
}

// Exported function
func Multiply(a, b float64) float64 {
	return a * b
}

// unexported helper starts with lowercase
// Only accessible within the same package
func clamp(val, min, max float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// Exported constant accessible from other packages
const Pi = 3.14159265358979

// unexported constant internal to this package
const maxIterations = 1000

// Exported type with exported and unexported fields
type Circle struct {
	Radius float64 // exported other packages can access
	label  string  // unexported only this package can see
}

// Exported method on Circle
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Exported method
func (c Circle) Circumference() float64 {
	return 2 * math.Pi * c.Radius
}

// unexported method helper, not part of public API
func (c Circle) describe() string {
	return fmt.Sprintf("circle(r=%.2f, label=%s)", c.Radius, c.label)
}

func main() {
	fmt.Println("--- Visibility Rules ---")
	fmt.Println("Uppercase = Exported (public):  Add, Multiply, Circle, Pi")
	fmt.Println("lowercase = unexported (private): clamp, maxIterations, label")

	fmt.Println("\n--- Using 'Exported' Functions ---")
	fmt.Printf("Add(3, 4) = %.1f\n", Add(3, 4))
	fmt.Printf("Multiply(3, 4) = %.1f\n", Multiply(3, 4))
	fmt.Printf("clamp(150, 0, 100) = %.1f (unexported helper)\n", clamp(150, 0, 100))

	fmt.Println("\n--- Using an Exported Type ---")
	c := Circle{Radius: 5.0, label: "unit-circle"}
	fmt.Printf("Circle: %+v\n", c)
	fmt.Printf("Area: %.4f\n", c.Area())
	fmt.Printf("Circumference: %.4f\n", c.Circumference())
	fmt.Printf("describe(): %s (unexported method)\n", c.describe())

	fmt.Println("\n--- go.mod Explained ---")
	fmt.Println("go.mod defines your module and its dependencies:")
	fmt.Println("  module github.com/yourname/project")
	fmt.Println("  go 1.22")
	fmt.Println("  require github.com/some/lib v1.2.3")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  go mod init <module-path>  create go.mod")
	fmt.Println("  go mod tidy               add/remove deps")
	fmt.Println("  go get <pkg>@<version>    add a dependency")

	fmt.Println("\n--- Project Layout Convention ---")
	fmt.Println("  myproject/")
	fmt.Println("  ├── go.mod")
	fmt.Println("  ├── main.go          (package main)")
	fmt.Println("  ├── mathutil/")
	fmt.Println("  │   └── mathutil.go  (package mathutil)")
	fmt.Println("  └── internal/        (private to this module)")
	fmt.Println("      └── secret.go")

	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. Uppercase first letter = exported (visible outside package)")
	fmt.Println("2. One package per directory, package name = directory name")
	fmt.Println("3. go.mod is the single source of truth for dependencies")
	fmt.Println("4. internal/ directory restricts access to parent module only")
	fmt.Println("5. No circular imports allowed forces clean architecture")
}
