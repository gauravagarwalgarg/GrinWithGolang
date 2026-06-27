/*
What this teaches:
    Your first Go program the absolute minimum to get code running.
    Covers: package declaration, import statement, func main, and the fmt package.

Beginner analogy:
    "package = your surname (every file in a folder shares the same family name),
     main = the front door (the entry point the OS knocks on to start your program)."

C++ comparison:
    "No header files, no linking step, one binary. `go build` compiles everything
     into a single statically-linked executable. No Makefile, no CMakeLists.txt."

Interview relevance:
    Interviewers check if you know that every Go executable must have package main
    and func main(). They also ask about the difference between package main and
    library packages, and how Go tooling discovers entry points.
*/

package main

import "fmt" // fmt = "format" Go's standard I/O formatting package

func main() {
	// 1. The classic every language starts here
	fmt.Println("Hello, World!")

	// 2. Package declaration
	// Every .go file begins with `package <name>`.
	// Files in the same directory must share the same package name.
	fmt.Println("\n--- Package Declaration ---")
	fmt.Println("This file belongs to 'package main'.")
	fmt.Println("'main' is special it tells Go: build an executable, not a library.")

	// 3. Import statement
	// Imports bring other packages into scope. Unused imports = compile error.
	fmt.Println("\n--- Import Statement ---")
	fmt.Println("We imported 'fmt' for formatted I/O.")
	fmt.Println("Go refuses to compile if you import something and never use it.")
	fmt.Println("This keeps binaries lean and dependencies honest.")

	// 4. func main the entry point
	fmt.Println("\n--- func main ---")
	fmt.Println("func main() takes no arguments and returns nothing.")
	fmt.Println("To pass args, use os.Args or the flag package.")
	fmt.Println("To signal failure, call os.Exit(code).")

	// 5. fmt package basics
	fmt.Println("\n--- fmt Package ---")
	name := "Gopher"
	age := 15 // Go was released in 2009

	fmt.Println("Println adds a newline automatically:", name)
	fmt.Printf("Printf uses verbs: %s is %d years old\n", name, age)
	fmt.Printf("%%v (value): %v | %%T (type): %T\n", age, age)

	// 6. Multiple values in Println
	fmt.Println("\n--- Multiple Values ---")
	fmt.Println("Go", "was", "designed", "at", "Google", "in", 2007)

	// 7. Key takeaways
	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. package main + func main() = executable")
	fmt.Println("2. Unused imports are compile errors (not warnings)")
	fmt.Println("3. fmt.Println is your go-to for quick output")
	fmt.Println("4. go run main.go compile and run in one step")
	fmt.Println("5. go build produces a static binary, ship it anywhere")
}
