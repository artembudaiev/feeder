package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/artembudaiev/feeder/internal/config"
	"github.com/artembudaiev/feeder/internal/message"
	"log"
	"net/http"
)

func main() {
	cfgManager, err := config.NewEnvAppManager()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to initialize config %w", err))
	}
	//dbConn, err := sql.Open("postgres", "postgresql://root@cockroachdb1:26257/defaultdb?sslmode=disable")
	dbConn, err := sql.Open("postgres", cfgManager.GetConfig().DbUrl)

	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	svc := message.NewService(
		message.NewCockroachDBRepository(dbConn),
	)
	// todo: implement graceful shutdown
	go func() {
		err = svc.Start(context.Background())
		if err != nil {
			log.Fatal(err)
			return
		}
	}()
	controller := message.NewController(
		svc,
	)

	http.Handle("/message", controller.HandleAdd())
	http.Handle("/messages", controller.HandleGet())
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("healthy"))
	})
	addr := fmt.Sprintf("%s:%s", cfgManager.GetConfig().AppHost, cfgManager.GetConfig().AppPort)
	log.Printf("starting http server on address %s...", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Printf("failed to start a server %s", err.Error())
	}
}
