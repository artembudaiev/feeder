package main

import (
	"database/sql"
	"github.com/artembudaiev/feeder/internal/message"
	"log"
	"net/http"
)

func main() {
	dbConn, err := sql.Open("postgres", "postgresql://root@cockroachdb1:26257/defaultdb?sslmode=disable")
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
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("healthy"))
	})
	err = http.ListenAndServe("app:8088", nil)
	log.Println("starting http server...")
	if err != nil {
		log.Printf("failed to start a server %s", err.Error())
	}
}
