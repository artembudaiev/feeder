package main

import (
	"log"
	"net/http"
	"strings"
)

func main() {
	for i := 0; i < 5; i++ {
		_, err := http.Post("http://app:8088/message", "application/json", strings.NewReader(`{"text":"zalupa"}`))
		if err != nil {
			log.Println(err)
		}
	}
}
