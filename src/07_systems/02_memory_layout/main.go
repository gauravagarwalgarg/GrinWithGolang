/*
What:       Struct padding, alignment, unsafe.Sizeof/Alignof/Offsetof, field reordering
Level:      Beginner
Analogy:    Padding = empty seats between passengers to fit rows evenly.
C++ Angle:  Same alignment rules as C/C++ but Go inserts padding automatically. Reorder large→small.
Interview:  "How do you optimize struct memory in Go?" → Reorder fields large to small.
*/
package main

import (
	"fmt"
	"unsafe"
)

// ─── Poorly ordered struct (lots of padding) ─────────────────────
type BadLayout struct {
	a bool    // 1 byte + 7 padding (next field needs 8-byte alignment)
	b float64 // 8 bytes
	c int32   // 4 bytes + 4 padding (struct must be multiple of largest align)
	d bool    // 1 byte + 3 padding
	e int16   // 2 bytes + 2 padding (total struct alignment padding)
}

// ─── Well-ordered struct (minimal padding) ───────────────────────
type GoodLayout struct {
	b float64 // 8 bytes (largest first)
	c int32   // 4 bytes
	e int16   // 2 bytes
	a bool    // 1 byte
	d bool    // 1 byte packs next to 'a'
} // total: 16 bytes (8 + 4 + 2 + 1 + 1 = 16, aligned to 8)

// ─── Real-world example: HTTP header struct ──────────────────────
type BadHTTPHeader struct {
	IsSecure    bool   // 1 byte + 7 padding
	ContentLen  int64  // 8 bytes
	StatusCode  int16  // 2 bytes + 6 padding
	KeepAlive   bool   // 1 byte + 7 padding
	MaxBodySize int64  // 8 bytes
}

type GoodHTTPHeader struct {
	ContentLen  int64 // 8 bytes
	MaxBodySize int64 // 8 bytes
	StatusCode  int16 // 2 bytes
	IsSecure    bool  // 1 byte
	KeepAlive   bool  // 1 byte + 4 padding to align struct
}

func inspectStruct(name string, size, align uintptr) {
	fmt.Printf("  %-20s Size: %2d bytes, Align: %d bytes\n", name, size, align)
}

func main() {
	fmt.Println("=== 1. Struct Sizes Bad vs Good Layout ===")
	inspectStruct("BadLayout", unsafe.Sizeof(BadLayout{}), unsafe.Alignof(BadLayout{}))
	inspectStruct("GoodLayout", unsafe.Sizeof(GoodLayout{}), unsafe.Alignof(GoodLayout{}))
	fmt.Printf("  Savings: %d bytes per struct!\n",
		unsafe.Sizeof(BadLayout{})-unsafe.Sizeof(GoodLayout{}))

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 2. Field Offsets Where Padding Lives ===")
	fmt.Println("  BadLayout field offsets:")
	var bad BadLayout
	fmt.Printf("    a (bool):    offset=%d, size=%d\n", unsafe.Offsetof(bad.a), unsafe.Sizeof(bad.a))
	fmt.Printf("    b (float64): offset=%d, size=%d\n", unsafe.Offsetof(bad.b), unsafe.Sizeof(bad.b))
	fmt.Printf("    c (int32):   offset=%d, size=%d\n", unsafe.Offsetof(bad.c), unsafe.Sizeof(bad.c))
	fmt.Printf("    d (bool):    offset=%d, size=%d\n", unsafe.Offsetof(bad.d), unsafe.Sizeof(bad.d))
	fmt.Printf("    e (int16):   offset=%d, size=%d\n", unsafe.Offsetof(bad.e), unsafe.Sizeof(bad.e))

	fmt.Println("\n  GoodLayout field offsets:")
	var good GoodLayout
	fmt.Printf("    b (float64): offset=%d, size=%d\n", unsafe.Offsetof(good.b), unsafe.Sizeof(good.b))
	fmt.Printf("    c (int32):   offset=%d, size=%d\n", unsafe.Offsetof(good.c), unsafe.Sizeof(good.c))
	fmt.Printf("    e (int16):   offset=%d, size=%d\n", unsafe.Offsetof(good.e), unsafe.Sizeof(good.e))
	fmt.Printf("    a (bool):    offset=%d, size=%d\n", unsafe.Offsetof(good.a), unsafe.Sizeof(good.a))
	fmt.Printf("    d (bool):    offset=%d, size=%d\n", unsafe.Offsetof(good.d), unsafe.Sizeof(good.d))

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 3. Real-World Example: HTTP Header ===")
	inspectStruct("BadHTTPHeader", unsafe.Sizeof(BadHTTPHeader{}), unsafe.Alignof(BadHTTPHeader{}))
	inspectStruct("GoodHTTPHeader", unsafe.Sizeof(GoodHTTPHeader{}), unsafe.Alignof(GoodHTTPHeader{}))

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 4. Primitive Type Sizes and Alignments ===")
	fmt.Println("  ┌────────────┬──────┬───────┐")
	fmt.Println("  │ Type       │ Size │ Align │")
	fmt.Println("  ├────────────┼──────┼───────┤")
	fmt.Printf("  │ bool       │ %4d │ %5d │\n", unsafe.Sizeof(true), unsafe.Alignof(true))
	fmt.Printf("  │ int8       │ %4d │ %5d │\n", unsafe.Sizeof(int8(0)), unsafe.Alignof(int8(0)))
	fmt.Printf("  │ int16      │ %4d │ %5d │\n", unsafe.Sizeof(int16(0)), unsafe.Alignof(int16(0)))
	fmt.Printf("  │ int32      │ %4d │ %5d │\n", unsafe.Sizeof(int32(0)), unsafe.Alignof(int32(0)))
	fmt.Printf("  │ int64      │ %4d │ %5d │\n", unsafe.Sizeof(int64(0)), unsafe.Alignof(int64(0)))
	fmt.Printf("  │ float64    │ %4d │ %5d │\n", unsafe.Sizeof(float64(0)), unsafe.Alignof(float64(0)))
	fmt.Printf("  │ string     │ %4d │ %5d │\n", unsafe.Sizeof(""), unsafe.Alignof(""))
	fmt.Printf("  │ []byte     │ %4d │ %5d │\n", unsafe.Sizeof([]byte{}), unsafe.Alignof([]byte{}))
	fmt.Println("  └────────────┴──────┴───────┘")

	// ─────────────────────────────────────────────
	fmt.Println("\n=== Optimization Rule ===")
	fmt.Println("  Order struct fields from LARGEST to SMALLEST alignment:")
	fmt.Println("    float64/int64/pointer → int32 → int16 → bool/int8")
	fmt.Println("  Tool: go vet -structlayout (or fieldalignment linter)")
}
