/*
What:       Bounded worker pool N workers pulling jobs from a channel, graceful shutdown
Level:      Beginner
Analogy:    A restaurant kitchen: N cooks take orders from a ticket queue. No cook is idle, no overflow.
C++ Angle:  Like a thread pool with a bounded queue, but trivial to implement with channels.
Interview:  "How do you limit concurrency in Go?" → Worker pool with buffered job channel.
*/
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Job represents a unit of work.
type Job struct {
	ID       int
	Duration time.Duration
}

// Result holds the output of a processed job.
type Result struct {
	JobID    int
	WorkerID int
	Output   string
}

// worker pulls jobs from the jobs channel until it's closed.
func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		// Simulate processing
		time.Sleep(job.Duration)
		results <- Result{
			JobID:    job.ID,
			WorkerID: id,
			Output:   fmt.Sprintf("processed in %v", job.Duration),
		}
	}
}

func main() {
	const (
		numWorkers = 3
		numJobs    = 10
	)

	fmt.Printf("=== Worker Pool: %d workers, %d jobs ===\n", numWorkers, numJobs)
	fmt.Println("  Workers pull from a shared job channel (bounded concurrency).")
	fmt.Println()

	jobs := make(chan Job, numJobs)       // buffered: enqueue all jobs up front
	results := make(chan Result, numJobs) // collect results

	// ─── Start Workers ───────────────────────────────────────────
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// ─── Submit Jobs ─────────────────────────────────────────────
	start := time.Now()
	for j := 1; j <= numJobs; j++ {
		jobs <- Job{
			ID:       j,
			Duration: time.Duration(rand.Intn(50)+10) * time.Millisecond,
		}
	}
	close(jobs) // signal workers: no more jobs coming

	// ─── Wait for completion and close results ───────────────────
	go func() {
		wg.Wait()
		close(results) // safe to close after all workers done
	}()

	// ─── Collect Results ─────────────────────────────────────────
	for r := range results {
		fmt.Printf("  Job %2d → Worker %d: %s\n", r.JobID, r.WorkerID, r.Output)
	}

	elapsed := time.Since(start)
	fmt.Printf("\n  Total time: %v (parallel speedup vs sequential)\n", elapsed)
	fmt.Println("\n=== Graceful Shutdown Notes ===")
	fmt.Println("  1. close(jobs) tells workers to exit their range loop")
	fmt.Println("  2. wg.Wait() ensures all in-flight jobs complete")
	fmt.Println("  3. Then close(results) so collector's range loop exits")
	fmt.Println("  4. For cancellation: add a done channel or use context.WithCancel")
}
