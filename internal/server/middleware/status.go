package middleware

import "net/http"

type responseWriterWithStatus struct {
	http.ResponseWriter
	status int
}

func (r *responseWriterWithStatus) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
