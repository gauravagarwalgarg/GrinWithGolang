/*
What this teaches:
    LRU Cache with O(1) get/put using container/list (doubly-linked list) plus a
    map. Generic LRUCache[K comparable, V any]. Thread-safe option with sync.RWMutex.

Beginner analogy:
    "Like a bookshelf with limited space: when you read a book, move it to the front.
     When the shelf is full and you buy a new book, remove the one at the back
     (least recently used)."

C++ comparison:
    "Same approach as C++ using std::list + std::unordered_map with splice operations.
     Go's container/list provides MoveToFront in O(1). The generic version mirrors
     C++ template<typename K, typename V> LRUCache."

Interview relevance:
    LRU Cache is a top-5 interview question (LeetCode #146). Interviewers expect
    O(1) get and put, understanding of doubly-linked list mechanics, and bonus
    points for thread safety.
*/

package main

import (
	"container/list"
	"fmt"
	"sync"
)

// --- Generic LRU Cache ---

type entry[K comparable, V any] struct {
	key   K
	value V
}

type LRUCache[K comparable, V any] struct {
	capacity int
	items    map[K]*list.Element
	order    *list.List // Front = most recent, Back = least recent
	mu       sync.RWMutex
}

func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity: capacity,
		items:    make(map[K]*list.Element),
		order:    list.New(),
	}
}

// Get retrieves a value and marks it as recently used. O(1).
func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.order.MoveToFront(elem)
		return elem.Value.(*entry[K, V]).value, true
	}
	var zero V
	return zero, false
}

// Put inserts or updates a key-value pair. Evicts LRU if at capacity. O(1).
func (c *LRUCache[K, V]) Put(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Update existing
	if elem, ok := c.items[key]; ok {
		c.order.MoveToFront(elem)
		elem.Value.(*entry[K, V]).value = value
		return
	}

	// Evict if at capacity
	if c.order.Len() >= c.capacity {
		c.evict()
	}

	// Insert new
	e := &entry[K, V]{key: key, value: value}
	elem := c.order.PushFront(e)
	c.items[key] = elem
}

// Delete removes a key from the cache. O(1).
func (c *LRUCache[K, V]) Delete(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.removeElement(elem)
		return true
	}
	return false
}

func (c *LRUCache[K, V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.order.Len()
}

func (c *LRUCache[K, V]) evict() {
	back := c.order.Back()
	if back != nil {
		c.removeElement(back)
	}
}

func (c *LRUCache[K, V]) removeElement(elem *list.Element) {
	c.order.Remove(elem)
	e := elem.Value.(*entry[K, V])
	delete(c.items, e.key)
}

// Keys returns all keys in LRU order (most recent first)
func (c *LRUCache[K, V]) Keys() []K {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]K, 0, c.order.Len())
	for elem := c.order.Front(); elem != nil; elem = elem.Next() {
		keys = append(keys, elem.Value.(*entry[K, V]).key)
	}
	return keys
}

func main() {
	fmt.Println("=== LRU Cache ===")

	cache := NewLRUCache[string, int](3)

	// Fill cache
	fmt.Println("\n--- Insert: a=1, b=2, c=3 (capacity=3) ---")
	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)
	fmt.Printf("  Keys (MRU→LRU): %v\n", cache.Keys())

	// Access 'a' moves to front
	fmt.Println("\n--- Get 'a' → moves to front ---")
	val, ok := cache.Get("a")
	fmt.Printf("  Get('a') = %d, found=%v\n", val, ok)
	fmt.Printf("  Keys (MRU→LRU): %v\n", cache.Keys())

	// Insert 'd' evicts 'b' (LRU)
	fmt.Println("\n--- Put 'd'=4 → evicts 'b' (LRU) ---")
	cache.Put("d", 4)
	fmt.Printf("  Keys (MRU→LRU): %v\n", cache.Keys())
	_, ok = cache.Get("b")
	fmt.Printf("  Get('b') = found=%v (evicted)\n", ok)

	// Update existing
	fmt.Println("\n--- Update 'c'=30 → moves to front ---")
	cache.Put("c", 30)
	val, _ = cache.Get("c")
	fmt.Printf("  Get('c') = %d\n", val)
	fmt.Printf("  Keys (MRU→LRU): %v\n", cache.Keys())

	// Delete
	fmt.Println("\n--- Delete 'a' ---")
	cache.Delete("a")
	fmt.Printf("  Keys: %v, Len: %d\n", cache.Keys(), cache.Len())

	// Key takeaways
	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. Map gives O(1) lookup; list gives O(1) reordering")
	fmt.Println("2. MoveToFront on access; PushFront on insert; Remove Back on evict")
	fmt.Println("3. sync.RWMutex makes it safe for concurrent use")
	fmt.Println("4. Generics: LRUCache[K comparable, V any] reusable for any type")
	fmt.Println("5. container/list uses interface{} internally; generics wrap it type-safely")
}
