/*
What:       Pipeline, fan-out, fan-in, and done/cancellation channel patterns
Level:      Beginner
Analogy:    Pipeline = assembly line. Fan-out = one boss, many workers. Fan-in = many reporters, one editor.
C++ Angle:  These patterns replace complex thread pool + queue architectures with composable channels.
Interview:  "Design a concurrent pipeline in Go" → chain of channels with goroutines per stage.
*/
package main

import (
	"fmt"
	"sync"
	"time"
)

// ─── Pipeline: stage1 → stage2 → stage3 ──────────────────────────
func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func addTen(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n + 10
		}
		close(out)
	}()
	return out
}

// ─── Fan-out: one producer, N workers ─────────────────────────────
func fanOutWorker(id int, jobs <-chan int, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		time.Sleep(10 * time.Millisecond) // simulate work
		results <- fmt.Sprintf("worker %d processed job %d → %d", id, job, job*2)
	}
}

// ─── Fan-in: merge N channels into one ────────────────────────────
func fanIn(channels ...<-chan string) <-chan string {
	merged := make(chan string)
	var wg sync.WaitGroup
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan string) {
			defer wg.Done()
			for msg := range c {
				merged <- msg
			}
		}(ch)
	}
	go func() {
		wg.Wait()
		close(merged)
	}()
	return merged
}

func makeProducer(name string, count int) <-chan string {
	ch := make(chan string)
	go func() {
		for i := 0; i < count; i++ {
			ch <- fmt.Sprintf("%s-%d", name, i)
			time.Sleep(20 * time.Millisecond)
		}
		close(ch)
	}()
	return ch
}

func main() {
	fmt.Println("=== 1. Pipeline Pattern ===")
	// Data flows: generate → square → addTen → consumer
	pipeline := addTen(square(generate(1, 2, 3, 4, 5)))
	for result := range pipeline {
		fmt.Printf("  %d\n", result) // 11, 14, 19, 26, 35
	}

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 2. Fan-Out (1 Producer, N Workers) ===")
	jobs := make(chan int, 10)
	results := make(chan string, 10)
	var wg sync.WaitGroup

	// Start 3 workers
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go fanOutWorker(w, jobs, results, &wg)
	}

	// Send jobs
	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	close(jobs)

	// Close results after workers finish
	go func() { wg.Wait(); close(results) }()
	for r := range results {
		fmt.Printf("  %s\n", r)
	}

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 3. Fan-In (N Producers, 1 Consumer) ===")
	p1 := makeProducer("alpha", 3)
	p2 := makeProducer("beta", 3)
	p3 := makeProducer("gamma", 3)

	for msg := range fanIn(p1, p2, p3) {
		fmt.Printf("  %s\n", msg)
	}

	// ─────────────────────────────────────────────
	fmt.Println("\n=== 4. Done/Cancellation Channel ===")
	done := make(chan struct{})
	out := make(chan int)

	go func() {
		defer close(out)
		for i := 0; ; i++ {
			select {
			case <-done:
				fmt.Println("  Producer: cancelled, cleaning up.")
				return
			case out <- i:
				time.Sleep(20 * time.Millisecond)
			}
		}
	}()

	// Consume 5 values then cancel
	for i := 0; i < 5; i++ {
		fmt.Printf("  Consumed: %d\n", <-out)
	}
	close(done) // broadcast cancellation
	time.Sleep(50 * time.Millisecond)
	fmt.Println("  Main: done.")
}
