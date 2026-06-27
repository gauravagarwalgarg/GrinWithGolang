/*
Module 10: Distributed - Event Sourcing Pattern

Demonstrates:
  - Event store: append-only log of domain events
  - Aggregate reconstruction from event history
  - Commands → Events → State (CQRS-lite)
  - Event types with metadata (timestamp, version)
  - Simple bank account example: deposit, withdraw, transfer
  - Snapshot optimization concept

Key insight: Instead of storing current state (CRUD), store the sequence
of events that led to current state. Benefits: complete audit trail,
temporal queries ("what was balance at time T?"), event replay for debugging.
Tradeoff: complexity + eventual consistency in read models.

Run: go run main.go
*/
package main

import (
	"fmt"
	"log"
	"time"
)

// --- Event types ---

type EventType string

const (
	AccountCreated  EventType = "AccountCreated"
	MoneyDeposited  EventType = "MoneyDeposited"
	MoneyWithdrawn  EventType = "MoneyWithdrawn"
	TransferSent    EventType = "TransferSent"
	TransferReceived EventType = "TransferReceived"
)

// Event represents a domain event in the append-only log.
type Event struct {
	ID        int
	Type      EventType
	Aggregate string // aggregate ID (e.g., account ID)
	Data      map[string]interface{}
	Timestamp time.Time
	Version   int
}

// --- Event Store (append-only log) ---

type EventStore struct {
	events []Event
	nextID int
}

func NewEventStore() *EventStore {
	return &EventStore{}
}

func (es *EventStore) Append(evt Event) {
	es.nextID++
	evt.ID = es.nextID
	evt.Timestamp = time.Now()
	es.events = append(es.events, evt)
}

func (es *EventStore) GetEvents(aggregateID string) []Event {
	var result []Event
	for _, evt := range es.events {
		if evt.Aggregate == aggregateID {
			result = append(result, evt)
		}
	}
	return result
}

// --- Aggregate: Bank Account (reconstructed from events) ---

type BankAccount struct {
	ID      string
	Owner   string
	Balance float64
	Version int
}

// Apply processes an event and updates account state.
func (a *BankAccount) Apply(evt Event) {
	switch evt.Type {
	case AccountCreated:
		a.ID = evt.Aggregate
		a.Owner = evt.Data["owner"].(string)
		a.Balance = 0
	case MoneyDeposited:
		a.Balance += evt.Data["amount"].(float64)
	case MoneyWithdrawn:
		a.Balance -= evt.Data["amount"].(float64)
	case TransferSent:
		a.Balance -= evt.Data["amount"].(float64)
	case TransferReceived:
		a.Balance += evt.Data["amount"].(float64)
	}
	a.Version = evt.Version
}

// Reconstruct rebuilds account state from event history.
func Reconstruct(store *EventStore, accountID string) *BankAccount {
	account := &BankAccount{}
	for _, evt := range store.GetEvents(accountID) {
		account.Apply(evt)
	}
	return account
}

// --- Command handlers (validate + produce events) ---

func CreateAccount(store *EventStore, id, owner string) {
	store.Append(Event{
		Type:      AccountCreated,
		Aggregate: id,
		Data:      map[string]interface{}{"owner": owner},
		Version:   1,
	})
}

func Deposit(store *EventStore, id string, amount float64) {
	account := Reconstruct(store, id)
	store.Append(Event{
		Type:      MoneyDeposited,
		Aggregate: id,
		Data:      map[string]interface{}{"amount": amount},
		Version:   account.Version + 1,
	})
}

func Withdraw(store *EventStore, id string, amount float64) error {
	account := Reconstruct(store, id)
	if account.Balance < amount {
		return fmt.Errorf("insufficient funds: balance=%.2f, requested=%.2f",
			account.Balance, amount)
	}
	store.Append(Event{
		Type:      MoneyWithdrawn,
		Aggregate: id,
		Data:      map[string]interface{}{"amount": amount},
		Version:   account.Version + 1,
	})
	return nil
}

func Transfer(store *EventStore, from, to string, amount float64) error {
	sender := Reconstruct(store, from)
	if sender.Balance < amount {
		return fmt.Errorf("insufficient funds for transfer")
	}
	store.Append(Event{
		Type:      TransferSent,
		Aggregate: from,
		Data:      map[string]interface{}{"amount": amount, "to": to},
		Version:   sender.Version + 1,
	})
	receiver := Reconstruct(store, to)
	store.Append(Event{
		Type:      TransferReceived,
		Aggregate: to,
		Data:      map[string]interface{}{"amount": amount, "from": from},
		Version:   receiver.Version + 1,
	})
	return nil
}

func main() {
	store := NewEventStore()

	// Execute commands (produce events)
	CreateAccount(store, "acc-001", "Alice")
	CreateAccount(store, "acc-002", "Bob")
	Deposit(store, "acc-001", 1000)
	Deposit(store, "acc-002", 500)
	Withdraw(store, "acc-001", 200)
	Transfer(store, "acc-001", "acc-002", 300)

	// Reconstruct state from events
	fmt.Println("=== Account States (reconstructed from events) ===")
	alice := Reconstruct(store, "acc-001")
	bob := Reconstruct(store, "acc-002")
	fmt.Printf("  Alice: balance=%.2f (version %d)\n", alice.Balance, alice.Version)
	fmt.Printf("  Bob:   balance=%.2f (version %d)\n", bob.Balance, bob.Version)

	// Show event log (audit trail)
	fmt.Println("\n=== Event Log (append-only) ===")
	for _, evt := range store.GetEvents("acc-001") {
		fmt.Printf("  #%d [%s] %s %v\n", evt.ID, evt.Type, evt.Aggregate, evt.Data)
	}

	// Demonstrate validation
	fmt.Println("\n=== Validation ===")
	if err := Withdraw(store, "acc-001", 9999); err != nil {
		log.Printf("Rejected: %v", err)
	}
}
