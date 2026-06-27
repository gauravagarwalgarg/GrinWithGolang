/*
What this teaches:
    Union-Find (Disjoint Set Union) with path compression and union by rank.
    Near O(1) amortized operations. Applications: connected components, cycle
    detection, Kruskal's MST algorithm. Generic over comparable types.

Beginner analogy:
    "Like friend groups: when two groups meet, one leader becomes the representative
     of the merged group. 'Find' asks 'who's your group leader?' and path
     compression makes future lookups instant."

C++ comparison:
    "Same algorithm as C++ implementations. Path compression flattens the tree during
     find. Union by rank keeps the tree balanced. Together they achieve inverse
     Ackermann α(n) ≈ O(1) amortized time per operation."

Interview relevance:
    Union-Find appears in: number of islands, redundant connections, accounts merge,
    and Kruskal's MST. Interviewers expect path compression + union by rank and
    understanding of the amortized complexity.
*/

package main

import "fmt"

// --- Generic Union-Find ---

type UnionFind[T comparable] struct {
	parent map[T]T
	rank   map[T]int
	count  int // Number of disjoint sets
}

func NewUnionFind[T comparable]() *UnionFind[T] {
	return &UnionFind[T]{
		parent: make(map[T]T),
		rank:   make(map[T]int),
	}
}

// MakeSet creates a new set with element x as its own representative.
func (uf *UnionFind[T]) MakeSet(x T) {
	if _, ok := uf.parent[x]; ok {
		return // Already exists
	}
	uf.parent[x] = x
	uf.rank[x] = 0
	uf.count++
}

// Find returns the representative of x's set with path compression.
func (uf *UnionFind[T]) Find(x T) T {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x]) // Path compression
	}
	return uf.parent[x]
}

// Union merges the sets containing x and y. Returns true if they were separate.
func (uf *UnionFind[T]) Union(x, y T) bool {
	rootX := uf.Find(x)
	rootY := uf.Find(y)

	if rootX == rootY {
		return false // Already in same set
	}

	// Union by rank: attach smaller tree under larger
	switch {
	case uf.rank[rootX] < uf.rank[rootY]:
		uf.parent[rootX] = rootY
	case uf.rank[rootX] > uf.rank[rootY]:
		uf.parent[rootY] = rootX
	default:
		uf.parent[rootY] = rootX
		uf.rank[rootX]++
	}
	uf.count--
	return true
}

// Connected checks if x and y are in the same set.
func (uf *UnionFind[T]) Connected(x, y T) bool {
	return uf.Find(x) == uf.Find(y)
}

// Count returns the number of disjoint sets.
func (uf *UnionFind[T]) Count() int {
	return uf.count
}

// --- Application: Kruskal's MST ---

type Edge struct {
	From, To string
	Weight   int
}

func kruskalMST(nodes []string, edges []Edge) ([]Edge, int) {
	// Sort edges by weight (simple insertion sort for demo)
	sorted := make([]Edge, len(edges))
	copy(sorted, edges)
	for i := 1; i < len(sorted); i++ {
		for j := i; j > 0 && sorted[j].Weight < sorted[j-1].Weight; j-- {
			sorted[j], sorted[j-1] = sorted[j-1], sorted[j]
		}
	}

	uf := NewUnionFind[string]()
	for _, n := range nodes {
		uf.MakeSet(n)
	}

	var mst []Edge
	totalWeight := 0

	for _, e := range sorted {
		if uf.Union(e.From, e.To) {
			mst = append(mst, e)
			totalWeight += e.Weight
			if len(mst) == len(nodes)-1 {
				break
			}
		}
	}
	return mst, totalWeight
}

func main() {
	fmt.Println("=== Union-Find (Disjoint Set Union) ===")

	// 1. Basic operations
	fmt.Println("\n--- Basic Operations ---")
	uf := NewUnionFind[int]()
	for i := 0; i < 7; i++ {
		uf.MakeSet(i)
	}
	fmt.Printf("  Initial sets: %d\n", uf.Count())

	uf.Union(0, 1)
	uf.Union(2, 3)
	uf.Union(4, 5)
	fmt.Printf("  After 3 unions: %d sets\n", uf.Count())

	uf.Union(1, 3) // Merges {0,1} and {2,3}
	fmt.Printf("  After merging groups: %d sets\n", uf.Count())
	fmt.Printf("  Connected(0, 3) = %v\n", uf.Connected(0, 3))
	fmt.Printf("  Connected(0, 4) = %v\n", uf.Connected(0, 4))

	// 2. Connected components (number of islands)
	fmt.Println("\n--- Connected Components ---")
	cities := NewUnionFind[string]()
	for _, c := range []string{"NYC", "LA", "Chicago", "SF", "Seattle", "Portland"} {
		cities.MakeSet(c)
	}
	// Add flight routes
	routes := [][2]string{{"NYC", "Chicago"}, {"LA", "SF"}, {"SF", "Seattle"}, {"Seattle", "Portland"}}
	for _, r := range routes {
		cities.Union(r[0], r[1])
	}
	fmt.Printf("  Number of disconnected regions: %d\n", cities.Count())
	fmt.Printf("  LA connected to Portland? %v\n", cities.Connected("LA", "Portland"))
	fmt.Printf("  NYC connected to LA? %v\n", cities.Connected("NYC", "LA"))

	// 3. Kruskal's MST
	fmt.Println("\n--- Kruskal's MST ---")
	nodes := []string{"A", "B", "C", "D", "E"}
	edges := []Edge{
		{"A", "B", 4}, {"A", "C", 2}, {"B", "C", 1},
		{"B", "D", 5}, {"C", "D", 8}, {"C", "E", 10},
		{"D", "E", 3},
	}
	mst, total := kruskalMST(nodes, edges)
	fmt.Printf("  MST edges (total weight = %d):\n", total)
	for _, e := range mst {
		fmt.Printf("    %s %s (weight %d)\n", e.From, e.To, e.Weight)
	}

	// Key takeaways
	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. Find with path compression: O(α(n)) ≈ O(1) amortized")
	fmt.Println("2. Union by rank keeps trees balanced")
	fmt.Println("3. Applications: connected components, cycle detection, MST")
	fmt.Println("4. Generic: works for any comparable type (int, string, etc.)")
	fmt.Println("5. Count() tracks number of disjoint sets as merges happen")
}
