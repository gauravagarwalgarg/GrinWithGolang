/*
What this teaches:
    Lock-free ring buffer (Single-Producer Single-Consumer) using atomic head/tail.
    Fixed capacity with wrap-around semantics. Used in logging, event queues, and
    inter-goroutine communication where allocation-free operation matters.

Beginner analogy:
    "A circular conveyor belt with fixed slots: the producer places items at the
     tail, the consumer picks from the head. When you reach the end, wrap around
     to the beginning. No locks needed when there's one producer and one consumer."

C++ comparison:
    "Like boost::lockfree::spsc_queue. Uses atomic load/store with memory ordering.
     Go's sync/atomic provides the same guarantees. The key insight: producer only
     writes tail, consumer only writes head no contention."

Interview relevance:
    Ring buffers test understanding of: modular arithmetic for wrap-around, memory
    ordering in lock-free code, and the SPSC constraint that enables lock-freedom.
    Common in systems/embedded interviews.
*/

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// --- Lock-free SPSC Ring Buffer ---

type RingBuffer[T any] struct {
	buffer   []T
	capacity uint64
	head     uint64 // Consumer reads from here (use atomic ops)
	tail     uint64 // Producer writes here (use atomic ops)
}

func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
	// Round up to power of 2 for fast modulo (bitwise AND)
	cap := nextPowerOf2(uint64(capacity))
	return &RingBuffer[T]{
		buffer:   make([]T, cap),
		capacity: cap,
	}
}

// Push adds an item. Returns false if buffer is full (non-blocking).
func (rb *RingBuffer[T]) Push(item T) bool {
	tail := atomic.LoadUint64(&rb.tail)
	head := atomic.LoadUint64(&rb.head)

	// Full when tail is one lap ahead of head
	if tail-head >= rb.capacity {
		return false
	}

	rb.buffer[tail&(rb.capacity-1)] = item // Bitwise AND for modulo
	atomic.StoreUint64(&rb.tail, tail+1)
	return true
}

// Pop removes and returns an item. Returns false if buffer is empty.
func (rb *RingBuffer[T]) Pop() (T, bool) {
	head := atomic.LoadUint64(&rb.head)
	tail := atomic.LoadUint64(&rb.tail)

	if head == tail {
		var zero T
		return zero, false
	}

	item := rb.buffer[head&(rb.capacity-1)]
	atomic.StoreUint64(&rb.head, head+1)
	return item, true
}

// Len returns the number of items currently in the buffer.
func (rb *RingBuffer[T]) Len() int {
	return int(atomic.LoadUint64(&rb.tail) - atomic.LoadUint64(&rb.head))
}

// Cap returns the buffer capacity.
func (rb *RingBuffer[T]) Cap() int {
	return int(rb.capacity)
}

// IsEmpty returns true if the buffer has no items.
func (rb *RingBuffer[T]) IsEmpty() bool {
	return atomic.LoadUint64(&rb.head) == atomic.LoadUint64(&rb.tail)
}

// IsFull returns true if the buffer is at capacity.
func (rb *RingBuffer[T]) IsFull() bool {
	return atomic.LoadUint64(&rb.tail)-atomic.LoadUint64(&rb.head) >= rb.capacity
}

// --- Helper: next power of 2 ---

func nextPowerOf2(n uint64) uint64 {
	if n == 0 {
		return 1
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	return n + 1
}

// --- Demo: log event type ---

type LogEvent struct {
	Level   string
	Message string
}

func main() {
	fmt.Println("=== Lock-Free Ring Buffer (SPSC) ===")

	// 1. Basic operations
	fmt.Println("\n--- Basic Operations ---")
	rb := NewRingBuffer[int](4)
	fmt.Printf("  Capacity: %d (rounded to power of 2)\n", rb.Cap())

	rb.Push(10)
	rb.Push(20)
	rb.Push(30)
	fmt.Printf("  After push 10,20,30: Len=%d\n", rb.Len())

	val, ok := rb.Pop()
	fmt.Printf("  Pop: %d (ok=%v), Len=%d\n", val, ok, rb.Len())

	rb.Push(40)
	rb.Push(50) // Wrap-around!
	fmt.Printf("  After wrap-around push: Len=%d, Full=%v\n", rb.Len(), rb.IsFull())

	// Drain
	fmt.Print("  Drain: ")
	for !rb.IsEmpty() {
		v, _ := rb.Pop()
		fmt.Printf("%d ", v)
	}
	fmt.Println()

	// 2. Full buffer behavior
	fmt.Println("\n--- Full Buffer ---")
	rb2 := NewRingBuffer[string](2)
	fmt.Printf("  Push 'a': %v\n", rb2.Push("a"))
	fmt.Printf("  Push 'b': %v\n", rb2.Push("b"))
	fmt.Printf("  Push 'c' (full): %v\n", rb2.Push("c")) // Returns false

	// 3. SPSC concurrent demo
	fmt.Println("\n--- SPSC Concurrent Demo ---")
	eventBuf := NewRingBuffer[LogEvent](64)
	var wg sync.WaitGroup

	// Producer goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		events := []LogEvent{
			{"INFO", "Server starting"},
			{"INFO", "Listening on :8080"},
			{"WARN", "High memory usage"},
			{"ERROR", "Connection refused"},
			{"INFO", "Request processed"},
		}
		for _, e := range events {
			for !eventBuf.Push(e) {
				time.Sleep(time.Microsecond) // Back-off if full
			}
		}
	}()

	// Consumer goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumed := 0
		for consumed < 5 {
			if event, ok := eventBuf.Pop(); ok {
				fmt.Printf("  [%s] %s\n", event.Level, event.Message)
				consumed++
			} else {
				time.Sleep(time.Microsecond) // Spin-wait
			}
		}
	}()

	wg.Wait()

	// Key takeaways
	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. SPSC = one producer, one consumer → no locks needed")
	fmt.Println("2. Power-of-2 capacity allows bitwise AND instead of modulo")
	fmt.Println("3. Atomic load/store on head/tail provides memory ordering")
	fmt.Println("4. Non-blocking: Push/Pop return false instead of blocking")
	fmt.Println("5. Use cases: logging, metrics, inter-goroutine message passing")
}
