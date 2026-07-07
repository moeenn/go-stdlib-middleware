package middleware

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

func (c Context) Context() context.Context {
	return c.Request.Context()
}

func (c *Context) ContextWithValue(key any, val any) {
	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, key, val)
	c.Request = c.Request.WithContext(ctx)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (c Context) Json(statusCode int, data any) error {
	c.Writer.WriteHeader(statusCode)
	return json.NewEncoder(c.Writer).Encode(data)
}

func (c Context) Error(statusCode int, err error) error {
	c.Writer.WriteHeader(statusCode)
	res := &ErrorResponse{err.Error()}
	return json.NewEncoder(c.Writer).Encode(res)
}

func toNativeHandler(h RequestHandler) NativeRequestHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(Context{Writer: w, Request: r})
		if err != nil {
			w.WriteHeader(500)
			res := ErrorResponse{Error: err.Error()}
			_ = json.NewEncoder(w).Encode(res)
		}
	}
}

type NativeRequestHandler func(http.ResponseWriter, *http.Request)
type RequestHandler func(Context) error
type Middleware func(RequestHandler) RequestHandler
type ContextKey struct{}

var AuthTokenContextKey ContextKey = struct{}{}

func LoggingMiddleware(logger *slog.Logger) Middleware {
	m := func(next RequestHandler) RequestHandler {
		handler := func(c Context) error {
			start := time.Now()
			_ = next(c)

			// FIXME: this never gets called.
			logger.Info("incoming request",
				"method", c.Request.Method,
				"uri", c.Request.RequestURI,
				"elapsedMs", time.Since(start).Milliseconds(),
			)
			return nil
		}
		return handler
	}
	return m
}

func AuthMiddleware(logger *slog.Logger) Middleware {
	m := func(next RequestHandler) RequestHandler {
		handler := func(c Context) error {
			bearerToken, err := getBearerToken(c.Request)
			if err != nil {
				ip := getRemoteIp(c.Request)
				logger.Warn("invalid auth bearer token access",
					"reason", err.Error(),
					"ip", ip,
				)
				return c.Error(http.StatusUnauthorized, err)
			}

			c.ContextWithValue(AuthTokenContextKey, *bearerToken)
			return next(c)
		}
		return handler
	}
	return m
}

func Chain(handler RequestHandler, middlewares ...Middleware) NativeRequestHandler {
	numMiddleware := len(middlewares)
	if numMiddleware == 0 {
		return toNativeHandler(handler)
	}

	fn := handler
	for i := 1; i < numMiddleware; i++ {
		fn = middlewares[i](fn)
	}

	return toNativeHandler(fn)
}
