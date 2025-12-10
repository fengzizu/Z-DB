package main

import (
	"fmt"
	"zdb/redis/internal/server"
)

func main() {
	fmt.Println("Starting Z-DB Redis Core...")
	server.Run()
}
