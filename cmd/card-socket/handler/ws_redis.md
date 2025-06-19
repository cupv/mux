package handler

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

// WebSocketServer holds the state of the WebSocket server
type WebSocketServer struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	mutex      sync.RWMutex
	upgrader   websocket.Upgrader
	redisClient *redis.Client
	ctx        context.Context
}

// NewWebSocketServer initializes a new WebSocket server
func NewWebSocketServer(redisAddr string) *WebSocketServer {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	return &WebSocketServer{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		redisClient: rdb,
		ctx:         context.Background(),
	}
}

// Start initializes the broadcasting goroutines
func (server *WebSocketServer) Start() {
	// Start local broadcaster
	go server.localBroadcast()

	// Start Redis subscriber for distributed broadcast
	go server.redisSubscribe()
}

// localBroadcast handles messages from the broadcast channel and sends them to clients
func (server *WebSocketServer) localBroadcast() {
	for message := range server.broadcast {
		server.mutex.RLock()
		for client := range server.clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				client.Close()
				server.unregister(client)
			}
		}
		server.mutex.RUnlock()
	}
}

// redisSubscribe subscribes to the Redis "messages" channel for cross-instance messaging
func (server *WebSocketServer) redisSubscribe() {
	pubsub := server.redisClient.Subscribe(server.ctx, "messages")
	defer pubsub.Close()

	for msg := range pubsub.Channel() {
		server.mutex.RLock()
		for client := range server.clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
			if err != nil {
				client.Close()
				server.unregister(client)
			}
		}
		server.mutex.RUnlock()
	}
}

// HandleConnections upgrades HTTP requests to WebSocket connections
func (server *WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	server.register(conn)
	go server.handleMessages(conn)
}

// register adds a new client to the server
func (server *WebSocketServer) register(conn *websocket.Conn) {
	server.mutex.Lock()
	server.clients[conn] = true
	server.mutex.Unlock()
}

// unregister removes a client from the server
func (server *WebSocketServer) unregister(conn *websocket.Conn) {
	server.mutex.Lock()
	delete(server.clients, conn)
	server.mutex.Unlock()
	conn.Close()
}

// handleMessages reads messages from a client and publishes them to Redis for cross-instance broadcast
func (server *WebSocketServer) handleMessages(conn *websocket.Conn) {
	defer server.unregister(conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		// Publish message to Redis for cross-instance distribution
		err = server.redisClient.Publish(server.ctx, "messages", message).Err()
		if err != nil {
			log.Println("Redis publish error:", err)
		}
	}
}
