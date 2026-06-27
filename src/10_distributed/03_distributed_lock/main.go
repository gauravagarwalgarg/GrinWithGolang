/*
Module 10: Distributed - Distributed Lock with Lease/TTL

Demonstrates:
  - Lock acquisition with key + TTL (time-to-live)
  - Automatic expiry goroutine (prevents deadlocks from crashed holders)
  - Unlock by key (only holder can release)
  - Pattern used by Redis SETNX + EXPIRE, etcd leases, ZooKeeper ephemerals
  - Fencing token concept to prevent stale lock holders

Key insight: Distributed locks MUST have TTL/lease. Without it, a crashed
process holds the lock forever (deadlock). The tradeoff: if TTL is too short,
the lock expires while work is still in progress (split-brain risk).
Solution: use fencing tokens to detect stale lock holders.

Run: go run main.go
*/
package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// LockEntry represents a held lock with expiry metadata.
type LockEntry struct {
	Key       string
	Holder    string
	Token     uint64 // fencing token (monotonically increasing)
	ExpiresAt time.Time
}

// DistributedLock simulates a distributed lock manager (like Redis/etcd).
type DistributedLock struct {
	mu      sync.Mutex
	locks   map[string]*LockEntry
	nextTok uint64
}

func NewDistributedLock() *DistributedLock {
	dl := &DistributedLock{
		locks: make(map[string]*LockEntry),
	}
	go dl.expiryLoop()
	return dl
}

// Lock acquires a lock on key with given TTL. Returns fencing token or error.
func (dl *DistributedLock) Lock(key, holder string, ttl time.Duration) (uint64, error) {
	dl.mu.Lock()
	defer dl.mu.Unlock()

	if existing, ok := dl.locks[key]; ok {
		if time.Now().Before(existing.ExpiresAt) {
			return 0, fmt.Errorf("lock %q held by %s (expires in %v)",
				key, existing.Holder, time.Until(existing.ExpiresAt).Round(time.Millisecond))
		}
		// Expired: allow acquisition
		log.Printf("[lock] expired lock on %q (was held by %s)", key, existing.Holder)
	}

	dl.nextTok++
	dl.locks[key] = &LockEntry{
		Key:       key,
		Holder:    holder,
		Token:     dl.nextTok,
		ExpiresAt: time.Now().Add(ttl),
	}
	log.Printf("[lock] acquired: key=%s holder=%s token=%d ttl=%v",
		key, holder, dl.nextTok, ttl)
	return dl.nextTok, nil
}

// Unlock releases a lock. Only the holder can unlock.
func (dl *DistributedLock) Unlock(key, holder string) error {
	dl.mu.Lock()
	defer dl.mu.Unlock()

	entry, ok := dl.locks[key]
	if !ok {
		return fmt.Errorf("lock %q not found", key)
	}
	if entry.Holder != holder {
		return fmt.Errorf("lock %q held by %s, not %s", key, entry.Holder, holder)
	}

	delete(dl.locks, key)
	log.Printf("[lock] released: key=%s holder=%s", key, holder)
	return nil
}

// expiryLoop periodically removes expired locks.
func (dl *DistributedLock) expiryLoop() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		dl.mu.Lock()
		now := time.Now()
		for key, entry := range dl.locks {
			if now.After(entry.ExpiresAt) {
				log.Printf("[lock] auto-expired: key=%s holder=%s", key, entry.Holder)
				delete(dl.locks, key)
			}
		}
		dl.mu.Unlock()
	}
}

func main() {
	dm := NewDistributedLock()

	// Scenario 1: Normal lock/unlock
	fmt.Println("=== Scenario 1: Normal Lock/Unlock ===")
	token, _ := dm.Lock("resource:orders", "worker-1", 5*time.Second)
	fmt.Printf("Worker-1 acquired lock with token: %d\n", token)

	// Scenario 2: Contention - second worker tries to acquire
	fmt.Println("\n=== Scenario 2: Lock Contention ===")
	_, err := dm.Lock("resource:orders", "worker-2", 5*time.Second)
	fmt.Printf("Worker-2 attempt: %v\n", err)

	// Release and retry
	dm.Unlock("resource:orders", "worker-1")
	token2, _ := dm.Lock("resource:orders", "worker-2", 5*time.Second)
	fmt.Printf("Worker-2 acquired after release, token: %d\n", token2)
	dm.Unlock("resource:orders", "worker-2")

	// Scenario 3: TTL expiry (simulates crashed process)
	fmt.Println("\n=== Scenario 3: TTL Expiry (crashed holder) ===")
	dm.Lock("resource:payments", "worker-3", 500*time.Millisecond)
	fmt.Println("Worker-3 acquired with 500ms TTL (simulating crash...)")
	time.Sleep(700 * time.Millisecond)

	token3, _ := dm.Lock("resource:payments", "worker-4", 5*time.Second)
	fmt.Printf("Worker-4 acquired expired lock, token: %d\n", token3)

	time.Sleep(100 * time.Millisecond)
}
