/*
What this teaches:
    Complete fan-in/fan-out pipeline with stage cancellation via context. Shows an
    image processing pipeline (resize → compress → upload) with bounded concurrency
    per stage.

Beginner analogy:
    "An assembly line with multiple stations: raw materials fan OUT to parallel
     workers at each station, then FAN IN to a single output. If one station jams,
     the whole line stops gracefully."

C++ comparison:
    "Like TBB's parallel_pipeline or a series of concurrent queues. Go channels
     replace the queue abstraction, and select+context replace cancellation tokens.
     No need for thread-safe queue libraries."

Interview relevance:
    Fan-in/fan-out is critical for data pipelines and microservice patterns.
    Interviewers ask about bounded concurrency, backpressure, and graceful shutdown
    this demonstrates all three.
*/

package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// --- Pipeline data types ---

type Image struct {
	ID   int
	Name string
	Size int // bytes
}

type ResizedImage struct {
	Image
	NewSize int
}

type CompressedImage struct {
	ResizedImage
	CompressedSize int
}

type UploadResult struct {
	ImageID int
	URL     string
	Err     error
}

// --- Stage 1: Generate source images ---

func generateImages(ctx context.Context, count int) <-chan Image {
	out := make(chan Image)
	go func() {
		defer close(out)
		for i := 1; i <= count; i++ {
			img := Image{ID: i, Name: fmt.Sprintf("photo_%03d.jpg", i), Size: 5000000 + rand.Intn(5000000)}
			select {
			case out <- img:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// --- Stage 2: Resize (fan-out to N workers) ---

func resize(ctx context.Context, images <-chan Image, workers int) <-chan ResizedImage {
	out := make(chan ResizedImage)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for img := range images {
				select {
				case <-ctx.Done():
					return
				default:
				}
				time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
				resized := ResizedImage{Image: img, NewSize: img.Size / 2}
				select {
				case out <- resized:
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	go func() { wg.Wait(); close(out) }()
	return out
}

// --- Stage 3: Compress (fan-out to N workers) ---

func compress(ctx context.Context, images <-chan ResizedImage, workers int) <-chan CompressedImage {
	out := make(chan CompressedImage)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for img := range images {
				select {
				case <-ctx.Done():
					return
				default:
				}
				time.Sleep(time.Duration(rand.Intn(30)) * time.Millisecond)
				compressed := CompressedImage{ResizedImage: img, CompressedSize: img.NewSize / 3}
				select {
				case out <- compressed:
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	go func() { wg.Wait(); close(out) }()
	return out
}

// --- Stage 4: Upload (fan-out to N workers, fan-in results) ---

func upload(ctx context.Context, images <-chan CompressedImage, workers int) <-chan UploadResult {
	out := make(chan UploadResult)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for img := range images {
				select {
				case <-ctx.Done():
					out <- UploadResult{ImageID: img.ID, Err: ctx.Err()}
					return
				default:
				}
				time.Sleep(time.Duration(rand.Intn(40)) * time.Millisecond)
				url := fmt.Sprintf("https://cdn.example.com/%s", img.Name)
				select {
				case out <- UploadResult{ImageID: img.ID, URL: url}:
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	go func() { wg.Wait(); close(out) }()
	return out
}

func main() {
	fmt.Println("=== Fan-In / Fan-Out Pipeline ===")
	rand.Seed(time.Now().UnixNano())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Pipeline: generate → resize(3) → compress(2) → upload(2)
	images := generateImages(ctx, 10)
	resized := resize(ctx, images, 3)
	compressed := compress(ctx, resized, 2)
	results := upload(ctx, compressed, 2)

	// Fan-in: collect all results
	fmt.Println("\n--- Upload Results ---")
	var success, errCount int
	for r := range results {
		if r.Err != nil {
			fmt.Printf("  Image %d: FAILED %v\n", r.ImageID, r.Err)
			errCount++
		} else {
			fmt.Printf("  Image %d: OK → %s\n", r.ImageID, r.URL)
			success++
		}
	}

	fmt.Printf("\n--- Pipeline Complete: %d uploaded, %d failed ---\n", success, errCount)
	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. Each stage is a goroutine reading from input channel, writing to output")
	fmt.Println("2. Fan-out: multiple goroutines read from same channel")
	fmt.Println("3. Fan-in: WaitGroup + close(out) merges worker outputs")
	fmt.Println("4. Context cancellation propagates through all stages")
	fmt.Println("5. Bounded concurrency per stage prevents resource exhaustion")
}
