/*
What:       select for multiplexing channels, timeouts, default, done pattern
Level:      Beginner
Analogy:    select = a waiter checking multiple tables serves whoever signals first.
C++ Angle:  Like epoll/select for channels. Blocks until one case is ready.
Interview:  "How do you implement a timeout in Go?" → select + time.After.
*/
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("=== 1. Basic Select First Ready Wins ===")
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(50 * time.Millisecond)
		ch1 <- "one"
	}()
	go func() {
		time.Sleep(30 * time.Millisecond)
		ch2 <- "two"
	}()

	// select blocks until one channel is ready
	select {
	case msg := <-ch1:
		fmt.Printf("  Received from ch1: %s\n", msg)
	case msg := <-ch2:
		fmt.Printf("  Received from ch2: %s\n", msg)
	}

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 2. Timeout with time.After ===")
	slowService := make(chan string)
	go func() {
		time.Sleep(200 * time.Millisecond) // simulate slow work
		slowService <- "result"
	}()

	select {
	case res := <-slowService:
		fmt.Printf("  Got result: %s\n", res)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("  Timeout! Service too slow.")
	}

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 3. Default Case (Non-blocking) ===")
	messages := make(chan string)

	// Non-blocking receive: if nothing is ready, take default
	select {
	case msg := <-messages:
		fmt.Printf("  Got message: %s\n", msg)
	default:
		fmt.Println("  No message available (non-blocking check)")
	}

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 4. Done Channel Pattern (Graceful Shutdown) ===")
	done := make(chan struct{}) // empty struct = zero memory signal

	go func() {
		ticker := time.NewTicker(30 * time.Millisecond)
		defer ticker.Stop()
		count := 0
		for {
			select {
			case <-done:
				fmt.Println("  Worker: received shutdown signal, exiting.")
				return
			case <-ticker.C:
				count++
				fmt.Printf("  Worker: tick %d\n", count)
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)
	close(done) // signal all goroutines listening on done
	time.Sleep(50 * time.Millisecond)

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 5. Select in a Loop Multi-source Aggregation ===")
	fastCh := make(chan int, 10)
	slowCh := make(chan int, 10)
	quitCh := make(chan struct{})

	// Producers
	go func() {
		for i := 0; i < 5; i++ {
			fastCh <- rand.Intn(100)
			time.Sleep(20 * time.Millisecond)
		}
	}()
	go func() {
		for i := 0; i < 3; i++ {
			slowCh <- rand.Intn(100) + 100
			time.Sleep(60 * time.Millisecond)
		}
		close(quitCh) // signal done after slow producer finishes
	}()

	// Consumer with select loop
loop:
	for {
		select {
		case v := <-fastCh:
			fmt.Printf("  Fast: %d\n", v)
		case v := <-slowCh:
			fmt.Printf("  Slow: %d\n", v)
		case <-quitCh:
			fmt.Println("  Quit signal received, breaking loop.")
			break loop
		case <-time.After(150 * time.Millisecond):
			fmt.Println("  No activity for 150ms, breaking.")
			break loop
		}
	}

	// ─────────────────────────────────────────────
	fmt.Println("\n=== Key Takeaways ===")
	fmt.Println("  • select picks a random ready case if multiple are ready")
	fmt.Println("  • time.After creates a one-shot timer channel")
	fmt.Println("  • default makes select non-blocking")
	fmt.Println("  • close(done) broadcasts to ALL listeners simultaneously")
}
