/*
What this teaches:
    Composition over inheritance in Go embedding structs to compose behavior
    instead of building class hierarchies. Shows Reader/Writer composition into
    ReadWriter, and how embedding a Logger into a Service provides delegation.

Beginner analogy:
    "Inheritance says 'I am a kind of X.' Embedding says 'I have an X inside me
     and I can do everything X can do.' It's the difference between being born
     royalty vs hiring a royal advisor."

C++ comparison:
    "No vtable chain. Embedding is delegation, not IS-A. In C++ you'd use private
     inheritance or member objects to achieve similar composition, but Go makes it
     first-class by promoting the embedded type's methods to the outer type."

Interview relevance:
    Interviewers ask why Go has no class keyword, how embedding differs from
    inheritance, and when promoted methods get shadowed. Understanding composition
    is essential for idiomatic Go design.
*/

package main

import (
	"fmt"
	"strings"
	"time"
)

// --- Interfaces: small, composable contracts ---

type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

// Composed interface Reader + Writer = ReadWriter
type ReadWriter interface {
	Reader
	Writer
}

// --- Logger: a reusable component ---

type Logger struct {
	Prefix string
}

func (l *Logger) Log(msg string) {
	fmt.Printf("[%s] %s: %s\n", time.Now().Format("15:04:05"), l.Prefix, msg)
}

func (l *Logger) Logf(format string, args ...any) {
	l.Log(fmt.Sprintf(format, args...))
}

// --- Service embeds Logger composition, not inheritance ---

type Service struct {
	Logger // Embedded: Service "has-a" Logger, gains Log/Logf methods
	Name   string
}

func (s *Service) Start() {
	s.Log("starting service: " + s.Name) // Promoted method from Logger
}

func (s *Service) Stop() {
	s.Log("stopping service: " + s.Name)
}

// --- Buffer implements ReadWriter via composition ---

type Buffer struct {
	data []byte
}

func (b *Buffer) Read(p []byte) (int, error) {
	n := copy(p, b.data)
	b.data = b.data[n:]
	return n, nil
}

func (b *Buffer) Write(p []byte) (int, error) {
	b.data = append(b.data, p...)
	return len(p), nil
}

// --- Demonstrating method shadowing ---

type AdvancedService struct {
	Service // Embed Service (which embeds Logger)
}

// Shadows the promoted Log method outer type wins
func (a *AdvancedService) Log(msg string) {
	fmt.Printf(">>> ADVANCED: %s\n", strings.ToUpper(msg))
}

func main() {
	fmt.Println("=== Composition Over Inheritance ===")

	// 1. Embedding Logger into Service
	fmt.Println("\n--- Embedding Logger into Service ---")
	svc := &Service{
		Logger: Logger{Prefix: "APP"},
		Name:   "UserService",
	}
	svc.Start()
	svc.Stop()
	svc.Logf("processed %d requests", 42) // Direct access to Logger's method

	// 2. ReadWriter composition
	fmt.Println("\n--- ReadWriter Composition ---")
	var rw ReadWriter = &Buffer{}
	rw.Write([]byte("Hello, composition!"))
	out := make([]byte, 20)
	n, _ := rw.Read(out)
	fmt.Printf("Read from buffer: %q\n", string(out[:n]))

	// 3. Method shadowing
	fmt.Println("\n--- Method Shadowing ---")
	adv := &AdvancedService{
		Service: Service{Logger: Logger{Prefix: "BASE"}, Name: "AdvSvc"},
	}
	adv.Log("this uses the shadowed method")     // Calls AdvancedService.Log
	adv.Service.Log("this bypasses to original") // Explicit access

	// 4. Key takeaways
	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. Go has no 'extends' keyword use embedding for composition")
	fmt.Println("2. Embedded type methods are promoted to the outer type")
	fmt.Println("3. The outer type can shadow embedded methods")
	fmt.Println("4. Embedding is NOT inheritance no polymorphic dispatch on the embedded type")
	fmt.Println("5. Prefer small interfaces + embedding over deep hierarchies")
}
