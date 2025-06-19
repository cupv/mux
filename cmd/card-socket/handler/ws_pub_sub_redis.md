package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9" // Redis client for Pub/Sub
)

// WebSocketServer holds the state of the WebSocket server
type WebSocketServer struct {
	clients   map[string]*websocket.Conn
	mutex     sync.Mutex
	upgrader  websocket.Upgrader
	redis     *redis.Client
	ctx       context.Context
}

// Message represents the structure for sending and receiving messages
type Message struct {
	RecipientID string `json:"recipient_id"`
	SenderID    string `json:"sender_id"`
	Content     string `json:"content"`
}

// NewWebSocketServer initializes a new WebSocket server with Redis client
func NewWebSocketServer(redisAddr string) *WebSocketServer {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	ctx := context.Background()

	server := &WebSocketServer{
		clients: make(map[string]*websocket.Conn),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		redis: rdb,
		ctx:   ctx,
	}

	// Subscribe to Redis channel for incoming messages
	go server.listenRedis()

	return server
}

// HandleConnections handles new WebSocket connections
func (server *WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		conn.Close()
		return
	}

	server.registerClient(userID, conn)
	go server.handleMessages(userID, conn)
}

// registerClient and unregisterClient manage client connections
func (server *WebSocketServer) registerClient(userID string, conn *websocket.Conn) {
	server.mutex.Lock()
	server.clients[userID] = conn
	server.mutex.Unlock()
}

func (server *WebSocketServer) unregisterClient(userID string) {
	server.mutex.Lock()
	if conn, exists := server.clients[userID]; exists {
		conn.Close()
		delete(server.clients, userID)
	}
	server.mutex.Unlock()
}

// handleMessages reads messages from a client and publishes them to Redis
func (server *WebSocketServer) handleMessages(userID string, conn *websocket.Conn) {
	defer server.unregisterClient(userID)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			fmt.Println("Invalid message format:", err)
			continue
		}

		msg.SenderID = userID
		server.publishToRedis(msg)
	}
}

// publishToRedis publishes a message to Redis
func (server *WebSocketServer) publishToRedis(msg Message) {
	messageData, _ := json.Marshal(msg)
	server.redis.Publish(server.ctx, "chat", messageData)
}

// listenRedis listens for messages on Redis and broadcasts them to local clients
func (server *WebSocketServer) listenRedis() {
	pubsub := server.redis.Subscribe(server.ctx, "chat")
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(server.ctx)
		if err != nil {
			fmt.Println("Error receiving message:", err)
			continue
		}

		var message Message
		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			fmt.Println("Invalid message format:", err)
			continue
		}

		// Send message to the recipient if they are connected to this instance
		server.sendToClient(message)
	}
}

// sendToClient sends a message to a specific client
func (server *WebSocketServer) sendToClient(message Message) {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	if client, exists := server.clients[message.RecipientID]; exists {
		err := client.WriteJSON(message)
		if err != nil {
			client.Close()
			delete(server.clients, message.RecipientID)
		}
	}
}
