package main

import (
	"database/sql"
	"github.com/artembudaiev/feeder/internal/message"
	"log"
	"net/http"
)

func main() {
	dbConn, err := sql.Open("postgres", "postgresql://root@roach1:26257/defaultdb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	controller := message.NewController(
		message.NewService(
			message.NewCockroachDBRepository(dbConn),
		),
	)

	http.Handle("/message", controller.HandleAdd())
	http.Handle("/messages", controller.HandleGet())
	http.ListenAndServe(":8081", nil)
}
