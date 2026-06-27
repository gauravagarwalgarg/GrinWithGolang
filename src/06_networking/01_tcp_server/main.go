/*
Module 6: Networking - TCP Echo Server

Demonstrates:
  - net.Listen for TCP socket binding
  - Accept loop handling multiple clients
  - Goroutine-per-connection model
  - Graceful shutdown with context cancellation
  - RAII-style cleanup with defer conn.Close()
  - Go's net package simplicity vs C/C++ socket programming

Key insight: In C++ you'd need socket(), bind(), listen(), accept() with
manual fd management. Go wraps this in two lines: net.Listen + Accept.

Run: go run main.go
Test: echo "hello" | nc localhost 9000
*/
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Trap SIGINT/SIGTERM for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()
	log.Println("TCP echo server listening on :9000")

	var wg sync.WaitGroup

	// Goroutine to close listener on shutdown signal
	go func() {
		select {
		case sig := <-sigCh:
			log.Printf("received signal %v, shutting down...", sig)
			cancel()
			listener.Close() // unblocks Accept()
		case <-ctx.Done():
		}
	}()

	// Accept loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				log.Println("server shutting down, waiting for connections...")
				wg.Wait()
				log.Println("all connections closed, goodbye")
				return
			default:
				log.Printf("accept error: %v", err)
				continue
			}
		}

		wg.Add(1)
		go handleConnection(ctx, conn, &wg)
	}
}

// handleConnection echoes data back to the client.
// Demonstrates RAII-style resource management with defer.
func handleConnection(ctx context.Context, conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close() // RAII: guaranteed cleanup regardless of exit path

	remoteAddr := conn.RemoteAddr().String()
	log.Printf("client connected: %s", remoteAddr)

	// Set a deadline that respects context cancellation
	go func() {
		<-ctx.Done()
		conn.SetDeadline(time.Now()) // unblocks any Read/Write
	}()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("read error from %s: %v", remoteAddr, err)
			}
			break
		}

		// Echo back with a prefix
		response := fmt.Sprintf("[echo] %s", buf[:n])
		if _, err := conn.Write([]byte(response)); err != nil {
			log.Printf("write error to %s: %v", remoteAddr, err)
			break
		}
	}

	log.Printf("client disconnected: %s", remoteAddr)
}
