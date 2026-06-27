// Kafka Consumer subscribes to a topic and processes events with a consumer group.
//
// LEARNING NOTES:
// - Consumer groups: Kafka assigns partitions to consumers in a group automatically
// - Rebalancing: when a consumer joins/leaves, partitions are redistributed
// - Offset commit: tracks where each consumer left off (survives restarts)
// - "At-least-once" by default: commit after processing (may reprocess on crash)
// - "Exactly-once" requires idempotent writes or transactional processing
//
// Run: go run consumer/main.go
// Requires: Kafka broker on localhost:9092
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type PaymentEvent struct {
	PaymentID string  `json:"payment_id"`
	OrderID   string  `json:"order_id"`
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
	Timestamp int64   `json:"timestamp"`
}

func main() {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "payment-processor",  // Consumer group
		"auto.offset.reset": "earliest",           // Start from beginning if no offset stored
		"enable.auto.commit": true,                 // Auto-commit offsets every 5s
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// Subscribe to the "payments" topic
	err = consumer.SubscribeTopics([]string{"payments"}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	fmt.Println("Kafka consumer started. Waiting for messages...")

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case <-sigChan:
			run = false
		default:
			msg, err := consumer.ReadMessage(-1) // Block until message
			if err != nil {
				log.Printf("Consumer error: %v", err)
				continue
			}

			var event PaymentEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("Unmarshal error: %v", err)
				continue
			}

			// Process the event
			fmt.Printf("[Partition %d | Offset %d] Payment: %s | User: %s | $%.2f | %s\n",
				msg.TopicPartition.Partition,
				msg.TopicPartition.Offset,
				event.PaymentID,
				event.UserID,
				event.Amount,
				event.Status,
			)
		}
	}

	fmt.Println("Consumer shutting down...")
}
