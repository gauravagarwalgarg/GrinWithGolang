/*
What this teaches:
    In-memory pub/sub event bus: Subscribe(topic, handler), Publish(topic, event).
    Type-safe with generics. Used for decoupled communication within a service
    where components shouldn't know about each other directly.

Beginner analogy:
    "Like a newspaper subscription: you subscribe to 'Sports' or 'Tech'. When a new
     article (event) is published to that topic, all subscribers get notified 
     without the writer knowing who's reading."

C++ comparison:
    "Similar to Qt's signal/slot or the Observer pattern with std::function callbacks.
     Go's generics + channels make it type-safe without templates or type erasure.
     The event bus decouples producers from consumers."

Interview relevance:
    Pub/sub is fundamental to event-driven architectures. Interviewers ask about
    thread safety, subscriber lifecycle, and how to prevent goroutine leaks when
    subscribers unsubscribe.
*/

package main

import (
	"fmt"
	"sync"
	"time"
)

// --- Generic Event Bus ---

type EventBus[T any] struct {
	mu          sync.RWMutex
	subscribers map[string][]chan T
	bufferSize  int
}

func NewEventBus[T any](bufferSize int) *EventBus[T] {
	return &EventBus[T]{
		subscribers: make(map[string][]chan T),
		bufferSize:  bufferSize,
	}
}

// Subscribe returns a channel that receives events for the given topic
func (eb *EventBus[T]) Subscribe(topic string) <-chan T {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	ch := make(chan T, eb.bufferSize)
	eb.subscribers[topic] = append(eb.subscribers[topic], ch)
	return ch
}

// Publish sends an event to all subscribers of the topic
func (eb *EventBus[T]) Publish(topic string, event T) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	subs, ok := eb.subscribers[topic]
	if !ok {
		return
	}
	for _, ch := range subs {
		select {
		case ch <- event:
		default:
			// Drop event if subscriber is slow (non-blocking)
			fmt.Printf("  [WARN] Subscriber on %q is slow, event dropped\n", topic)
		}
	}
}

// Unsubscribe removes a channel from a topic
func (eb *EventBus[T]) Unsubscribe(topic string, ch <-chan T) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	subs := eb.subscribers[topic]
	for i, sub := range subs {
		if sub == ch {
			eb.subscribers[topic] = append(subs[:i], subs[i+1:]...)
			close(sub)
			return
		}
	}
}

// Close shuts down all subscriber channels
func (eb *EventBus[T]) Close() {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	for topic, subs := range eb.subscribers {
		for _, ch := range subs {
			close(ch)
		}
		delete(eb.subscribers, topic)
	}
}

// --- Domain events ---

type OrderEvent struct {
	OrderID   string
	Action    string
	Amount    float64
	Timestamp time.Time
}

func main() {
	fmt.Println("=== Pub/Sub Event Bus ===")

	bus := NewEventBus[OrderEvent](10)

	// Subscriber 1: Analytics service
	fmt.Println("\n--- Setting up subscribers ---")
	analyticsCh := bus.Subscribe("orders")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for event := range analyticsCh {
			fmt.Printf("  [Analytics] %s: order %s ($%.2f)\n",
				event.Action, event.OrderID, event.Amount)
		}
		fmt.Println("  [Analytics] Subscriber shut down")
	}()

	// Subscriber 2: Notification service
	notifyCh := bus.Subscribe("orders")
	wg.Add(1)
	go func() {
		defer wg.Done()
		for event := range notifyCh {
			if event.Amount > 100 {
				fmt.Printf("  [Notify] High-value order %s: $%.2f!\n",
					event.OrderID, event.Amount)
			}
		}
		fmt.Println("  [Notify] Subscriber shut down")
	}()

	// Subscriber 3: Inventory (different topic)
	inventoryCh := bus.Subscribe("inventory")
	wg.Add(1)
	go func() {
		defer wg.Done()
		for event := range inventoryCh {
			fmt.Printf("  [Inventory] Stock change for order %s\n", event.OrderID)
		}
		fmt.Println("  [Inventory] Subscriber shut down")
	}()

	// Publish events
	fmt.Println("\n--- Publishing events ---")
	events := []OrderEvent{
		{OrderID: "ORD-001", Action: "created", Amount: 49.99, Timestamp: time.Now()},
		{OrderID: "ORD-002", Action: "created", Amount: 250.00, Timestamp: time.Now()},
		{OrderID: "ORD-003", Action: "shipped", Amount: 75.00, Timestamp: time.Now()},
	}

	for _, e := range events {
		bus.Publish("orders", e)
		bus.Publish("inventory", e) // Also notify inventory
	}

	// Allow subscribers to process
	time.Sleep(100 * time.Millisecond)

	// Graceful shutdown
	fmt.Println("\n--- Shutting down ---")
	bus.Close()
	wg.Wait()

	// Key takeaways
	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. EventBus[T] is generic type-safe events without interface{}")
	fmt.Println("2. Non-blocking publish: slow subscribers don't block producers")
	fmt.Println("3. Topic-based routing: subscribers only get relevant events")
	fmt.Println("4. Close() propagates shutdown to all subscribers via channel close")
	fmt.Println("5. Use for in-process decoupling; for distributed, use Kafka/NATS")
}
