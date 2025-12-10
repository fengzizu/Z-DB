package tests

import (
	"bufio"
	"net"
	"testing"
	"time"
)

func TestServerConnection(t *testing.T) {
	// Attempt to connect. Retry a few times in case server is starting up.
	var conn net.Conn
	var err error
	for i := 0; i < 5; i++ {
		conn, err = net.Dial("tcp", "localhost:6379")
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Send PING as RESP Array: *1\r\n$4\r\nPING\r\n
	payload := []byte("*1\r\n$4\r\nPING\r\n")
	_, err = conn.Write(payload)
	if err != nil {
		t.Fatalf("Failed to write: %v", err)
	}

	// Read response (Expect +OK\r\n)
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read: %v", err)
	}

	expected := "+OK\r\n"
	if line != expected {
		t.Errorf("Expected %q, got %q", expected, line)
	}
}

