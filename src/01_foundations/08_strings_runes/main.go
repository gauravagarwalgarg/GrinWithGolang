/*
What this teaches:
    Strings (immutable byte slices), runes (Unicode code points), iterating
    strings with range (yields runes), the strings package, and strings.Builder
    for efficient concatenation.

Beginner analogy:
    "A string is a necklace of beads; a rune is one bead (which might be wide 
     like an emoji taking multiple bytes). len() counts the thread length in
     bytes, not the number of beads."

C++ comparison:
    "Go strings are UTF-8 by default. len() gives bytes, not characters.
     Immutable like std::string_view but owned. strings.Builder ≈ std::ostringstream."

Interview relevance:
    Commonly asked: Why does len('こんにちは') return 15, not 5? How to count
    runes? How does range differ from indexing? Why use Builder over += ?
*/

package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func main() {
	// 1. Strings are immutable byte slices
	fmt.Println("--- Strings as Byte Slices ---")
	s := "Hello, 世界!" // mix of ASCII and multi-byte characters
	fmt.Printf("String: %s\n", s)
	fmt.Printf("len(s) = %d bytes (NOT characters)\n", len(s))
	fmt.Printf("utf8.RuneCountInString(s) = %d runes\n", utf8.RuneCountInString(s))

	// 2. Indexing gives bytes, not characters
	fmt.Println("\n--- Byte Indexing ---")
	ascii := "Go"
	fmt.Printf("ascii[0] = %d (%c) one byte per ASCII char\n", ascii[0], ascii[0])
	multi := "世界"
	fmt.Printf("multi[0] = %d (0x%x) first BYTE of '世', not the rune!\n",
		multi[0], multi[0])

	// 3. Runes Unicode code points (int32)
	fmt.Println("\n--- Runes ---")
	r := '🚀'
	fmt.Printf("Rune '🚀': value=%d (U+%04X), type=%T\n", r, r, r)
	fmt.Printf("Size in UTF-8: %d bytes\n", utf8.RuneLen(r))

	// 4. Range iteration yields runes (not bytes!)
	fmt.Println("\n--- Range Iteration (Rune by Rune) ---")
	word := "Gö🎉"
	for i, ch := range word {
		fmt.Printf("  byte_offset=%d rune=%c (U+%04X)\n", i, ch, ch)
	}

	// 5. Converting between strings, bytes, and runes
	fmt.Println("\n--- Conversions ---")
	original := "café"
	byteSlice := []byte(original)
	runeSlice := []rune(original)
	fmt.Printf("[]byte: %v (len=%d)\n", byteSlice, len(byteSlice))
	fmt.Printf("[]rune: %v (len=%d)\n", runeSlice, len(runeSlice))
	fmt.Printf("Back to string: %s\n", string(runeSlice))

	// 6. strings package common operations
	fmt.Println("\n--- strings Package ---")
	text := "  Go is Simple and Powerful  "
	fmt.Printf("TrimSpace: %q\n", strings.TrimSpace(text))
	fmt.Printf("ToUpper:   %q\n", strings.ToUpper("hello"))
	fmt.Printf("Contains:  %t\n", strings.Contains("seafood", "foo"))
	fmt.Printf("Split:     %v\n", strings.Split("a,b,c", ","))
	fmt.Printf("Join:      %s\n", strings.Join([]string{"Go", "is", "fun"}, " "))
	fmt.Printf("Replace:   %s\n", strings.ReplaceAll("foo.bar.baz", ".", "/"))
	fmt.Printf("HasPrefix: %t\n", strings.HasPrefix("Golang", "Go"))
	fmt.Printf("Index:     %d\n", strings.Index("hello", "ll"))

	// 7. strings.Builder efficient concatenation (avoids O(n²) copies)
	fmt.Println("\n--- strings.Builder ---")
	var b strings.Builder
	for i := 0; i < 5; i++ {
		fmt.Fprintf(&b, "item_%d ", i)
	}
	result := b.String()
	fmt.Printf("Built: %q\n", result)
	fmt.Printf("Builder avoids reallocating on every += operation\n")

	// 8. Multiline strings with backticks (raw strings)
	fmt.Println("\n--- Raw Strings (Backticks) ---")
	raw := `This is a raw string.
No escape sequences: \n \t are literal.
Great for regex, SQL, JSON templates.`
	fmt.Println(raw)

	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. len() = byte count. Use utf8.RuneCountInString for char count")
	fmt.Println("2. range on strings iterates runes, not bytes")
	fmt.Println("3. Strings are immutable modification creates a new string")
	fmt.Println("4. Use strings.Builder for loops; += is O(n²)")
	fmt.Println("5. Backtick strings are raw no escape processing")
}
