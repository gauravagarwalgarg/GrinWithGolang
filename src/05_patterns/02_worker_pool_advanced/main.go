/*
What this teaches:
    Production-grade worker pool: generic job processor with result collection,
    error handling, graceful drain, and context cancellation. Reusable Pool[T, R]
    struct that can process any job type.

Beginner analogy:
    "A factory assembly line: jobs arrive on a conveyor belt (channel), N workers
     pick them up, process them, and place results on the output belt. If the boss
     says stop (context cancel), everyone wraps up their current item and goes home."

C++ comparison:
    "Like a thread pool with std::async/futures, but channels replace condition
     variables and the generic Pool[T, R] replaces templated executors. No mutex
     contention on the queue channels handle synchronization."

Interview relevance:
    Worker pools are the #1 concurrency pattern asked in Go interviews. Demonstrate
    bounded concurrency, graceful shutdown, error propagation, and generic reuse.
*/

package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// --- Generic Worker Pool ---

type Job[T any] struct {
	ID      int
	Payload T
}

type Result[R any] struct {
	JobID int
	Value R
	Err   error
}

type Pool[T any, R any] struct {
	workers    int
	processor  func(context.Context, T) (R, error)
	jobs       chan Job[T]
	results    chan Result[R]
}

func NewPool[T any, R any](workers int, bufferSize int, fn func(context.Context, T) (R, error)) *Pool[T, R] {
	return &Pool[T, R]{
		workers:   workers,
		processor: fn,
		jobs:      make(chan Job[T], bufferSize),
		results:   make(chan Result[R], bufferSize),
	}
}

func (p *Pool[T, R]) Start(ctx context.Context) {
	var wg sync.WaitGroup
	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range p.jobs {
				select {
				case <-ctx.Done():
					p.results <- Result[R]{JobID: job.ID, Err: ctx.Err()}
					return
				default:
				}
				value, err := p.processor(ctx, job.Payload)
				p.results <- Result[R]{JobID: job.ID, Value: value, Err: err}
			}
		}(i)
	}

	// Close results when all workers are done
	go func() {
		wg.Wait()
		close(p.results)
	}()
}

func (p *Pool[T, R]) Submit(job Job[T]) {
	p.jobs <- job
}

func (p *Pool[T, R]) Close() {
	close(p.jobs)
}

func (p *Pool[T, R]) Results() <-chan Result[R] {
	return p.results
}

// --- Example: Image processing simulation ---

type ImageTask struct {
	Filename string
	Width    int
}

type ImageResult struct {
	Filename string
	Size     int
}

func processImage(ctx context.Context, task ImageTask) (ImageResult, error) {
	// Simulate work
	duration := time.Duration(rand.Intn(100)) * time.Millisecond
	select {
	case <-time.After(duration):
	case <-ctx.Done():
		return ImageResult{}, ctx.Err()
	}

	// Simulate occasional failures
	if rand.Float64() < 0.1 {
		return ImageResult{}, fmt.Errorf("codec error for %s", task.Filename)
	}

	return ImageResult{
		Filename: task.Filename,
		Size:     task.Width * task.Width / 4, // Simulated compressed size
	}, nil
}

func main() {
	fmt.Println("=== Advanced Worker Pool ===")
	rand.Seed(time.Now().UnixNano())

	// Create a cancellable context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Create pool: 4 workers, process ImageTask → ImageResult
	pool := NewPool[ImageTask, ImageResult](4, 20, processImage)
	pool.Start(ctx)

	// Submit jobs
	fmt.Println("\n--- Submitting 15 image jobs ---")
	go func() {
		for i := 1; i <= 15; i++ {
			pool.Submit(Job[ImageTask]{
				ID:      i,
				Payload: ImageTask{Filename: fmt.Sprintf("img_%03d.png", i), Width: 1920},
			})
		}
		pool.Close() // Signal no more jobs
	}()

	// Collect results
	fmt.Println("\n--- Results ---")
	var succeeded, failed int
	for result := range pool.Results() {
		if result.Err != nil {
			fmt.Printf("  Job %2d: ERROR %v\n", result.JobID, result.Err)
			failed++
		} else {
			fmt.Printf("  Job %2d: OK %s (%d bytes)\n", result.JobID, result.Value.Filename, result.Value.Size)
			succeeded++
		}
	}

	fmt.Printf("\n--- Summary: %d succeeded, %d failed ---\n", succeeded, failed)

	// Key takeaways
	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. Pool[T, R] is generic works for any job/result type")
	fmt.Println("2. Context propagation enables timeout and cancellation")
	fmt.Println("3. Buffered channels prevent producer/consumer blocking")
	fmt.Println("4. WaitGroup + close(results) signals completion to collector")
	fmt.Println("5. Graceful drain: workers finish current job before exiting")
}
