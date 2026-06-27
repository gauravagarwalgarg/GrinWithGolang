// Kafka Producer publishes events to a Kafka topic with partitioning.
//
// LEARNING NOTES:
// - Kafka = distributed commit log. Messages are ordered within a partition.
// - Partitions enable parallelism: more partitions = more consumers in a group
// - Key-based partitioning: same key → same partition → guaranteed ordering for that key
// - Kafka retains messages (default 7 days) consumers can replay from any offset
// - vs NSQ: Kafka has stronger ordering, replay ability, but heavier ops
//
// For C++ devs: Kafka is like a distributed ring buffer with persistent storage
// and consumer offset tracking (similar to shared memory segments with seek).
//
// Run: go run producer/main.go
// Requires: Kafka broker on localhost:9092
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type PaymentEvent struct {
	PaymentID string  `json:"payment_id"`
	OrderID   string  `json:"order_id"`
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"` // pending, completed, failed
	Timestamp int64   `json:"timestamp"`
}

func main() {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"acks":              "all",  // Wait for all replicas (strongest durability)
		"retries":           5,      // Retry on transient failures
		"linger.ms":         10,     // Batch messages for 10ms (throughput vs latency tradeoff)
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	topic := "payments"

	// Delivery report goroutine (async confirmation)
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("FAILED delivery: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered to partition %d at offset %v\n",
						ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			}
		}
	}()

	// Publish 10 payment events
	for i := 1; i <= 10; i++ {
		event := PaymentEvent{
			PaymentID: fmt.Sprintf("pay_%d", i),
			OrderID:   fmt.Sprintf("order_%d", i),
			UserID:    fmt.Sprintf("user_%d", i%3+1),
			Amount:    float64(i) * 49.99,
			Status:    "completed",
			Timestamp: time.Now().UnixMilli(),
		}

		body, _ := json.Marshal(event)

		// Key = UserID → all events for the same user go to the same partition (ordered)
		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            []byte(event.UserID),
			Value:          body,
		}, nil)

		if err != nil {
			log.Printf("Produce error: %v", err)
		}

		fmt.Printf("Sent: %s ($%.2f)\n", event.PaymentID, event.Amount)
		time.Sleep(200 * time.Millisecond)
	}

	// Wait for all messages to be delivered
	producer.Flush(5000)
	fmt.Println("\nAll payment events published!")
}
