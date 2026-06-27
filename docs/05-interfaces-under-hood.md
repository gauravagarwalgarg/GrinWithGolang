# Interfaces Under the Hood

## Interface Representation

Go interfaces are represented as a two-word struct:

```
┌─────────────────────────────────────┐
│  Non-empty interface (iface)        │
│  ┌─────────┬─────────┐             │
│  │  *itab  │  *data  │             │
│  └────┬────┴────┬────┘             │
│       │         └── pointer to concrete value
│       └── type info + method table
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│  Empty interface (eface)            │
│  ┌─────────┬─────────┐             │
│  │  *_type │  *data  │             │
│  └────┬────┴────┬────┘             │
│       │         └── pointer to concrete value
│       └── type descriptor only (no methods)
└─────────────────────────────────────┘
```

## itab Structure

```
type itab struct {
    inter  *interfacetype  // static interface type
    _type  *_type          // concrete type stored
    hash   uint32          // copy of _type.hash (fast type switch)
    fun    [1]uintptr      // method table (variable size)
}
```

- `fun[]` contains pointers to the concrete type's methods matching the interface
- itabs are cached globally: once computed, reused for same (interface, type) pair
- Empty interface (`interface{}`) skips itab entirely just stores `*_type`

## Type Assertion Cost

| Operation | Cost | Mechanism |
|-----------|------|-----------|
| `i.(ConcreteType)` | O(1) | Compare itab._type pointer |
| `i.(Interface)` | O(1) amortized | itab cache lookup |
| Type switch | O(1) per case | Hash comparison |
| `reflect.TypeOf()` | O(1) | Read _type pointer |

Type assertions are cheap they compare pointers, not method sets.

## Method Set Rules

| Receiver | Value method set | Pointer method set |
|----------|-----------------|-------------------|
| `func (t T) Foo()` | ✓ has Foo | ✓ has Foo |
| `func (t *T) Bar()` | ✗ no Bar | ✓ has Bar |

**Why?** A value might not be addressable (e.g., map values, function returns).
You can't take `&mapValue["key"]`. So value types can't call pointer receivers.

```go
var w io.Writer
w = os.Stdout   // *os.File has Write (pointer receiver) ✓
w = myBuffer    // Buffer has Write (value receiver) ✓
// w = MyType{} // if Write is on *MyType → compile error
```

## Interface Satisfaction: Compile-Time vs Runtime

- **Compile time**: Assigning concrete type to interface variable
- **Runtime**: Type assertions, type switches, reflect
- No explicit `implements` keyword satisfaction is structural (duck typing)
- Verified at assignment, not declaration

## Performance Implications

1. **Small values (≤ pointer size)**: Stored inline in data field no allocation
2. **Larger values**: Heap-allocated, data field points to heap copy
3. **Method calls through interface**: One extra indirection (itab.fun[n])
4. **Escape to interface**: Causes heap allocation for value types > 1 word

Avoid hot-path interface{} for value types if allocation-sensitive.
