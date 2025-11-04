package main

import (
	"fmt"
	"log/slog"
	"nativeweb/lib"
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
	mux := http.NewServeMux()
	loggingMiddleware := lib.NewLoggingMiddleware(logger)
	authMiddleware := lib.NewAuthMiddleware(logger)

	mux.Handle("GET /hello", lib.Chain(helloHandler, loggingMiddleware, authMiddleware))

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

func helloHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(lib.AuthContextKey).(string)

	res := helloResponse{Message: "Welcome to our website: " + token}
	lib.Respond(w, http.StatusOK, res)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}
