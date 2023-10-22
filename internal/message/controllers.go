package message

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type Service interface {
	Start(ctx context.Context) error
	Add(ctx context.Context, message Message) error
	Attach(ctx context.Context) (<-chan Message, error)
}

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) HandleAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var msg Message
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&msg); err != nil {
			log.Println(err.Error())
			http.Error(w, "invalid JSON payload", http.StatusBadRequest)
			return
		}
		if err := c.service.Add(r.Context(), msg); err != nil {
			log.Println(err.Error())
			http.Error(w, "failed to add message", http.StatusInternalServerError)
			return
		}
		log.Printf("added message: %s", msg)
		w.WriteHeader(http.StatusOK)
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

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Transfer-Encoding", "chunked")

		flusher.Flush()

		enc := json.NewEncoder(w)

		msgChan, err := c.service.Attach(ctx)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "failed to listen to the new messages", http.StatusInternalServerError)
			return
		}
		for {
			select {
			case msg := <-msgChan:
				_ = enc.Encode(msg)
				flusher.Flush()
				log.Printf("sent message: %s", msg)
			case <-ctx.Done():
				// on client cancelled request
				return
			}
		}

	}
}
