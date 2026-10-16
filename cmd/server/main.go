package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"cep-gcp/internal/handler"
	"cep-gcp/internal/viacep"
	"cep-gcp/internal/weather"
)

func main() {
	key := os.Getenv("WEATHER_API_KEY")
	if key == "" {
		log.Fatal("WEATHER_API_KEY is required")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           handler.New(viacep.NewClient(), weather.NewClient(key)),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("listening on :%s", port)
	log.Fatal(srv.ListenAndServe())
}
