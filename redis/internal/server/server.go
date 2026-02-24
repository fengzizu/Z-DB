package server

import (
	"fmt"
	"io"
	"net"
	"zdb/redis/internal/core"
	"zdb/redis/internal/resp"
)

// Run starts the TCP server on port 6379.
// Run 启动监听端口 6379 的 TCP 服务器。
func Run(store *core.Store) {
	// 1. Bind to the port / 绑定端口
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Printf("Failed to bind to port 6379: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Println("Listening on port 6379...")

	// 2. Accept loop / 接收连接循环
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		// 3. Spawn a Goroutine for each connection / 为每个连接启动一个 Goroutine
		go handleConnection(conn, store)
	}
}

// handleConnection manages the lifecycle of a single client connection.
// handleConnection 管理单个客户端连接的生命周期。
func handleConnection(conn net.Conn, store *core.Store) {
	defer conn.Close()

	// Initialize RESP Reader and Writer for this connection
	// 为该连接初始化 RESP 读取器和写入器
	reader := resp.NewReader(conn)
	writer := resp.NewWriter(conn)

	// Command processing loop / 命令处理循环
	for {
		// 1. Read command / 读取命令
		value, err := reader.Read()
		if err != nil {
			if err != io.EOF {
				fmt.Println("Connection error:", err)
			}
			return
		}

		fmt.Printf("Received: %+v\n", value)

		// 2. Execute command (Placeholder for Phase 2) / 执行命令 (Phase 2 的占位符)
		// Currently just responds with OK to keep connection alive / 目前仅回复 OK 以保持连接
		result := core.EvalCommand(value.Array, store)

		if err := writer.Write(*result); err != nil {
			fmt.Println("Write error:", err)
			return
		}
	}
}
