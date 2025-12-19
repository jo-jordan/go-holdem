package entities

type ChatMessage struct {
	SenderID string `json:"sender_id"`
	Content  string `json:"content"`
}
