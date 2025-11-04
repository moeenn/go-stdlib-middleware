package lib

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func Respond[T any](w http.ResponseWriter, statusCode int, data T) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ErrorRespond(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	res := &ErrorResponse{err.Error()}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
