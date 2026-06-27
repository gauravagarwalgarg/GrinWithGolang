/*
What this teaches:
    Segment tree for range sum queries with point updates. Build in O(n), update
    and query in O(log n). Demonstrates the array-based representation where
    children of node i are at 2i and 2i+1.

Beginner analogy:
    "Imagine a tournament bracket where each match stores the sum of its sub-bracket.
     To find the sum of any range, you only check O(log n) bracket nodes instead of
     adding every element individually."

C++ comparison:
    "Same algorithm as competitive programming C++ implementations. Array-based
     segment tree with 1-indexed nodes. Go slice replaces C++ vector; the logic
     is identical. Some C++ libs use iterative (bottom-up) trees for cache efficiency."

Interview relevance:
    Segment trees appear in: range query problems, interval scheduling, count of
    inversions, and competitive programming. Interviewers test understanding of
    build/update/query operations and the O(log n) complexity guarantee.
*/

package main

import "fmt"

// --- Segment Tree ---

type SegmentTree struct {
	tree []int
	n    int
}

// Build constructs the segment tree from an input array. O(n).
func NewSegmentTree(arr []int) *SegmentTree {
	n := len(arr)
	tree := make([]int, 4*n) // 4n is safe upper bound
	st := &SegmentTree{tree: tree, n: n}
	if n > 0 {
		st.build(arr, 1, 0, n-1)
	}
	return st
}

func (st *SegmentTree) build(arr []int, node, start, end int) {
	if start == end {
		st.tree[node] = arr[start]
		return
	}
	mid := (start + end) / 2
	st.build(arr, 2*node, start, mid)
	st.build(arr, 2*node+1, mid+1, end)
	st.tree[node] = st.tree[2*node] + st.tree[2*node+1]
}

// Update sets arr[idx] = val and updates the tree. O(log n).
func (st *SegmentTree) Update(idx, val int) {
	st.update(1, 0, st.n-1, idx, val)
}

func (st *SegmentTree) update(node, start, end, idx, val int) {
	if start == end {
		st.tree[node] = val
		return
	}
	mid := (start + end) / 2
	if idx <= mid {
		st.update(2*node, start, mid, idx, val)
	} else {
		st.update(2*node+1, mid+1, end, idx, val)
	}
	st.tree[node] = st.tree[2*node] + st.tree[2*node+1]
}

// Query returns sum of arr[l..r] (inclusive). O(log n).
func (st *SegmentTree) Query(l, r int) int {
	return st.query(1, 0, st.n-1, l, r)
}

func (st *SegmentTree) query(node, start, end, l, r int) int {
	// No overlap
	if r < start || end < l {
		return 0
	}
	// Complete overlap
	if l <= start && end <= r {
		return st.tree[node]
	}
	// Partial overlap
	mid := (start + end) / 2
	leftSum := st.query(2*node, start, mid, l, r)
	rightSum := st.query(2*node+1, mid+1, end, l, r)
	return leftSum + rightSum
}

// --- Utility: print array ---

func printArr(label string, arr []int) {
	fmt.Printf("  %s: %v\n", label, arr)
}

func main() {
	fmt.Println("=== Segment Tree (Range Sum Query) ===")

	arr := []int{1, 3, 5, 7, 9, 11}
	st := NewSegmentTree(arr)
	printArr("Array", arr)

	// 1. Range queries
	fmt.Println("\n--- Range Sum Queries ---")
	queries := [][2]int{{0, 5}, {1, 3}, {2, 4}, {0, 0}, {3, 5}}
	for _, q := range queries {
		sum := st.Query(q[0], q[1])
		fmt.Printf("  Sum[%d..%d] = %d\n", q[0], q[1], sum)
	}

	// 2. Point update
	fmt.Println("\n--- Point Update: arr[2] = 10 ---")
	st.Update(2, 10) // Change 5 → 10
	arr[2] = 10
	printArr("Updated array", arr)
	fmt.Printf("  Sum[0..5] = %d (was 36, now 41)\n", st.Query(0, 5))
	fmt.Printf("  Sum[1..3] = %d (was 15, now 20)\n", st.Query(1, 3))

	// 3. Multiple updates
	fmt.Println("\n--- Multiple Updates ---")
	st.Update(0, 100)
	arr[0] = 100
	st.Update(5, 50)
	arr[5] = 50
	printArr("Array after updates", arr)
	fmt.Printf("  Sum[0..5] = %d\n", st.Query(0, 5))
	fmt.Printf("  Sum[0..2] = %d\n", st.Query(0, 2))

	// 4. Edge cases
	fmt.Println("\n--- Edge Cases ---")
	fmt.Printf("  Single element Sum[3..3] = %d\n", st.Query(3, 3))
	fmt.Printf("  Full range Sum[0..5] = %d\n", st.Query(0, 5))

	// 5. Build a new tree to verify
	fmt.Println("\n--- Fresh Build: Powers of 2 ---")
	powers := []int{1, 2, 4, 8, 16, 32, 64, 128}
	st2 := NewSegmentTree(powers)
	printArr("Array", powers)
	fmt.Printf("  Sum[0..7] = %d (should be 255)\n", st2.Query(0, 7))
	fmt.Printf("  Sum[4..7] = %d (should be 240)\n", st2.Query(4, 7))
	fmt.Printf("  Sum[2..5] = %d (should be 60)\n", st2.Query(2, 5))

	// Key takeaways
	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. Build: O(n) | Update: O(log n) | Query: O(log n)")
	fmt.Println("2. Array-based: node i has children 2i and 2i+1 (1-indexed)")
	fmt.Println("3. 4n space is safe; exact is 2*nextPow2(n)")
	fmt.Println("4. Extend for: min/max queries, lazy propagation, range updates")
	fmt.Println("5. Use cases: running statistics, interval problems, competitive programming")
}
