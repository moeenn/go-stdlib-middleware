package main

import (
	"api/internal/server"
	"api/internal/server/middleware"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

const (
	SERVER_ADDR    = "0.0.0.0:5000"
	SERVER_TIMEOUT = time.Second * 10
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

type ProtectedResponse struct {
	RequestId   string `json:"requestId"`
	BearerToken string `json:"bearerToken"`
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	requestId, _ := middleware.GetRequestId(r.Context())
	bearerToken := middleware.GetBearerToken(r.Context())
	res := ProtectedResponse{requestId, bearerToken}
	server.JsonResponse(w, http.StatusOK, res)
}

func run() error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /home", homeHandler)
	mux.Handle("GET /protected", middleware.Chain(protectedHandler, middleware.BearerToken))

	reqIdMiddleware := middleware.RequestId(middleware.RequestIdArgs{
		Factory:        func() string { return uuid.NewString() },
		ReadFromHeader: true,
	})

	s := server.New(server.Config{
		Address: SERVER_ADDR,
		Timeout: SERVER_TIMEOUT,
		Handler: middleware.ChainMux(mux, reqIdMiddleware, middleware.Logging),
	})

	logger.Info("starting server", "address", SERVER_ADDR)
	return s.ListenAndServe()
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}
