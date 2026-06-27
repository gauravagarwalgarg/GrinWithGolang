/*
What this teaches:
    Generic priority queue using container/heap interface. Implements both MinHeap
    and MaxHeap. Use cases: task scheduling, Dijkstra's algorithm, top-K elements.

Beginner analogy:
    "A hospital ER waiting room: patients aren't served FIFO the most critical
     (highest priority) always goes next. The heap ensures extracting the highest
     priority is always O(log n)."

C++ comparison:
    "Like std::priority_queue<T, Container, Compare>. Go's container/heap requires
     implementing the heap.Interface (Len, Less, Swap, Push, Pop). The generic
     wrapper provides a cleaner API over the raw interface."

Interview relevance:
    Priority queues appear in: merge K sorted lists, task schedulers, median finding,
    and graph algorithms (Dijkstra, Prim). Interviewers expect O(log n) push/pop
    and understanding of the underlying binary heap.
*/

package main

import (
	"container/heap"
	"fmt"
)

// --- Generic Heap (internal) ---

type heapSlice[T any] struct {
	data []T
	less func(a, b T) bool
}

func (h *heapSlice[T]) Len() int           { return len(h.data) }
func (h *heapSlice[T]) Less(i, j int) bool { return h.less(h.data[i], h.data[j]) }
func (h *heapSlice[T]) Swap(i, j int)      { h.data[i], h.data[j] = h.data[j], h.data[i] }

func (h *heapSlice[T]) Push(x any) {
	h.data = append(h.data, x.(T))
}

func (h *heapSlice[T]) Pop() any {
	old := h.data
	n := len(old)
	x := old[n-1]
	h.data = old[:n-1]
	return x
}

// --- Public PriorityQueue API ---

type PriorityQueue[T any] struct {
	h *heapSlice[T]
}

func NewMinHeap[T any](less func(a, b T) bool) *PriorityQueue[T] {
	pq := &PriorityQueue[T]{h: &heapSlice[T]{less: less}}
	heap.Init(pq.h)
	return pq
}

func NewMaxHeap[T any](less func(a, b T) bool) *PriorityQueue[T] {
	// Invert comparison for max heap
	return NewMinHeap(func(a, b T) bool { return less(b, a) })
}

func (pq *PriorityQueue[T]) Push(item T) {
	heap.Push(pq.h, item)
}

func (pq *PriorityQueue[T]) Pop() T {
	return heap.Pop(pq.h).(T)
}

func (pq *PriorityQueue[T]) Peek() T {
	return pq.h.data[0]
}

func (pq *PriorityQueue[T]) Len() int {
	return pq.h.Len()
}

// --- Domain types ---

type Task struct {
	Name     string
	Priority int
}

type Edge struct {
	Node string
	Dist int
}

func main() {
	fmt.Println("=== Generic Priority Queue (Heap) ===")

	// 1. MinHeap of integers
	fmt.Println("\n--- MinHeap[int] ---")
	minPQ := NewMinHeap(func(a, b int) bool { return a < b })
	for _, v := range []int{42, 15, 88, 3, 67, 22} {
		minPQ.Push(v)
	}
	fmt.Print("  Sorted (min first): ")
	for minPQ.Len() > 0 {
		fmt.Printf("%d ", minPQ.Pop())
	}
	fmt.Println()

	// 2. MaxHeap of integers
	fmt.Println("\n--- MaxHeap[int] ---")
	maxPQ := NewMaxHeap(func(a, b int) bool { return a < b })
	for _, v := range []int{42, 15, 88, 3, 67, 22} {
		maxPQ.Push(v)
	}
	fmt.Print("  Sorted (max first): ")
	for maxPQ.Len() > 0 {
		fmt.Printf("%d ", maxPQ.Pop())
	}
	fmt.Println()

	// 3. Task scheduler (higher priority number = more urgent)
	fmt.Println("\n--- Task Scheduler (MaxHeap by priority) ---")
	taskPQ := NewMaxHeap(func(a, b Task) bool { return a.Priority < b.Priority })
	tasks := []Task{
		{"Send email", 1},
		{"Fix crash", 10},
		{"Update docs", 3},
		{"Deploy hotfix", 9},
		{"Refactor tests", 2},
	}
	for _, t := range tasks {
		taskPQ.Push(t)
	}
	fmt.Println("  Execution order:")
	for taskPQ.Len() > 0 {
		t := taskPQ.Pop()
		fmt.Printf("    [P%d] %s\n", t.Priority, t.Name)
	}

	// 4. Top-K smallest (use MaxHeap of size K)
	fmt.Println("\n--- Top-3 Smallest Elements ---")
	k := 3
	topK := NewMaxHeap(func(a, b int) bool { return a < b })
	stream := []int{50, 10, 80, 30, 90, 5, 70, 20}
	for _, v := range stream {
		topK.Push(v)
		if topK.Len() > k {
			topK.Pop() // Remove the largest in our K-sized max heap
		}
	}
	fmt.Printf("  Stream: %v\n", stream)
	fmt.Print("  Top-3 smallest: ")
	for topK.Len() > 0 {
		fmt.Printf("%d ", topK.Pop())
	}
	fmt.Println()

	// 5. Dijkstra-style edge processing
	fmt.Println("\n--- Dijkstra Edge Queue (MinHeap by distance) ---")
	edgePQ := NewMinHeap(func(a, b Edge) bool { return a.Dist < b.Dist })
	edges := []Edge{{"B", 4}, {"C", 1}, {"D", 7}, {"E", 2}}
	for _, e := range edges {
		edgePQ.Push(e)
	}
	fmt.Println("  Processing edges nearest-first:")
	for edgePQ.Len() > 0 {
		e := edgePQ.Pop()
		fmt.Printf("    → %s (dist=%d)\n", e.Node, e.Dist)
	}

	// Key takeaways
	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. container/heap requires: Len, Less, Swap, Push, Pop")
	fmt.Println("2. MinHeap: Less(a,b) = a < b; MaxHeap: invert the comparator")
	fmt.Println("3. Push/Pop are O(log n); Peek is O(1)")
	fmt.Println("4. Top-K: maintain a heap of size K, evict when exceeding")
	fmt.Println("5. Generic wrapper hides the interface{} ceremony of container/heap")
}
