package routing

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEmptyFilterChain(t *testing.T) {

	chain := FilterChain{}
	handler := chain.Build(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("EXECUTED"))
	})
	req := httptest.NewRequest("GET", "http://example.com/subrouting/1234", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, req)

	if recorder.Body.String() != "EXECUTED" {
		t.Fail()
	}
}

func TestSingletonFilterChain(t *testing.T) {

	chain := FilterChain{
		func(w http.ResponseWriter, r *http.Request, handlerFunc http.HandlerFunc) {
			w.Write([]byte("FILTER-1::"))
			handlerFunc(w, r)
		},
		func(w http.ResponseWriter, r *http.Request, handlerFunc http.HandlerFunc) {
			w.Write([]byte("FILTER-2::"))
			handlerFunc(w, r)
		},
	}

	handler := chain.Build(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("EXECUTED"))
	})
	req := httptest.NewRequest("GET", "http://example.com/subrouting/1234", nil)
	recorder := httptest.NewRecorder()
	handler(recorder, req)

	if recorder.Body.String() != "FILTER-1::FILTER-2::EXECUTED" {
		t.Fail()
	}
}
