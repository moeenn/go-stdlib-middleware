package middleware

import (
	"context"
	"net/http"
)

var contextKeyRemoteIp ContextKey = "REMOTE_IP"

func RemoteIp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipAddr := r.Header.Get("X-Real-Ip")
		if ipAddr == "" {
			ipAddr = r.Header.Get("X-Forwarded-For")
		}
		if ipAddr == "" {
			ipAddr = r.RemoteAddr
		}

		ctx := context.WithValue(r.Context(), contextKeyRemoteIp, ipAddr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRemoteIp(ctx context.Context) string {
	return ctx.Value(contextKeyRemoteIp).(string)
}
