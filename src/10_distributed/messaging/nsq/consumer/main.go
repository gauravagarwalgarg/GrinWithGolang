// NSQ Consumer subscribes to a topic and processes messages.
//
// LEARNING NOTES:
// - A consumer connects to nsqlookupd to discover which nsqd instances have the topic
// - "Channel" = logical consumer group (like Kafka consumer groups)
// - Multiple consumers on the same channel = load-balanced (round-robin)
// - Multiple consumers on different channels = fan-out (each gets all messages)
// - msg.Finish() = ACK; msg.Requeue() = NACK (retry later)
//
// Run: go run consumer/main.go
// Requires: nsqlookupd on localhost:4161, nsqd on localhost:4150
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nsqio/go-nsq"
)

type OrderEvent struct {
	OrderID   string    `json:"order_id"`
	UserID    string    `json:"user_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

// OrderHandler implements nsq.Handler interface.
type OrderHandler struct{}

// HandleMessage is called for each message received.
func (h *OrderHandler) HandleMessage(msg *nsq.Message) error {
	var event OrderEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("Error unmarshaling: %v", err)
		return err // Returning error triggers requeue
	}

	// Process the order (in reality: update DB, send email, etc.)
	fmt.Printf("Processing: %s | User: %s | Amount: $%.2f\n",
		event.OrderID, event.UserID, event.Amount)

	// msg.Finish() is called automatically when handler returns nil
	return nil
}

func main() {
	config := nsq.NewConfig()
	config.MaxInFlight = 5 // Process up to 5 messages concurrently

	// Create consumer for topic "orders", channel "order-processor"
	consumer, err := nsq.NewConsumer("orders", "order-processor", config)
	if err != nil {
		log.Fatal(err)
	}

	// Set the handler (can also use AddConcurrentHandlers for parallelism)
	consumer.AddHandler(&OrderHandler{})

	// Connect to nsqlookupd for topic discovery
	err = consumer.ConnectToNSQLookupd("localhost:4161")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Consumer started. Waiting for messages...")

	// Wait for interrupt signal for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	consumer.Stop()
	fmt.Println("Consumer stopped.")
}
