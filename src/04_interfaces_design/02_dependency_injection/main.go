/*
What this teaches:
    Interface-based dependency injection in Go. A NotificationService depends on
    an abstract Sender interface, allowing EmailSender, SMSSender, and MockSender
    to be injected at runtime or in tests.

Beginner analogy:
    "DI is like a power outlet any appliance (implementation) can plug in as
     long as it fits the socket (interface). The wall doesn't care if it's a lamp
     or a toaster."

C++ comparison:
    "Like the Dependency Inversion Principle (DIP) depend on abstractions, not
     concretions. In C++ you'd use abstract base classes with pure virtual methods.
     In Go, interfaces are implicitly satisfied no 'implements' keyword needed."

Interview relevance:
    DI is fundamental to testable Go code. Interviewers ask how you'd mock external
    services, structure constructor injection, and why Go's implicit interfaces make
    DI particularly elegant without DI frameworks.
*/

package main

import (
	"fmt"
	"strings"
)

// --- The abstraction: what we depend on ---

type Sender interface {
	Send(to, subject, body string) error
}

// --- Concrete implementation: Email ---

type EmailSender struct {
	SMTPHost string
	Port     int
}

func (e *EmailSender) Send(to, subject, body string) error {
	fmt.Printf("[EMAIL via %s:%d] To: %s | Subject: %s | Body: %s\n",
		e.SMTPHost, e.Port, to, subject, body)
	return nil
}

// --- Concrete implementation: SMS ---

type SMSSender struct {
	APIKey string
}

func (s *SMSSender) Send(to, subject, body string) error {
	msg := subject + ": " + body
	if len(msg) > 160 {
		msg = msg[:157] + "..."
	}
	fmt.Printf("[SMS via API] To: %s | Message: %s\n", to, msg)
	return nil
}

// --- Mock for testing ---

type MockSender struct {
	Calls []MockCall
}

type MockCall struct {
	To, Subject, Body string
}

func (m *MockSender) Send(to, subject, body string) error {
	m.Calls = append(m.Calls, MockCall{To: to, Subject: subject, Body: body})
	return nil
}

func (m *MockSender) AssertCalled(t string, times int) bool {
	count := 0
	for _, c := range m.Calls {
		if c.To == t {
			count++
		}
	}
	return count == times
}

// --- The service that depends on Sender interface ---

type NotificationService struct {
	sender Sender // Injected dependency unexported for encapsulation
}

// Constructor injection the idiomatic Go pattern
func NewNotificationService(s Sender) *NotificationService {
	return &NotificationService{sender: s}
}

func (ns *NotificationService) NotifyUser(userEmail, event string) error {
	subject := "Notification: " + strings.Title(event)
	body := fmt.Sprintf("Hello! Event '%s' occurred.", event)
	return ns.sender.Send(userEmail, subject, body)
}

func (ns *NotificationService) NotifyBatch(emails []string, event string) error {
	for _, email := range emails {
		if err := ns.NotifyUser(email, event); err != nil {
			return fmt.Errorf("failed to notify %s: %w", email, err)
		}
	}
	return nil
}

func main() {
	fmt.Println("=== Dependency Injection via Interfaces ===")

	// 1. Inject EmailSender
	fmt.Println("\n--- Using EmailSender ---")
	emailSvc := NewNotificationService(&EmailSender{
		SMTPHost: "smtp.example.com",
		Port:     587,
	})
	emailSvc.NotifyUser("alice@example.com", "signup")

	// 2. Inject SMSSender same service, different transport
	fmt.Println("\n--- Using SMSSender ---")
	smsSvc := NewNotificationService(&SMSSender{APIKey: "sk_live_xxx"})
	smsSvc.NotifyUser("+1-555-0123", "password_reset")

	// 3. Inject MockSender for testing
	fmt.Println("\n--- Using MockSender (testing) ---")
	mock := &MockSender{}
	testSvc := NewNotificationService(mock)
	testSvc.NotifyBatch([]string{"bob@test.com", "carol@test.com"}, "deployment")

	fmt.Printf("Mock received %d calls\n", len(mock.Calls))
	fmt.Printf("Bob notified: %v\n", mock.AssertCalled("bob@test.com", 1))
	fmt.Printf("Unknown notified: %v\n", mock.AssertCalled("unknown@test.com", 1))

	// 4. Key takeaways
	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. Depend on interfaces (Sender), not concrete types (EmailSender)")
	fmt.Println("2. Use constructor injection: NewXxx(dep Interface) *Xxx")
	fmt.Println("3. No DI framework needed Go's implicit interfaces suffice")
	fmt.Println("4. Mocks become trivial: implement the interface, record calls")
	fmt.Println("5. Swap implementations at runtime without changing service code")
}
