package main

import (
	"github.com/artembudaiev/feeder/internal/message"
	"net/http"
)

func main() {
	controller := message.NewController(
		message.NewService(
			message.NewInMemoryRepository(),
		),
	)

	http.Handle("/message", controller.HandleAdd())
	http.Handle("/messages", controller.HandleGet())
	http.ListenAndServe(":8080", nil)
}
