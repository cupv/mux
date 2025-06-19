package main

import (
	"fmt"
	"net/http"
	"github.com/cupv/mux/cmd/card-socket/handler"
)

func main() {
	server := handler.NewWebSocketServer("localhost:6379","abcde12345-")
	server.Start()
	http.HandleFunc("/ws", server.HandleConnections)
	fmt.Println("WebSocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
