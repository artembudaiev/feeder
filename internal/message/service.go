package message

import (
	"context"
	"github.com/google/uuid"
	"log"
	"sync"
	"time"
)

type Repository interface {
	Add(context.Context, Message) (string, error)
	GetAll(context.Context) ([]Message, error)
}

type DLQProducer interface {
	Produce(ctx context.Context, msg []byte) error
}

type service struct {
	repository     Repository
	dlqProducer    DLQProducer
	sendChannels   map[string]chan Message
	receiveChannel chan Message
	m              sync.RWMutex
}

func NewService(repository Repository, dlqProducer DLQProducer) Service {
	return &service{
		repository:     repository,
		dlqProducer:    dlqProducer,
		sendChannels:   make(map[string]chan Message, 10),
		receiveChannel: make(chan Message, 2),
		m:              sync.RWMutex{},
	}
}

func (s *service) Add(ctx context.Context, message Message) error {
	select {
	case s.receiveChannel <- message:
	default:
		// todo: proper serialization, for now ok
		if err := s.dlqProducer.Produce(ctx, []byte(message.Text)); err != nil {
			log.Println(err)
		}
		log.Printf("message moved to dlq: %s", message.Text)
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
			// todo: remove the next line, it is used only for demonstration purpose to simulate high workload
			time.Sleep(time.Millisecond * 3000)

			// todo: consider where to unlock
			s.m.Lock()
			{
				_, err := s.repository.Add(ctx, msg)
				if err != nil {
					// todo: maybe move to dlq
					log.Println(err)
					continue
				}
				log.Printf("added message: %s", msg)
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
