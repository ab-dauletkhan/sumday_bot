package repository

import (
	"sync"

	"github.com/ab-dauletkhan/sumday_bot/internal/models"
)

type MapRepo struct {
	messages map[int64][]models.Message
	mu       sync.RWMutex
}

func NewMapRepo() *MapRepo {
	return &MapRepo{
		messages: make(map[int64][]models.Message),
	}
}

func (r *MapRepo) Save(userID int64, message models.Message) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.messages[userID] = append(r.messages[userID], message)
}

func (r *MapRepo) GetMessages(userID int64) []models.Message {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.messages[userID]
}
