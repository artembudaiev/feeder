package message

import (
	"context"
	"time"
)

type Repository interface {
	Add(context.Context, Message) (string, error)
	GetAll(context.Context) ([]Message, error)
	Get(context.Context, string) (Message, error)
}

type Observer interface {
	Send(ctx context.Context, message Message)
	ID() string
}

type service struct {
	repository Repository
	observers  map[string]Observer
}

func (s *service) Add(ctx context.Context, message Message) error {
	_, err := s.repository.Add(ctx, message)
	if err != nil {
		return err
	}
	for i := range s.observers {
		sendCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		s.observers[i].Send(sendCtx, message)
	}
	return nil
}

func (s *service) Attach(ctx context.Context, observer Observer) error {
	messages, _ := s.repository.GetAll(ctx)
	s.observers[observer.ID()] = observer
	for _, message := range messages {
		observer.Send(ctx, message)
	}
	return nil
}

func (s *service) Detach(_ context.Context, id string) error {
	delete(s.observers, id)
	return nil
}
