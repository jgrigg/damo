package util

import (
	"encoding/json"
	"net/http"
)

type checkResponse struct {
	Message string `json:"message"`
}

// WriteError writes to the stream in the common error format
func WriteError(rw *http.ResponseWriter, status int, message string) {
	resultBody, _ := json.Marshal(checkResponse{Message: message})
	(*rw).Header().Set("Content-Type", "application/json")
	(*rw).WriteHeader(status)
	(*rw).Write(resultBody)
}
