package main

import (
	"fmt"
	"github.com/artembudaiev/feeder/internal/config"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	cfgManager, err := config.NewEnvClientManager()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to initialize config %w", err))
	}
	counter := 0
	for {
		time.Sleep(time.Millisecond * time.Duration(cfgManager.GetConfig().SpamTimeoutMs))
		_, err := http.Post(fmt.Sprintf("http://%s:%s/message", cfgManager.GetConfig().AppHost, cfgManager.GetConfig().AppPort), "application/json", strings.NewReader(fmt.Sprintf(`{"text":"message%d"}`, counter)))
		if err != nil {
			log.Println(err)
		}
		counter++
	}
}
