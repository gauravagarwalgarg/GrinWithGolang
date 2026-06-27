/*
What:       unsafe.Pointer, unsafe.Sizeof, type casting, when needed, why it breaks GC
Level:      Beginner
Analogy:    unsafe = taking off the safety helmet. Only for experts.
C++ Angle:  Like reinterpret_cast breaks type safety. Needed for FFI and performance hacks.
Interview:  "When would you use unsafe in Go?" → syscalls, mmap, zero-copy parsing, struct hacking.
*/
package main

import (
	"fmt"
	"unsafe"
)

// ─── Example struct for pointer arithmetic ───────────────────────
type Pixel struct {
	R, G, B, A uint8
}

type Header struct {
	Version uint32
	Length  uint64
	Flags   uint16
}

func main() {
	fmt.Println("=== 1. unsafe.Sizeof Type Memory Sizes ===")
	fmt.Printf("  int:          %d bytes\n", unsafe.Sizeof(int(0)))
	fmt.Printf("  int32:        %d bytes\n", unsafe.Sizeof(int32(0)))
	fmt.Printf("  int64:        %d bytes\n", unsafe.Sizeof(int64(0)))
	fmt.Printf("  string:       %d bytes (header: ptr + len)\n", unsafe.Sizeof(""))
	fmt.Printf("  []byte:       %d bytes (header: ptr + len + cap)\n", unsafe.Sizeof([]byte{}))
	fmt.Printf("  Pixel:        %d bytes\n", unsafe.Sizeof(Pixel{}))
	fmt.Printf("  Header:       %d bytes (with padding)\n", unsafe.Sizeof(Header{}))

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 2. unsafe.Pointer Type Conversion ===")
	// Convert between unrelated pointer types (like reinterpret_cast)
	var pixel Pixel = Pixel{R: 255, G: 128, B: 64, A: 255}
	fmt.Printf("  Pixel: %+v\n", pixel)

	// View the 4-byte Pixel as a single uint32 (RGBA packed)
	packed := *(*uint32)(unsafe.Pointer(&pixel))
	fmt.Printf("  As uint32: 0x%08X\n", packed)

	// Convert back
	var pixel2 Pixel
	*(*uint32)(unsafe.Pointer(&pixel2)) = packed
	fmt.Printf("  Back to Pixel: %+v\n", pixel2)

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 3. Pointer Arithmetic (Reading Struct Fields) ===")
	h := Header{Version: 1, Length: 1024, Flags: 0x00FF}
	fmt.Printf("  Header: %+v\n", h)

	// Access Length field via pointer arithmetic
	// offset of Length = unsafe.Offsetof(h.Length)
	base := unsafe.Pointer(&h)
	lengthPtr := (*uint64)(unsafe.Add(base, unsafe.Offsetof(h.Length)))
	fmt.Printf("  Length via pointer arithmetic: %d\n", *lengthPtr)

	// Modify Flags via pointer arithmetic
	flagsPtr := (*uint16)(unsafe.Add(base, unsafe.Offsetof(h.Flags)))
	*flagsPtr = 0x1234
	fmt.Printf("  Modified Flags: 0x%04X\n", h.Flags)

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 4. String ↔ []byte Zero-Copy (Dangerous!) ===")
	// Normal conversion copies data: []byte(str) allocates
	// This hack avoids the copy but is UNSAFE don't mutate!
	str := "hello unsafe world"
	// String header: {ptr, len}. Slice header: {ptr, len, cap}
	// We can read the string's underlying bytes without copying
	strBytes := unsafe.Slice(unsafe.StringData(str), len(str))
	fmt.Printf("  Zero-copy bytes: %v\n", strBytes)
	fmt.Printf("  As string: %s\n", string(strBytes))

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 5. Why unsafe Breaks GC Guarantees ===")
	fmt.Println("  The GC tracks pointers to know what's alive.")
	fmt.Println("  unsafe.Pointer can:")
	fmt.Println("    • Create pointers the GC can't see (memory leaks)")
	fmt.Println("    • Point into the middle of objects (invalid after GC moves them)")
	fmt.Println("    • Alias types incorrectly (data corruption)")
	fmt.Println()
	fmt.Println("  Rules for safe usage:")
	fmt.Println("    1. Never store unsafe.Pointer in a uintptr variable across statements")
	fmt.Println("       (GC may move the object between statements)")
	fmt.Println("    2. Conversion must happen in a single expression:")
	fmt.Println("       OK:  p := (*T)(unsafe.Pointer(&x))")
	fmt.Println("       BAD: u := uintptr(unsafe.Pointer(&x)); p := (*T)(unsafe.Pointer(u))")
	fmt.Println("    3. Use go vet to detect violations")

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 6. Legitimate Use Cases ===")
	fmt.Println("  • syscall/mmap: raw memory from OS")
	fmt.Println("  • Binary protocol parsing: zero-copy header decode")
	fmt.Println("  • reflect: runtime type introspection internals")
	fmt.Println("  • sync/atomic: generic atomic operations")
	fmt.Println("  • cgo/FFI: passing pointers to C code")
	fmt.Println("  • Performance: avoid allocation in hot paths")
	fmt.Println()
	fmt.Println("  In 99% of Go code, you should NEVER need unsafe.")
	fmt.Println("  If you think you need it, you probably don't.")
}
