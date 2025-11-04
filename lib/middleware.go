package lib

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type RequestHandler func(http.ResponseWriter, *http.Request)
type Middleware interface {
	Handle(http.Handler) http.Handler
}

func getRemoteIp(r *http.Request) string {
	ipAddress := r.Header.Get("X-Real-Ip")
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}
	return ipAddress
}

type LoggingMiddleware struct {
	logger *slog.Logger
}

func NewLoggingMiddleware(logger *slog.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{logger}
}

func (m *LoggingMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		m.logger.Info("incoming request",
			"method", r.Method,
			"uri", r.RequestURI,
			"elapsed_ms", time.Since(start).Milliseconds(),
		)
	})
}

type AuthMiddleware struct {
	logger *slog.Logger
}

func NewAuthMiddleware(logger *slog.Logger) *AuthMiddleware {
	return &AuthMiddleware{logger}
}

type ContextKey struct{}

var AuthContextKey ContextKey = struct{}{}

func (m *AuthMiddleware) getBearerToken(r *http.Request) (*string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("authorization header missing")
	}

	prefix := "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return nil, errors.New("unexpected authorization scheme provided")
	}

	authHeader = strings.TrimPrefix(authHeader, prefix)
	if len(authHeader) == 0 {
		return nil, errors.New("missing bearer token value in header")
	}

	authHeader = strings.TrimSpace(authHeader)
	return &authHeader, nil
}

func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken, err := m.getBearerToken(r)
		if err != nil {
			ip := getRemoteIp(r)
			m.logger.Warn("invalid auth bearer token access",
				"reason", err.Error(),
				"ip", ip,
			)
			ErrorRespond(w, http.StatusUnauthorized, err)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, AuthContextKey, bearerToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Chain(handler RequestHandler, middlewares ...Middleware) http.Handler {
	numMiddleware := len(middlewares)

	fn := http.HandlerFunc(handler)
	if numMiddleware == 0 {
		return fn
	}

	h := middlewares[0].Handle(fn)
	for i := 1; i < numMiddleware; i++ {
		h = middlewares[i].Handle(h)
	}

	return h
}
