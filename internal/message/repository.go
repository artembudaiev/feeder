package message

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

type InMemoryRepository struct {
	storage map[string]Message
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{storage: make(map[string]Message, 10)}
}

func (r *InMemoryRepository) Add(_ context.Context, message Message) (string, error) {
	id := uuid.New().String()
	r.storage[id] = message
	return id, nil
}

func (r *InMemoryRepository) GetAll(_ context.Context) ([]Message, error) {
	messages := make([]Message, len(r.storage))
	for _, message := range r.storage {
		messages = append(messages, message)
	}
	return messages, nil
}

func (r *InMemoryRepository) Get(_ context.Context, id string) (Message, error) {
	message, ok := r.storage[id]
	if !ok {
		return "", fmt.Errorf("no message found by id %s", id)
	}
	return message, nil
}
