package db

import (
	"sync"
	"time"

	"MessasingApp/backend/models"
)

var (
	usersMu sync.RWMutex
	users   = map[string]string{} // username -> password (plaintext for demo only!)

	msgMu    sync.RWMutex
	messages = []models.Message{}
)

// CreateUser returns false if already exists
func CreateUser(username, password string) bool {
	usersMu.Lock()
	defer usersMu.Unlock()
	if _, ok := users[username]; ok {
		return false
	}
	users[username] = password
	return true
}

func ValidateUser(username, password string) bool {
	usersMu.RLock()
	defer usersMu.RUnlock()
	if pw, ok := users[username]; ok && pw == password {
		return true
	}
	return false
}

func AddMessage(m models.Message) {
	msgMu.Lock()
	defer msgMu.Unlock()
	messages = append(messages, m)
}

func GetAllMessages() []models.Message {
	msgMu.RLock()
	defer msgMu.RUnlock()
	// copy to avoid races
	result := make([]models.Message, len(messages))
	copy(result, messages)
	return result
}

// Factory for creating message IDs (simple timestamp)
func NewMessage(sender, receiver, content string) models.Message {
	return models.Message{
		ID:        time.Now().Format("20060102150405.000"),
		Sender:    sender,
		Receiver:  receiver,
		Content:   content,
		Timestamp: time.Now().Unix(),
	}
}
