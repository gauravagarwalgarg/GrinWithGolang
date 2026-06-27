/*
What:       Unbuffered vs buffered channels, send/receive, range, closing
Level:      Beginner
Analogy:    Unbuffered channel = a relay baton handoff. Buffered = a mailbox with N slots.
C++ Angle:  Channel ≈ thread-safe queue with blocking semantics built-in.
Interview:  "What happens if you send on a closed channel?" → panic.
*/
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== 1. Unbuffered Channel (Synchronous Handshake) ===")
	// Sender blocks until receiver is ready, and vice versa.
	handshake := make(chan string) // unbuffered: capacity 0

	go func() {
		time.Sleep(100 * time.Millisecond)
		handshake <- "hello from goroutine" // blocks until main receives
	}()

	msg := <-handshake // blocks until goroutine sends
	fmt.Printf("  Received: %q\n", msg)

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 2. Buffered Channel (Mailbox with N Slots) ===")
	// Sender blocks only when buffer is full.
	mailbox := make(chan int, 3) // capacity 3

	// Can send 3 values without a receiver ready
	mailbox <- 10
	mailbox <- 20
	mailbox <- 30
	// mailbox <- 40 // would block here buffer full!

	fmt.Printf("  Buffer len=%d, cap=%d\n", len(mailbox), cap(mailbox))
	fmt.Printf("  Received: %d, %d, %d\n", <-mailbox, <-mailbox, <-mailbox)

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 3. Directional Channels (Send-only, Receive-only) ===")
	// Type safety: restrict what a function can do with a channel.
	producer := func(out chan<- int) { // can only send
		for i := 1; i <= 5; i++ {
			out <- i
		}
		close(out)
	}

	consumer := func(in <-chan int) { // can only receive
		for val := range in { // range exits when channel is closed
			fmt.Printf("  Got: %d\n", val)
		}
	}

	ch := make(chan int, 5)
	producer(ch)
	consumer(ch)

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 4. Range Over Channel ===")
	numbers := make(chan int, 5)
	go func() {
		for i := 10; i <= 50; i += 10 {
			numbers <- i
		}
		close(numbers) // MUST close or range blocks forever
	}()

	fmt.Print("  Numbers: ")
	for n := range numbers {
		fmt.Printf("%d ", n)
	}
	fmt.Println()

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 5. Detecting Close with ok Idiom ===")
	ch2 := make(chan string, 2)
	ch2 <- "first"
	ch2 <- "second"
	close(ch2)

	for {
		val, ok := <-ch2
		if !ok {
			fmt.Println("  Channel closed, done reading.")
			break
		}
		fmt.Printf("  val=%q ok=%v\n", val, ok)
	}

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 6. Channel Axioms (Memorize These) ===")
	fmt.Println("  ┌─────────────────────┬────────────────────────────┐")
	fmt.Println("  │ Operation           │ Nil Chan  │ Closed Chan    │")
	fmt.Println("  ├─────────────────────┼────────────────────────────┤")
	fmt.Println("  │ Send                │ Block ∞   │ PANIC          │")
	fmt.Println("  │ Receive             │ Block ∞   │ Zero value, ok=false │")
	fmt.Println("  │ Close               │ PANIC     │ PANIC          │")
	fmt.Println("  └─────────────────────┴────────────────────────────┘")
}
