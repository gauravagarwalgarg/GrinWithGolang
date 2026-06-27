/*
What this teaches:
    Maps in Go creation, CRUD operations, the comma-ok idiom for existence
    checks, random iteration order, delete(), and practical patterns.

Beginner analogy:
    "A map is a phone book name → number. You look up a name and instantly
     get the number. But the pages aren't in any guaranteed order!"

C++ comparison:
    "Like unordered_map. Hash buckets with 8 key-value pairs per bucket.
     Grows via evacuation (incremental rehash). Not thread-safe for concurrent
     writes use sync.Map or a mutex."

Interview relevance:
    Comma-ok idiom, zero value for missing keys, iteration order randomization,
    maps not being safe for concurrent writes, and using maps as sets.
*/

package main

import (
	"fmt"
	"sort"
)

func main() {
	// 1. Creating maps
	fmt.Println("--- Creating Maps ---")
	// Literal syntax
	ages := map[string]int{
		"Alice": 30,
		"Bob":   25,
		"Carol": 35,
	}
	fmt.Printf("ages = %v\n", ages)

	// make() for empty map with size hint
	scores := make(map[string]float64, 10) // hint: expect ~10 entries
	scores["math"] = 95.5
	scores["science"] = 88.0
	fmt.Printf("scores = %v\n", scores)

	// 2. CRUD operations
	fmt.Println("\n--- CRUD Operations ---")
	// Create / Update
	ages["Dave"] = 28
	fmt.Printf("After adding Dave: %v\n", ages)

	// Read
	aliceAge := ages["Alice"]
	fmt.Printf("Alice's age: %d\n", aliceAge)

	// Read missing key → zero value (not an error!)
	missingAge := ages["Zara"]
	fmt.Printf("Zara's age (missing): %d ← zero value!\n", missingAge)

	// Delete
	delete(ages, "Bob")
	fmt.Printf("After deleting Bob: %v\n", ages)

	// 3. Comma-ok idiom distinguish "exists with zero" from "missing"
	fmt.Println("\n--- Comma-Ok Idiom ---")
	inventory := map[string]int{
		"apples":  0, // explicitly zero
		"bananas": 5,
	}
	if count, ok := inventory["apples"]; ok {
		fmt.Printf("apples exist, count = %d\n", count)
	}
	if _, ok := inventory["oranges"]; !ok {
		fmt.Println("oranges not in inventory")
	}

	// 4. Iteration order is RANDOM by design
	fmt.Println("\n--- Iteration (Random Order) ---")
	colors := map[string]string{
		"red":   "#FF0000",
		"green": "#00FF00",
		"blue":  "#0000FF",
		"white": "#FFFFFF",
	}
	fmt.Println("Iterating (order may differ each run):")
	for name, hex := range colors {
		fmt.Printf("  %s → %s\n", name, hex)
	}

	// 5. Sorted iteration collect keys, sort, then iterate
	fmt.Println("\n--- Sorted Iteration ---")
	keys := make([]string, 0, len(colors))
	for k := range colors {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("  %s → %s\n", k, colors[k])
	}

	// 6. Maps as sets
	fmt.Println("\n--- Map as Set ---")
	seen := make(map[string]bool)
	words := []string{"go", "is", "go", "fun", "is", "go"}
	for _, w := range words {
		seen[w] = true
	}
	fmt.Printf("Unique words: ")
	for word := range seen {
		fmt.Printf("%s ", word)
	}
	fmt.Println()

	// 7. Map of slices grouping pattern
	fmt.Println("\n--- Map of Slices (Grouping) ---")
	students := map[string][]string{
		"math":    {"Alice", "Bob"},
		"science": {"Carol"},
	}
	students["math"] = append(students["math"], "Dave")
	students["art"] = append(students["art"], "Eve") // auto-creates slice
	for class, names := range students {
		fmt.Printf("  %s: %v\n", class, names)
	}

	// 8. Nil map vs empty map
	fmt.Println("\n--- Nil vs Empty Map ---")
	var nilMap map[string]int
	emptyMap := map[string]int{}
	fmt.Printf("nil map: %v, ==nil? %t\n", nilMap, nilMap == nil)
	fmt.Printf("empty map: %v, ==nil? %t\n", emptyMap, emptyMap == nil)
	fmt.Println("Reading nil map is safe (returns zero). Writing panics!")

	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. Always use comma-ok to distinguish zero from absent")
	fmt.Println("2. Iteration order is randomized don't depend on it")
	fmt.Println("3. Maps are NOT safe for concurrent read/write")
	fmt.Println("4. A nil map panics on write always initialize before storing")
	fmt.Println("5. delete() on missing key is a no-op (no panic)")
}
