/*
What this teaches:
    Sharded concurrent map better than sync.Map for write-heavy workloads. Shards
    by hash of key, each shard with its own mutex. Reduces lock contention by
    distributing writes across N independent locks.

Beginner analogy:
    "Like a library with 16 sections: instead of one librarian handling all returns,
     each section has its own desk. Multiple people can return books simultaneously
     as long as they're in different sections."

C++ comparison:
    "Like Java's ConcurrentHashMap or tbb::concurrent_hash_map with segment-level
     locking. Go's sync.Map is optimized for read-heavy/write-once; this sharded
     approach wins when writes are frequent."

Interview relevance:
    Interviewers ask when sync.Map is appropriate vs alternatives. Demonstrating a
    sharded map shows understanding of lock granularity, hash distribution, and
    concurrency vs contention tradeoffs.
*/

package main

import (
	"fmt"
	"hash/fnv"
	"sync"
	"time"
)

// --- Sharded Concurrent Map ---

const defaultShardCount = 16

type shard[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]V
}

type ShardedMap[K comparable, V any] struct {
	shards    []*shard[K, V]
	shardCount int
	hashFn    func(K) uint64
}

func NewShardedMap[K comparable, V any](shardCount int, hashFn func(K) uint64) *ShardedMap[K, V] {
	if shardCount <= 0 {
		shardCount = defaultShardCount
	}
	shards := make([]*shard[K, V], shardCount)
	for i := range shards {
		shards[i] = &shard[K, V]{items: make(map[K]V)}
	}
	return &ShardedMap[K, V]{
		shards:     shards,
		shardCount: shardCount,
		hashFn:     hashFn,
	}
}

func (m *ShardedMap[K, V]) getShard(key K) *shard[K, V] {
	h := m.hashFn(key)
	return m.shards[h%uint64(m.shardCount)]
}

// Set stores a key-value pair.
func (m *ShardedMap[K, V]) Set(key K, value V) {
	s := m.getShard(key)
	s.mu.Lock()
	s.items[key] = value
	s.mu.Unlock()
}

// Get retrieves a value by key.
func (m *ShardedMap[K, V]) Get(key K) (V, bool) {
	s := m.getShard(key)
	s.mu.RLock()
	val, ok := s.items[key]
	s.mu.RUnlock()
	return val, ok
}

// Delete removes a key.
func (m *ShardedMap[K, V]) Delete(key K) {
	s := m.getShard(key)
	s.mu.Lock()
	delete(s.items, key)
	s.mu.Unlock()
}

// Len returns total count across all shards.
func (m *ShardedMap[K, V]) Len() int {
	total := 0
	for _, s := range m.shards {
		s.mu.RLock()
		total += len(s.items)
		s.mu.RUnlock()
	}
	return total
}

// Keys returns all keys (snapshot, not live).
func (m *ShardedMap[K, V]) Keys() []K {
	keys := make([]K, 0)
	for _, s := range m.shards {
		s.mu.RLock()
		for k := range s.items {
			keys = append(keys, k)
		}
		s.mu.RUnlock()
	}
	return keys
}

// --- Hash function for string keys ---

func stringHash(key string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(key))
	return h.Sum64()
}

// --- Benchmark helper ---

func benchmarkMap(name string, fn func()) time.Duration {
	start := time.Now()
	fn()
	elapsed := time.Since(start)
	fmt.Printf("  %-20s %v\n", name+":", elapsed)
	return elapsed
}

func main() {
	fmt.Println("=== Sharded Concurrent Map ===")

	// 1. Basic operations
	fmt.Println("\n--- Basic Operations ---")
	m := NewShardedMap[string, int](16, stringHash)
	m.Set("alpha", 1)
	m.Set("beta", 2)
	m.Set("gamma", 3)

	v, ok := m.Get("beta")
	fmt.Printf("  Get('beta') = %d, found=%v\n", v, ok)
	fmt.Printf("  Len = %d\n", m.Len())

	m.Delete("beta")
	_, ok = m.Get("beta")
	fmt.Printf("  After Delete('beta'): found=%v, Len=%d\n", ok, m.Len())

	// 2. Concurrent write stress test
	fmt.Println("\n--- Concurrent Write Benchmark ---")
	const numOps = 100_000
	const numGoroutines = 8

	// Sharded map
	shardedMap := NewShardedMap[string, int](16, stringHash)
	benchmarkMap("ShardedMap", func() {
		var wg sync.WaitGroup
		for g := 0; g < numGoroutines; g++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for i := 0; i < numOps/numGoroutines; i++ {
					key := fmt.Sprintf("key-%d-%d", id, i)
					shardedMap.Set(key, i)
				}
			}(g)
		}
		wg.Wait()
	})

	// sync.Map comparison
	var syncMap sync.Map
	benchmarkMap("sync.Map", func() {
		var wg sync.WaitGroup
		for g := 0; g < numGoroutines; g++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for i := 0; i < numOps/numGoroutines; i++ {
					key := fmt.Sprintf("key-%d-%d", id, i)
					syncMap.Store(key, i)
				}
			}(g)
		}
		wg.Wait()
	})

	// Single-mutex map comparison
	var muMap sync.Mutex
	plainMap := make(map[string]int)
	benchmarkMap("Mutex+Map", func() {
		var wg sync.WaitGroup
		for g := 0; g < numGoroutines; g++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for i := 0; i < numOps/numGoroutines; i++ {
					key := fmt.Sprintf("key-%d-%d", id, i)
					muMap.Lock()
					plainMap[key] = i
					muMap.Unlock()
				}
			}(g)
		}
		wg.Wait()
	})

	fmt.Printf("\n  ShardedMap final size: %d\n", shardedMap.Len())

	// Key takeaways
	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. Sharding distributes lock contention across N independent mutexes")
	fmt.Println("2. sync.Map is optimized for read-heavy, write-once patterns")
	fmt.Println("3. ShardedMap wins for write-heavy concurrent workloads")
	fmt.Println("4. Shard count should be ≥ GOMAXPROCS for best parallelism")
	fmt.Println("5. FNV hash gives good distribution; power-of-2 shards allow bitwise AND")
}
