# Go vs C++: Philosophy & Key Differences

## Design Philosophy

| Aspect | C++ | Go |
|--------|-----|-----|
| Philosophy | Zero-cost abstractions, maximum control | Simplicity, fast compilation, pragmatism |
| Paradigm | Multi-paradigm (OOP, generic, functional) | Composition + interfaces + concurrency |
| Compilation | Slow (templates, headers) | Fast (seconds for large projects) |
| Memory | Manual (RAII, smart pointers) | Garbage collected (tri-color mark & sweep) |
| Error handling | Exceptions + error codes | Error values (`error` interface) |
| Inheritance | Class hierarchies, virtual dispatch | No inheritance composition only |
| Generics | Templates (Turing-complete) | Type parameters (Go 1.18+, simpler) |

## No Inheritance Composition Instead

```
C++: class Dog : public Animal { ... }           // inheritance hierarchy
Go:  type Dog struct { Animal }                   // embedding (has-a, not is-a)
```

Go uses **struct embedding** for code reuse and **interfaces** for polymorphism.
No vtables, no diamond problem, no fragile base class.

## No Constructors Factory Functions

```
C++: MyClass obj(arg1, arg2);                     // constructor
Go:  obj := NewMyClass(arg1, arg2)                // factory function (convention)
```

Zero-value initialization means most types are usable without constructors.

## Error Values vs Exceptions

```
C++: try { risky(); } catch (std::exception& e) { ... }
Go:  result, err := risky(); if err != nil { return err }
```

Go errors are explicit, checked at each call site. No hidden control flow.
`errors.Is()` and `errors.As()` replace catch hierarchies.

## Concurrency Model

| Feature | C++ | Go |
|---------|-----|-----|
| Primitive | `std::thread` | goroutine (2KB stack) |
| Communication | `std::mutex`, atomics | channels (CSP model) |
| Cost | ~1MB stack per thread | ~2KB, M:N scheduled |
| Paradigm | Shared memory + locks | "Share by communicating" |

## No Pointer Arithmetic

```
C++: int* p = arr; p += 5;    // pointer arithmetic
Go:  // Not possible. Use slices with bounds checking.
```

Go has pointers (`*T`) but no arithmetic. Slices replace pointer + length patterns.

## What Go Deliberately Omits

- No header files (import path = package)
- No preprocessor (#define, #ifdef)
- No operator overloading
- No implicit conversions
- No default parameters
- No function overloading
- No ternary operator

Each omission reduces cognitive load and eliminates a class of bugs.
The tradeoff: more verbose code, but every Go developer reads the same idioms.
