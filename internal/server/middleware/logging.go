package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := fmt.Sprintf("%s %s", r.Method, r.URL.String())
		reqId, reqIdOk := GetRequestId(r.Context())
		if !reqIdOk {
			slog.Warn("RequestId middleware should be registered before Logging middleware")
			reqId = "missing-request-id"
		}

		ws := &responseWriterWithStatus{w, http.StatusOK}
		next.ServeHTTP(ws, r)
		elapsedMs := time.Since(start).Microseconds()

		level := slog.LevelInfo
		switch {
		case ws.status >= 400:
			level = slog.LevelWarn
		case ws.status >= 500:
			level = slog.LevelError
		}

		slog.Log(r.Context(), level, "new request",
			"path", path,
			"requestId", reqId,
			"status", ws.status,
			"elapsedMs", elapsedMs,
		)
	})
}
