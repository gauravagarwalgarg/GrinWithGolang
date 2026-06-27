// NSQ Producer publishes messages to an NSQ topic.
//
// LEARNING NOTES:
// - NSQ is a real-time distributed messaging platform (created at Bitly)
// - Compared to Kafka: simpler ops, no ZooKeeper, but less ordering guarantees
// - Architecture: nsqd (message daemon) + nsqlookupd (discovery) + nsqadmin (UI)
// - Messages are "at least once" delivery consumers must be idempotent
//
// For C++ devs: Think of NSQ as a lock-free MPMC queue distributed across machines.
//
// Run: go run producer/main.go
// Requires: nsqd running on localhost:4150
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nsqio/go-nsq"
)

// OrderEvent represents a domain event published when an order is placed.
type OrderEvent struct {
	OrderID   string    `json:"order_id"`
	UserID    string    `json:"user_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	// Create NSQ producer (connects to nsqd)
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer("localhost:4150", config)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Stop()

	// Publish 10 order events
	for i := 1; i <= 10; i++ {
		event := OrderEvent{
			OrderID:   fmt.Sprintf("order_%d", i),
			UserID:    fmt.Sprintf("user_%d", i%3+1),
			Amount:    float64(i) * 29.99,
			CreatedAt: time.Now(),
		}

		body, err := json.Marshal(event)
		if err != nil {
			log.Fatal(err)
		}

		// Publish to "orders" topic
		err = producer.Publish("orders", body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Published: %s ($%.2f)\n", event.OrderID, event.Amount)
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("\nAll messages published!")
}
