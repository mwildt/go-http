package routing

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEmptyFilterChain(t *testing.T) {
	chain := FilterChain{}
	handler := chain.Build(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("EXECUTED"))
	}))
	req := httptest.NewRequest("GET", "http://example.com/subrouting/1234", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	if recorder.Body.String() != "EXECUTED" {
		t.Fail()
	}
}

func TestSingletonFilterChain(t *testing.T) {
	chain := FilterChain{
		func(w http.ResponseWriter, r *http.Request, next http.Handler) {
			w.Write([]byte("FILTER-1::"))
			next.ServeHTTP(w, r)
		},
		func(w http.ResponseWriter, r *http.Request, next http.Handler) {
			w.Write([]byte("FILTER-2::"))
			next.ServeHTTP(w, r)
		},
	}

	handler := chain.Build(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("EXECUTED"))
	}))
	req := httptest.NewRequest("GET", "http://example.com/subrouting/1234", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	if recorder.Body.String() != "FILTER-1::FILTER-2::EXECUTED" {
		t.Fail()
	}
}
