package middleware

import (
	"context"
	"log/slog"
	"nativeweb/lib/responses"
	"net/http"
	"time"
)

type RequestHandler func(http.ResponseWriter, *http.Request)
type Middleware func(RequestHandler) RequestHandler
type ContextKey struct{}

var AuthTokenContextKey ContextKey = struct{}{}

func LoggingMiddleware(logger *slog.Logger) Middleware {
	m := func(next RequestHandler) RequestHandler {
		handler := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next(w, r)
			logger.Info("incoming request",
				"method", r.Method,
				"uri", r.RequestURI,
				"elapsed_ms", time.Since(start).Milliseconds(),
			)
		}

		return handler
	}

	return m
}

func AuthMiddleware(logger *slog.Logger) Middleware {
	m := func(next RequestHandler) RequestHandler {
		handler := func(w http.ResponseWriter, r *http.Request) {
			bearerToken, err := getBearerToken(r)
			if err != nil {
				ip := getRemoteIp(r)
				logger.Warn("invalid auth bearer token access",
					"reason", err.Error(),
					"ip", ip,
				)
				responses.Error(w, http.StatusUnauthorized, err)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, AuthTokenContextKey, *bearerToken)
			next(w, r.WithContext(ctx))
		}

		return handler
	}

	return m
}

func Chain(handler RequestHandler, middlewares ...Middleware) RequestHandler {
	numMiddleware := len(middlewares)
	if numMiddleware == 0 {
		return handler
	}

	fn := handler
	for i := 1; i < numMiddleware; i++ {
		fn = middlewares[i](fn)
	}

	return fn
}
