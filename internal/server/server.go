package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type Config struct {
	Handler http.Handler
	Address string
	Timeout time.Duration
}

func New(config Config) *http.Server {
	//nolint:exhaustruct_v5
	return &http.Server{
		Addr:              config.Address,
		Handler:           config.Handler,
		ReadTimeout:       config.Timeout,
		WriteTimeout:      config.Timeout,
		IdleTimeout:       config.Timeout,
		ReadHeaderTimeout: config.Timeout,
	}
}

func JsonResponse(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("failed encode response", "error", err.Error())
	}
}

type errResponse struct {
	Error string `json:"error"`
}

func ErrorResponse(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := errResponse{err.Error()}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		slog.Error("failed encode error response", "error", err.Error())
	}
}
