/*
Module 10: Distributed - Service Discovery

Demonstrates:
  - In-memory service registry (Register/Discover pattern)
  - Health checking with periodic goroutines
  - Service deregistration on failure
  - Pattern used by Consul, etcd, ZooKeeper
  - Thread-safe registry with sync.RWMutex
  - Load balancing: round-robin selection from instances

Key insight: Service discovery decouples service locations from consumers.
Instead of hardcoded IPs, services register themselves and clients discover
available instances dynamically. Health checks remove dead instances.

Run: go run main.go
*/
package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// ServiceInstance represents a registered service endpoint.
type ServiceInstance struct {
	ID      string
	Name    string
	Address string
	Port    int
	Healthy bool
	LastSeen time.Time
}

// Registry is an in-memory service discovery registry.
type Registry struct {
	mu       sync.RWMutex
	services map[string][]*ServiceInstance // name → instances
}

func NewRegistry() *Registry {
	return &Registry{
		services: make(map[string][]*ServiceInstance),
	}
}

// Register adds a service instance to the registry.
func (r *Registry) Register(instance *ServiceInstance) {
	r.mu.Lock()
	defer r.mu.Unlock()
	instance.Healthy = true
	instance.LastSeen = time.Now()
	r.services[instance.Name] = append(r.services[instance.Name], instance)
	log.Printf("[registry] registered: %s at %s:%d", instance.Name, instance.Address, instance.Port)
}

// Discover returns all healthy instances for a service name.
func (r *Registry) Discover(name string) []*ServiceInstance {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var healthy []*ServiceInstance
	for _, inst := range r.services[name] {
		if inst.Healthy {
			healthy = append(healthy, inst)
		}
	}
	return healthy
}

// Heartbeat updates the last-seen timestamp for an instance.
func (r *Registry) Heartbeat(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, instances := range r.services {
		for _, inst := range instances {
			if inst.ID == id {
				inst.LastSeen = time.Now()
				inst.Healthy = true
				return
			}
		}
	}
}

// healthChecker periodically marks instances as unhealthy if no heartbeat.
func (r *Registry) healthChecker(interval, timeout time.Duration, stop <-chan struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.mu.Lock()
			now := time.Now()
			for _, instances := range r.services {
				for _, inst := range instances {
					if inst.Healthy && now.Sub(inst.LastSeen) > timeout {
						inst.Healthy = false
						log.Printf("[health] marked unhealthy: %s (%s)", inst.ID, inst.Name)
					}
				}
			}
			r.mu.Unlock()
		case <-stop:
			return
		}
	}
}

// roundRobin picks an instance using round-robin (simplified with random).
func roundRobin(instances []*ServiceInstance) *ServiceInstance {
	if len(instances) == 0 {
		return nil
	}
	return instances[rand.Intn(len(instances))]
}

func main() {
	registry := NewRegistry()
	stop := make(chan struct{})

	// Start health checker
	go registry.healthChecker(500*time.Millisecond, 2*time.Second, stop)

	// Register services
	registry.Register(&ServiceInstance{ID: "api-1", Name: "api", Address: "10.0.0.1", Port: 8080})
	registry.Register(&ServiceInstance{ID: "api-2", Name: "api", Address: "10.0.0.2", Port: 8080})
	registry.Register(&ServiceInstance{ID: "db-1", Name: "database", Address: "10.0.1.1", Port: 5432})

	// Simulate heartbeats for api-1 only (api-2 will become unhealthy)
	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(500 * time.Millisecond)
			registry.Heartbeat("api-1")
			registry.Heartbeat("db-1")
		}
	}()

	// Discover services
	time.Sleep(100 * time.Millisecond)
	apis := registry.Discover("api")
	fmt.Printf("\nDiscovered %d healthy API instances\n", len(apis))
	for _, inst := range apis {
		fmt.Printf("  → %s at %s:%d\n", inst.ID, inst.Address, inst.Port)
	}

	// Wait for api-2 to become unhealthy
	time.Sleep(3 * time.Second)
	apis = registry.Discover("api")
	fmt.Printf("\nAfter health check: %d healthy API instances\n", len(apis))
	if inst := roundRobin(apis); inst != nil {
		fmt.Printf("  Selected: %s at %s:%d\n", inst.ID, inst.Address, inst.Port)
	}

	close(stop)
}
