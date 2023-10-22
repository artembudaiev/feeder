package message

import (
	"context"
	"github.com/google/uuid"
	"log"
	"sync"
)

type Repository interface {
	Add(context.Context, Message) (string, error)
	GetAll(context.Context) ([]Message, error)
	Get(context.Context, string) (Message, error)
}

type service struct {
	repository     Repository
	sendChannels   map[string]chan Message
	receiveChannel chan Message
	m              sync.RWMutex
}

func NewService(repository Repository) Service {
	return &service{
		repository:     repository,
		sendChannels:   make(map[string]chan Message, 10),
		receiveChannel: make(chan Message, 2),
		m:              sync.RWMutex{},
	}
}

func (s *service) Add(ctx context.Context, message Message) error {

	select {
	case s.receiveChannel <- message:
	default:
		// todo: move to dlq
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
	// todo: maybe return existing messages as slice and use chan only for new
	messageChan := make(chan Message, len(messages)*2)
	for _, message := range messages {
		messageChan <- message
	}
	s.sendChannels[uuid.New().String()] = messageChan
	return messageChan, nil
}

func (s *service) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-s.receiveChannel:
			// todo: consider where to unlock
			s.m.Lock()
			{
				_, err := s.repository.Add(ctx, msg)
				if err != nil {
					// todo: maybe move to dlq
					log.Println(err)
					continue
				}
				// broadcast new message to all channels
				for id, messageChannel := range s.sendChannels {
					select {
					case messageChannel <- msg:
					default:
						// message channel is full - consider that receiver does not cope, drop it not to block or not to produce infinite number of goroutines
						close(messageChannel)
						delete(s.sendChannels, id)
					}
				}
			}
			s.m.Unlock()

		}
	}
}
