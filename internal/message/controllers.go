package message

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Service interface {
	Add(ctx context.Context, message Message) error
	Attach(ctx context.Context) (<-chan Message, error)
}

type Controller struct {
	service Service
}

func (c *Controller) HandleAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var msg Message
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&msg); err != nil {
			http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		}
		if err := c.service.Add(r.Context(), msg); err != nil {
			http.Error(w, "failed to add message", http.StatusInternalServerError)
		}
		// todo: remove
		fmt.Printf("added message: %s", msg)
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
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.WriteHeader(http.StatusOK)
		flusher.Flush()

		enc := json.NewEncoder(w)

		msgChan, err := c.service.Attach(ctx)
		if err != nil {
			http.Error(w, "failed to listen to the new messages", http.StatusInternalServerError)
		}
		for {
			select {
			case msg := <-msgChan:
				_ = enc.Encode(msg)
				flusher.Flush()
			case <-ctx.Done():
				// on client cancelled request
				return
			}
		}

	}
}
