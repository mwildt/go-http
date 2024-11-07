package httputils

import (
	"encoding/json"
	"net/http"
)

func Send(w http.ResponseWriter, request *http.Request, code int) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
}

func Unauthorized(w http.ResponseWriter, request *http.Request) {
	Send(w, request, http.StatusUnauthorized)
}

func NotFound(w http.ResponseWriter, request *http.Request) {
	Send(w, request, http.StatusNotFound)
}

func BadRequest(w http.ResponseWriter, request *http.Request) {
	Send(w, request, http.StatusBadRequest)
}

func InternalServerError(w http.ResponseWriter, request *http.Request) {
	Send(w, request, http.StatusInternalServerError)
}

func SendJson(w http.ResponseWriter, request *http.Request, code int, data interface{}) {
	payload, err := json.Marshal(data)
	if err != nil {
		InternalServerError(w, request)
	} else {
		w.Header().Set("Content-Type", "application/json")
		Send(w, request, code)
		w.Write(payload)
	}
}
func OkJson(w http.ResponseWriter, request *http.Request, data interface{}) {
	SendJson(w, request, http.StatusOK, data)
}

func CreatedJson(w http.ResponseWriter, request *http.Request, data interface{}) {
	SendJson(w, request, http.StatusCreated, data)
}
