package message

import (
	"context"
	"github.com/google/uuid"
	"sync"
)

type Repository interface {
	Add(context.Context, Message) (string, error)
	GetAll(context.Context) ([]Message, error)
	Get(context.Context, string) (Message, error)
}

type service struct {
	repository      Repository
	messageChannels map[string]chan Message
	m               sync.RWMutex
}

func (s *service) Add(ctx context.Context, message Message) error {
	s.m.Lock()
	defer s.m.Unlock()

	_, err := s.repository.Add(ctx, message)
	if err != nil {
		return err
	}
	for id, messageChannel := range s.messageChannels {
		select {
		case messageChannel <- message:
		default:
			// message channel is full - consider that receiver does not cope, drop it not to block or not to produce infinite number of goroutines
			close(messageChannel)
			delete(s.messageChannels, id)
		}
	}
	return nil
}

func (s *service) Attach(ctx context.Context) (<-chan Message, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	// get all existing messages from repo and feed with them channel
	messages, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	messageChan := make(chan Message)
	for _, message := range messages {
		messageChan <- message
	}
	s.messageChannels[uuid.New().String()] = messageChan
	return messageChan, nil
}
