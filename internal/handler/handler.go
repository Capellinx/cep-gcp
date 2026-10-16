// Package handler expõe a API HTTP de clima por CEP.
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"cep-gcp/internal/temperature"
	"cep-gcp/internal/viacep"
)

type LocationFinder interface {
	FindCity(ctx context.Context, cep string) (string, error)
}

type TemperatureFinder interface {
	CurrentCelsius(ctx context.Context, city string) (float64, error)
}

var cepPattern = regexp.MustCompile(`^\d{8}$`)

type server struct {
	locations LocationFinder
	temps     TemperatureFinder
}

type weatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

// New devolve o handler com a rota GET /weather/{cep}.
func New(l LocationFinder, t TemperatureFinder) http.Handler {
	s := &server{locations: l, temps: t}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /weather/{cep}", s.weather)
	return mux
}

func (s *server) weather(w http.ResponseWriter, r *http.Request) {
	cep := r.PathValue("cep")
	if !cepPattern.MatchString(cep) {
		writeText(w, http.StatusUnprocessableEntity, "invalid zipcode")
		return
	}

	city, err := s.locations.FindCity(r.Context(), cep)
	if errors.Is(err, viacep.ErrNotFound) {
		writeText(w, http.StatusNotFound, "can not find zipcode")
		return
	}
	if err != nil {
		upstream(w, err)
		return
	}

	c, err := s.temps.CurrentCelsius(r.Context(), city)
	if err != nil {
		upstream(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(weatherResponse{
		TempC: c,
		TempF: temperature.CelsiusToFahrenheit(c),
		TempK: temperature.CelsiusToKelvin(c),
	})
}

func upstream(w http.ResponseWriter, err error) {
	log.Printf("upstream error: %v", err)
	writeText(w, http.StatusBadGateway, "upstream error")
}

func writeText(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(msg))
}
