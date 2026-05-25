// Package httputils provides utility functions for HTTP responses.
package httputils

import (
	"encoding/json"
	"log"
	"net/http"
)

// Send sends an HTTP response with the given status code and security headers.
func Send(w http.ResponseWriter, request *http.Request, code int) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
}

// Unauthorized sends a 401 Unauthorized response.
func Unauthorized(w http.ResponseWriter, request *http.Request) {
	Send(w, request, http.StatusUnauthorized)
}

// NotFound sends a 404 Not Found response.
func NotFound(w http.ResponseWriter, request *http.Request) {
	Send(w, request, http.StatusNotFound)
}

// BadRequest sends a 400 Bad Request response.
func BadRequest(w http.ResponseWriter, request *http.Request) {
	Send(w, request, http.StatusBadRequest)
}

// InternalServerError sends a 500 Internal Server Error response.
func InternalServerError(w http.ResponseWriter, request *http.Request) {
	Send(w, request, http.StatusInternalServerError)
}

// SendJson sends a JSON response with the given status code and data.
// If marshaling fails, it logs the error and sends a 500 response.
func SendJson(w http.ResponseWriter, request *http.Request, code int, data interface{}) {
	payload, err := json.Marshal(data)
	if err != nil {
		log.Printf("json marshal error: %v", err)
		InternalServerError(w, request)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	Send(w, request, code)
	if _, err := w.Write(payload); err != nil {
		log.Printf("write error: %v", err)
	}
}

// OkJson sends a JSON response with status 200 OK.
func OkJson(w http.ResponseWriter, request *http.Request, data interface{}) {
	SendJson(w, request, http.StatusOK, data)
}

// CreatedJson sends a JSON response with status 201 Created.
func CreatedJson(w http.ResponseWriter, request *http.Request, data interface{}) {
	SendJson(w, request, http.StatusCreated, data)
}
