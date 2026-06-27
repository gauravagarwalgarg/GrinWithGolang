/*
What this teaches:
    Control flow in Go if/else, for (Go's ONLY loop keyword), switch (no
    fallthrough by default), and defer (cleanup scheduling).

Beginner analogy:
    "defer = 'I'll clean up when I leave the room' you schedule cleanup at
     the top, and Go guarantees it runs when the function exits, even on panic."

C++ comparison:
    "defer is like RAII destructors but explicit and readable at the call site.
     No while/do-while for covers all loop patterns. Switch doesn't fall
     through by default (no forgotten break bugs)."

Interview relevance:
    defer ordering (LIFO stack), defer in loops, the init statement in if,
    and switch with no condition (clean if-else chains) are common topics.
*/

package main

import "fmt"

func main() {
	// 1. if/else with init statement
	fmt.Println("--- if/else ---")
	score := 85
	if score >= 90 {
		fmt.Println("Grade: A")
	} else if score >= 80 {
		fmt.Println("Grade: B")
	} else {
		fmt.Println("Grade: C or below")
	}

	// if with init statement variable scoped to if/else block
	fmt.Println("\n--- if with init ---")
	if x := 42 % 2; x == 0 {
		fmt.Println("42 is even")
	} else {
		fmt.Println("42 is odd")
	}
	// x is not accessible here scoped to the if block

	// 2. for the ONLY loop in Go (replaces while, do-while, for)
	fmt.Println("\n--- for (traditional) ---")
	for i := 0; i < 5; i++ {
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// for as "while"
	fmt.Println("\n--- for as while ---")
	n := 1
	for n < 32 {
		fmt.Printf("%d ", n)
		n *= 2
	}
	fmt.Println()

	// for range iterating over a slice
	fmt.Println("\n--- for range ---")
	fruits := []string{"apple", "banana", "cherry"}
	for i, fruit := range fruits {
		fmt.Printf("  [%d] %s\n", i, fruit)
	}

	// infinite loop with break
	fmt.Println("\n--- infinite loop + break ---")
	count := 0
	for {
		count++
		if count == 3 {
			break
		}
	}
	fmt.Printf("Broke out at count = %d\n", count)

	// 3. switch clean, no fallthrough by default
	fmt.Println("\n--- switch ---")
	day := "Wednesday"
	switch day {
	case "Monday", "Tuesday", "Wednesday", "Thursday", "Friday":
		fmt.Println(day, "is a weekday")
	case "Saturday", "Sunday":
		fmt.Println(day, "is a weekend")
	default:
		fmt.Println("Unknown day")
	}

	// switch with no condition (replaces if-else chains)
	fmt.Println("\n--- switch (no condition) ---")
	temp := 35
	switch {
	case temp > 40:
		fmt.Println("Extreme heat")
	case temp > 30:
		fmt.Println("Hot")
	case temp > 20:
		fmt.Println("Warm")
	default:
		fmt.Println("Cool or cold")
	}

	// 4. defer LIFO execution, runs when function exits
	fmt.Println("\n--- defer ---")
	fmt.Println("Start")
	defer fmt.Println("Deferred 1 (runs last)")
	defer fmt.Println("Deferred 2 (runs second-to-last)")
	fmt.Println("End of main body")

	// Practical defer: simulating resource cleanup
	fmt.Println("\n--- defer for cleanup ---")
	processFile()

	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. 'for' is Go's only loop it replaces while and do-while")
	fmt.Println("2. switch doesn't fall through use 'fallthrough' keyword if needed")
	fmt.Println("3. defer runs in LIFO order when the function returns")
	fmt.Println("4. if/switch can have init statements for tighter scoping")
}

func processFile() {
	fmt.Println("  Opening file...")
	defer fmt.Println("  Closing file... (deferred)")
	fmt.Println("  Reading data...")
	fmt.Println("  Processing data...")
	// Even if a panic happened here, defer would still run
}
