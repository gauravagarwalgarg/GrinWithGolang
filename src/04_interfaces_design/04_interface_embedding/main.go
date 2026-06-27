/*
What this teaches:
    Composing interfaces by embedding smaller interfaces into larger ones. Shows
    how io.ReadWriteCloser is built from Reader + Writer + Closer, how to create
    custom composed interfaces, and the rules around method sets with value vs
    pointer receivers.

Beginner analogy:
    "Interface embedding is like stacking job descriptions: a 'Senior Dev' role
     includes all duties of 'Dev' plus extra. The composed interface includes all
     methods of its embedded interfaces."

C++ comparison:
    "Similar to multiple inheritance of pure abstract classes, but without the
     diamond problem. Go's composition is flat no virtual base class gymnastics.
     Method sets replace C++'s const/non-const overload resolution."

Interview relevance:
    Interviewers test knowledge of method sets: why a value receiver satisfies both
    T and *T method sets, but a pointer receiver only satisfies *T. Understanding
    this explains compile errors when assigning to interfaces.
*/

package main

import (
	"fmt"
	"io"
	"strings"
)

// --- Building blocks: small interfaces ---

type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type Closer interface {
	Close() error
}

// --- Composed interfaces via embedding ---

type ReadCloser interface {
	Reader
	Closer
}

type ReadWriteCloser interface {
	Reader
	Writer
	Closer
}

// --- Custom domain interfaces ---

type Flusher interface {
	Flush() error
}

type BufferedWriter interface {
	Writer
	Flusher
}

// --- Implementation: DeviceStream ---

type DeviceStream struct {
	name   string
	data   []byte
	closed bool
}

func NewDeviceStream(name string) *DeviceStream {
	return &DeviceStream{name: name}
}

func (d *DeviceStream) Read(p []byte) (int, error) {
	if d.closed {
		return 0, fmt.Errorf("stream %s is closed", d.name)
	}
	n := copy(p, d.data)
	d.data = d.data[n:]
	return n, nil
}

func (d *DeviceStream) Write(p []byte) (int, error) {
	if d.closed {
		return 0, fmt.Errorf("stream %s is closed", d.name)
	}
	d.data = append(d.data, p...)
	return len(p), nil
}

func (d *DeviceStream) Close() error {
	fmt.Printf("  Closing stream: %s\n", d.name)
	d.closed = true
	return nil
}

// --- Method set demonstration ---

type Counter struct {
	count int
}

// Value receiver in method set of both Counter and *Counter
func (c Counter) Value() int {
	return c.count
}

// Pointer receiver ONLY in method set of *Counter
func (c *Counter) Increment() {
	c.count++
}

type ValueGetter interface {
	Value() int
}

type Incrementer interface {
	Increment()
}

type FullCounter interface {
	ValueGetter
	Incrementer
}

func main() {
	fmt.Println("=== Interface Embedding & Method Sets ===")

	// 1. Composed interfaces in action
	fmt.Println("\n--- Composed Interfaces ---")
	stream := NewDeviceStream("sensor-01")

	// DeviceStream satisfies ReadWriteCloser
	var rwc ReadWriteCloser = stream
	rwc.Write([]byte("temperature=23.5"))
	buf := make([]byte, 32)
	n, _ := rwc.Read(buf)
	fmt.Printf("  Read: %s\n", buf[:n])
	rwc.Close()

	// 2. Standard library example: io.ReadCloser
	fmt.Println("\n--- stdlib io.ReadCloser ---")
	rc := io.NopCloser(strings.NewReader("hello from io"))
	data := make([]byte, 20)
	n, _ = rc.Read(data)
	fmt.Printf("  io.ReadCloser read: %s\n", data[:n])
	rc.Close()

	// 3. Method sets: value vs pointer receiver
	fmt.Println("\n--- Method Sets ---")

	// Value type satisfies ValueGetter (value receiver)
	var vg ValueGetter = Counter{count: 10}
	fmt.Printf("  Counter value: %d\n", vg.Value())

	// Value type does NOT satisfy Incrementer (pointer receiver)
	// var inc Incrementer = Counter{count: 5} // COMPILE ERROR
	var inc Incrementer = &Counter{count: 5} // Pointer works
	inc.Increment()
	fmt.Printf("  After increment: %d\n", inc.(*Counter).count)

	// Pointer satisfies FullCounter (both value + pointer receivers)
	var fc FullCounter = &Counter{count: 100}
	fc.Increment()
	fmt.Printf("  FullCounter value: %d\n", fc.Value())

	// 4. Key takeaways
	fmt.Println("\n--- Key Takeaways ---")
	fmt.Println("1. Embed interfaces to compose: ReadWriteCloser = Reader + Writer + Closer")
	fmt.Println("2. Value receiver → method set of T and *T")
	fmt.Println("3. Pointer receiver → method set of *T only")
	fmt.Println("4. This is why you can't assign T to interface requiring pointer methods")
	fmt.Println("5. Keep interfaces small, compose for complex contracts")
}
