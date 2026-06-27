/*
Module 6: Networking - WebSocket Basics (Standard Library Only)

Demonstrates:
  - WebSocket handshake using net/http + crypto/sha1
  - The WebSocket upgrade protocol (HTTP → persistent connection)
  - Bidirectional communication pattern
  - Frame reading (simplified text frames)
  - Why WebSocket exists: full-duplex over single TCP connection

Key insight: WebSocket is just HTTP upgrade + a thin framing protocol.
The handshake is: client sends Sec-WebSocket-Key, server responds with
SHA-1(key + magic GUID) base64-encoded. Then it's raw framed TCP.

Run: go run main.go
Test: Use a WebSocket client to connect to ws://localhost:8081/ws
*/
package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

const websocketMagicGUID = "258EAFA5-E914-47DA-95CA-5AB5DC085B11"

// computeAcceptKey generates the Sec-WebSocket-Accept value per RFC 6455.
func computeAcceptKey(clientKey string) string {
	h := sha1.New()
	h.Write([]byte(strings.TrimSpace(clientKey) + websocketMagicGUID))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// upgradeToWebSocket hijacks the HTTP connection and performs WS handshake.
func upgradeToWebSocket(w http.ResponseWriter, r *http.Request) (net.Conn, error) {
	clientKey := r.Header.Get("Sec-WebSocket-Key")
	if clientKey == "" {
		http.Error(w, "missing Sec-WebSocket-Key", http.StatusBadRequest)
		return nil, fmt.Errorf("missing websocket key")
	}

	// Hijack the connection from net/http
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		return nil, fmt.Errorf("server doesn't support hijacking")
	}
	conn, buf, err := hijacker.Hijack()
	if err != nil {
		return nil, err
	}

	// Send upgrade response
	acceptKey := computeAcceptKey(clientKey)
	response := "HTTP/1.1 101 Switching Protocols\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Accept: " + acceptKey + "\r\n\r\n"
	buf.WriteString(response)
	buf.Flush()

	return conn, nil
}

// readFrame reads a simplified WebSocket text frame (no masking for clarity).
func readFrame(conn net.Conn) (string, error) {
	header := make([]byte, 2)
	if _, err := io.ReadFull(conn, header); err != nil {
		return "", err
	}

	masked := (header[1] & 0x80) != 0
	length := int(header[1] & 0x7F)

	// Extended payload lengths
	if length == 126 {
		ext := make([]byte, 2)
		io.ReadFull(conn, ext)
		length = int(binary.BigEndian.Uint16(ext))
	}

	var maskKey [4]byte
	if masked {
		io.ReadFull(conn, maskKey[:])
	}

	payload := make([]byte, length)
	io.ReadFull(conn, payload)

	if masked {
		for i := range payload {
			payload[i] ^= maskKey[i%4]
		}
	}
	return string(payload), nil
}

// writeFrame sends a WebSocket text frame (server frames are unmasked).
func writeFrame(conn net.Conn, msg string) error {
	frame := []byte{0x81} // FIN + text opcode
	payload := []byte(msg)
	if len(payload) < 126 {
		frame = append(frame, byte(len(payload)))
	}
	frame = append(frame, payload...)
	_, err := conn.Write(frame)
	return err
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgradeToWebSocket(w, r)
	if err != nil {
		log.Printf("upgrade failed: %v", err)
		return
	}
	defer conn.Close()
	log.Printf("WebSocket connected: %s", conn.RemoteAddr())

	for {
		msg, err := readFrame(conn)
		if err != nil {
			log.Printf("read error: %v", err)
			return
		}
		log.Printf("received: %s", msg)
		echo := fmt.Sprintf("[echo] %s", msg)
		if err := writeFrame(conn, echo); err != nil {
			return
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleWS)
	log.Println("WebSocket server on :8081/ws")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
