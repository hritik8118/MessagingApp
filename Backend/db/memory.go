package db

import (
	"MessasingApp/Backend/models"
	"sync"
)

var (
	users    = make(map[string]models.User)
	messages = []models.Message{}
	mu       sync.RWMutex
)

func UserExists(username string) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, exists := users[username]
	return exists
}

func AddUser(user models.User) {
	mu.Lock()
	defer mu.Unlock()
	users[user.Username] = user
}

func GetUser(username string) (models.User, bool) {
	mu.RLock()
	defer mu.RUnlock()
	user, exists := users[username]
	return user, exists
}

func AddMessage(msg models.Message) {
	mu.Lock()
	defer mu.Unlock()
	messages = append(messages, msg)
}

func GetMessages() []models.Message {
	mu.RLock()
	defer mu.RUnlock()
	return messages
}
