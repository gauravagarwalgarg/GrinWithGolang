/*
What this teaches:
    Arrays (fixed-size) vs slices (dynamic, growable). The slice header:
    pointer, length, capacity. append(), copy(), slicing syntax, and how
    capacity grows on reallocation.

Beginner analogy:
    "Array = fixed park bench (5 seats, no more). Slice = elastic waistband
     (stretches as you eat more data). Under the hood, a slice is just a
     window into an array."

C++ comparison:
    "Slice ≈ a fat pointer (like std::span but growable). append may
     reallocate like vector::push_back. Capacity doubles (approximately)
     on growth, just like std::vector."

Interview relevance:
    Slice internals (3-word header), append semantics (when does it copy?),
    slice aliasing bugs, and nil slice vs empty slice are top questions.
*/

package main

import "fmt"

func main() {
	// 1. Arrays fixed size, value type
	fmt.Println("--- Arrays (Fixed Size) ---")
	var arr [5]int // zero-valued: [0 0 0 0 0]
	arr[0] = 10
	arr[4] = 50
	fmt.Printf("arr = %v, len=%d\n", arr, len(arr))

	arr2 := [3]string{"Go", "is", "fun"}
	fmt.Printf("arr2 = %v\n", arr2)
	// Arrays are VALUE types assignment copies!
	arr3 := arr
	arr3[0] = 999
	fmt.Printf("arr=%v arr3=%v (independent copies)\n", arr, arr3)

	// 2. Slices dynamic, reference to underlying array
	fmt.Println("\n--- Slices (Dynamic) ---")
	s := []int{1, 2, 3, 4, 5} // slice literal (no size in brackets)
	fmt.Printf("s = %v, len=%d, cap=%d\n", s, len(s), cap(s))

	// 3. The 3-word slice header: pointer, length, capacity
	fmt.Println("\n--- Slice Header: ptr | len | cap ---")
	base := make([]int, 3, 8) // len=3, cap=8
	fmt.Printf("make([]int, 3, 8) → len=%d cap=%d %v\n", len(base), cap(base), base)

	// 4. Slicing creates a new header, shares the array
	fmt.Println("\n--- Slicing (Shared Backing Array) ---")
	original := []int{10, 20, 30, 40, 50}
	slice := original[1:4] // elements at index 1, 2, 3
	fmt.Printf("original = %v\n", original)
	fmt.Printf("slice = original[1:4] = %v, len=%d, cap=%d\n",
		slice, len(slice), cap(slice))
	slice[0] = 999 // modifies original[1]!
	fmt.Printf("After slice[0]=999: original = %v (shared!)\n", original)

	// 5. append may reallocate if capacity exceeded
	fmt.Println("\n--- append() ---")
	a := make([]int, 0, 4)
	fmt.Printf("Initial: len=%d cap=%d\n", len(a), cap(a))
	for i := 1; i <= 6; i++ {
		a = append(a, i*10)
		fmt.Printf("  append(%d): len=%d cap=%d %v\n", i*10, len(a), cap(a), a)
	}
	fmt.Println("Notice: capacity doubled when exceeded (4 → 8)")

	// 6. append multiple and spread
	fmt.Println("\n--- append multiple / spread ---")
	x := []int{1, 2, 3}
	y := []int{4, 5, 6}
	x = append(x, y...) // spread y into append
	fmt.Printf("x = %v\n", x)

	// 7. copy independent duplication
	fmt.Println("\n--- copy() ---")
	src := []int{10, 20, 30, 40}
	dst := make([]int, len(src))
	n := copy(dst, src)
	dst[0] = 999
	fmt.Printf("Copied %d elements. src=%v dst=%v (independent)\n", n, src, dst)

	// 8. nil slice vs empty slice
	fmt.Println("\n--- nil vs empty slice ---")
	var nilSlice []int
	emptySlice := []int{}
	fmt.Printf("nil slice: %v, len=%d, cap=%d, ==nil? %t\n",
		nilSlice, len(nilSlice), cap(nilSlice), nilSlice == nil)
	fmt.Printf("empty slice: %v, len=%d, cap=%d, ==nil? %t\n",
		emptySlice, len(emptySlice), cap(emptySlice), emptySlice == nil)
	fmt.Println("Both work with append prefer nil slice as zero value")

	// 9. Deleting from a slice (order-preserving)
	fmt.Println("\n--- Delete element (keep order) ---")
	items := []string{"a", "b", "c", "d", "e"}
	i := 2 // remove "c"
	items = append(items[:i], items[i+1:]...)
	fmt.Printf("After removing index 2: %v\n", items)

	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. Arrays are fixed + value type; slices are dynamic + reference-like")
	fmt.Println("2. Slice = {pointer, length, capacity} 3 machine words")
	fmt.Println("3. append returns a new slice always reassign: s = append(s, x)")
	fmt.Println("4. Sub-slices SHARE the backing array mutations are visible!")
	fmt.Println("5. Use copy() for truly independent duplicates")
}
