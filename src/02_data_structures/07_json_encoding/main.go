/*
What this teaches:
    Struct tags for JSON, json.Marshal/Unmarshal, custom MarshalJSON methods,
    handling optional fields with pointers and omitempty, and common patterns
    for API data interchange.

Beginner analogy:
    "Tags are instructions on how to pack a suitcase they tell the JSON
     encoder which pocket to put each field in, and whether to skip empty ones."

C++ comparison:
    "Reflection-based serialization via struct tags. No external codegen needed
     for basic cases (unlike protobuf). Runtime reflection reads tags similar
     to Java annotations but at the field level."

Interview relevance:
    struct tag syntax, omitempty semantics, difference between nil pointer and
    zero value for optional fields, custom marshaling, and handling unknown fields.
*/

package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// 1. Basic struct with JSON tags
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Age       int    `json:"age,omitempty"`     // omit if zero value
	Phone     string `json:"phone,omitempty"`   // omit if empty string
	internal  string // unexported NEVER serialized
}

// 2. Struct with optional fields (pointer = distinguishes null from zero)
type Config struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Debug   *bool  `json:"debug,omitempty"`   // nil = absent, *false = explicitly false
	Timeout *int   `json:"timeout,omitempty"` // nil = absent, *0 = explicitly zero
}

// 3. Nested structs
type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
	Zip    string `json:"zip"`
}

type Employee struct {
	Name    string  `json:"name"`
	Role    string  `json:"role"`
	Address Address `json:"address"`
}

// 4. Custom MarshalJSON control the output format
type Event struct {
	Name string
	Date time.Time
}

func (e Event) MarshalJSON() ([]byte, error) {
	// Custom format: encode date as "YYYY-MM-DD" string
	type Alias struct {
		Name string `json:"name"`
		Date string `json:"date"`
	}
	return json.Marshal(Alias{
		Name: e.Name,
		Date: e.Date.Format("2006-01-02"),
	})
}

func (e *Event) UnmarshalJSON(data []byte) error {
	type Alias struct {
		Name string `json:"name"`
		Date string `json:"date"`
	}
	var aux Alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	e.Name = aux.Name
	t, err := time.Parse("2006-01-02", aux.Date)
	if err != nil {
		return err
	}
	e.Date = t
	return nil
}

// 5. Using json.RawMessage for deferred parsing
type APIResponse struct {
	Status string          `json:"status"`
	Data   json.RawMessage `json:"data"` // parse later based on Status
}

func main() {
	// Marshal: struct → JSON
	fmt.Println("--- Marshal (struct → JSON) ---")
	user := User{
		ID:        1,
		FirstName: "Go",
		LastName:  "Pher",
		Email:     "gopher@go.dev",
		Age:       0,     // omitempty: will be omitted
		Phone:     "",    // omitempty: will be omitted
		internal:  "secret",
	}
	data, _ := json.MarshalIndent(user, "", "  ")
	fmt.Println(string(data))

	// Unmarshal: JSON → struct
	fmt.Println("\n--- Unmarshal (JSON → struct) ---")
	jsonStr := `{"id":2,"first_name":"Jane","last_name":"Doe","email":"jane@go.dev","age":28}`
	var parsed User
	_ = json.Unmarshal([]byte(jsonStr), &parsed)
	fmt.Printf("Parsed: %+v\n", parsed)

	// Optional fields with pointers
	fmt.Println("\n--- Optional Fields (Pointer) ---")
	falseVal := false
	cfg := Config{Host: "localhost", Port: 8080, Debug: &falseVal}
	data, _ = json.MarshalIndent(cfg, "", "  ")
	fmt.Println(string(data))
	fmt.Println("Debug is *false (shown), Timeout is nil (omitted)")

	// Nested struct
	fmt.Println("\n--- Nested Structs ---")
	emp := Employee{
		Name: "Alice",
		Role: "Engineer",
		Address: Address{
			Street: "123 Main St",
			City:   "Gopher City",
			Zip:    "12345",
		},
	}
	data, _ = json.MarshalIndent(emp, "", "  ")
	fmt.Println(string(data))

	// Custom MarshalJSON
	fmt.Println("\n--- Custom MarshalJSON ---")
	event := Event{Name: "GopherCon", Date: time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC)}
	data, _ = json.Marshal(event)
	fmt.Printf("Custom format: %s\n", data)

	// Unmarshal back with custom UnmarshalJSON
	var decoded Event
	_ = json.Unmarshal(data, &decoded)
	fmt.Printf("Decoded: %s on %s\n", decoded.Name, decoded.Date.Format("Jan 2, 2006"))

	// Unknown/extra fields are silently ignored by default
	fmt.Println("\n--- Unknown Fields (Ignored) ---")
	extra := `{"id":3,"first_name":"Bob","unknown_field":"ignored","last_name":"X","email":"b@x.com"}`
	var u2 User
	_ = json.Unmarshal([]byte(extra), &u2)
	fmt.Printf("Parsed (extras ignored): %+v\n", u2)

	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. Struct tags: `json:\"name,omitempty\"` control serialization")
	fmt.Println("2. omitempty skips zero values (0, \"\", nil, empty slice/map)")
	fmt.Println("3. Use *T for optional fields to distinguish nil from zero")
	fmt.Println("4. Custom MarshalJSON/UnmarshalJSON for special formats")
	fmt.Println("5. Only exported (uppercase) fields are serialized")
}
