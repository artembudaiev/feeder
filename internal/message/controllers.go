package message

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

type Service interface {
	Add(ctx context.Context, message Message) error
	Attach(ctx context.Context, observer Observer) error
	Detach(ctx context.Context, observer Observer) error
}

type Controller struct {
	service Service
}

func (c *Controller) HandleAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// todo: parse, call c.service.Add(r.Context(),message)
		http.NotFound(w, r)
	}
}

func (c *Controller) HandleGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.NotFound(w, r)
			return
		}
		messageChan := make(chan Message, 10)
		observer := &observer{id: uuid.New().String(), messageChan: messageChan}
		c.service.Attach(ctx, observer)
		defer c.service.Detach(ctx, observer)
		for {
			select {
			case <-messageChan:
				// todo: write msg
				flusher.Flush()
			case <-ctx.Done():
				// finish
				return
			}
		}

	}
}

type observer struct {
	id          string
	messageChan chan Message
}

func (o *observer) Send(ctx context.Context, message Message) {
	select {
	case o.messageChan <- message:

	case <-ctx.Done():

	}

}

func (o *observer) ID() string {
	return o.id
}
