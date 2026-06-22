package main

import (
	"fmt"
	"log/slog"
	"nativeweb/lib/middleware"
	"net/http"
	"os"
	"time"
)

const (
	ADDRESS        string        = "0.0.0.0:3000"
	SERVER_TIMEOUT time.Duration = time.Second * 10
)

func run() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mux := http.NewServeMux()
	loggingMiddleware := middleware.LoggingMiddleware(logger)
	authMiddleware := middleware.AuthMiddleware(logger)

	mux.HandleFunc("GET /hello", middleware.Chain(helloHandler, loggingMiddleware, authMiddleware))

	server := &http.Server{
		Addr:              ADDRESS,
		Handler:           mux,
		ReadTimeout:       SERVER_TIMEOUT,
		WriteTimeout:      SERVER_TIMEOUT,
		ReadHeaderTimeout: SERVER_TIMEOUT,
		IdleTimeout:       SERVER_TIMEOUT,
	}

	logger.Info("starting server", "address", ADDRESS)
	return server.ListenAndServe()
}

type helloResponse struct {
	Message string `json:"message"`
}

func helloHandler(c middleware.Context) error {
	token := c.Context().Value(middleware.AuthTokenContextKey).(string)
	res := helloResponse{Message: "Welcome to our website: " + token}
	return c.Json(http.StatusOK, res)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}
