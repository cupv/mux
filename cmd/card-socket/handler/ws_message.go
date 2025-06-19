package handler

// Message represents the message structure for sending and receiving
type Message struct {
	RecipientID string `json:"recipient_id"`
	SenderID    string `json:"sender_id"`
	Content     string `json:"content"`
}