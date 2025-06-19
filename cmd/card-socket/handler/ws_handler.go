package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

// WebSocketServer holds the state of the WebSocket server
type WebSocketServer struct {
	clients     map[int]*websocket.Conn
	broadcast   chan []byte
	mutex       sync.RWMutex
	upgrader    websocket.Upgrader
	redisClient *redis.Client
	ctx         context.Context
}

// NewWebSocketServer initializes a new WebSocket server
func NewWebSocketServer(redisAddr string, redisPassword string) *WebSocketServer {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
	})
	return &WebSocketServer{
		clients:   make(map[int]*websocket.Conn),
		broadcast: make(chan []byte),
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
		for userId := range server.clients {
			client := server.clients[userId]
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				client.Close()
				server.unregister(client, userId)
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
		for userId := range server.clients {
			client := server.clients[userId]
			err := client.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
			if err != nil {
				client.Close()
				server.unregister(client, userId)
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
	xUserId := r.Header.Get("X-User-Id")
	fmt.Println("Client " + xUserId + " connect to ws server.")

	userId, _ := strconv.Atoi(xUserId)
	server.register(conn, userId)
	go server.handleMessages(conn, userId)
}

// register adds a new client to the server
func (server *WebSocketServer) register(conn *websocket.Conn, userId int) {
	server.mutex.Lock()
	server.clients[userId] = conn
	server.mutex.Unlock()
}

// unregister removes a client from the server
func (server *WebSocketServer) unregister(conn *websocket.Conn, userId int) {
	server.mutex.Lock()
	delete(server.clients, userId)
	server.mutex.Unlock()
	conn.Close()
}

// handleMessages reads messages from a client and publishes them to Redis for cross-instance broadcast
func (server *WebSocketServer) handleMessages(conn *websocket.Conn, userId int) {
	defer server.unregister(conn, userId)

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
