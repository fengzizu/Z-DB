package server

import (
	"fmt"
	"io"
	"net"
	"zdb/redis/internal/resp"
)

func Run() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Printf("Failed to bind to port 6379: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Println("Listening on port 6379...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := resp.NewReader(conn)
	writer := resp.NewWriter(conn)

	for {
		value, err := reader.Read()
		if err != nil {
			if err != io.EOF {
				fmt.Println("Connection error:", err)
			}
			return
		}

		fmt.Printf("Received: %+v\n", value)

		// For Phase 1 validation: Respond with +OK
		if err := writer.Write(resp.Value{Type: "string", Str: "OK"}); err != nil {
			fmt.Println("Write error:", err)
			return
		}
	}
}
