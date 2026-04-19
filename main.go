package main

import (
	"log"
	"net/http"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := NewApp(cfg)

	log.Printf("Listening on :%s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, app.Handler()); err != nil {
		log.Fatal(err)
	}
}
