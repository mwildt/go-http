package routing

import (
	go_http "github.com/mwildt/go-http"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultNotFound(t *testing.T) {
	router := NewRouter(DefaultNotFound())

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	go_http.Assert(t, w.Body.String() == "404 page not found\n", "unexpected response body %s", w.Body.String())
	go_http.Assert(t, w.Code == 404, "unexpected response status %d", w.Code)
}

func TestSubRouting(t *testing.T) {
	router := NewRouter()
	router.Route(Path("/subrouting")).HandleFunc(Path("/{id}"), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(418)
		writer.Write([]byte(NewSegments("/id/{id}").Print(GetParameters(request.Context()))))
	})

	req := httptest.NewRequest("GET", "http://example.com/subrouting/1234", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	go_http.Assert(t, recorder.Body.String() == "/id/1234", "unexpected response body %s", recorder.Body.String())
	go_http.Assert(t, recorder.Code == 418, "unexpected response status %d", recorder.Code)
}

func TestSubRoutingConsume(t *testing.T) {
	router := NewRouter()
	router.Route(NewRouteBuilder(), func(r Routing) {
		r.HandleFunc(Get("/api/test"), func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte("TEST"))
		})
	})

	req := httptest.NewRequest("GET", "http://example.com/api/test", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	go_http.Assert(t, recorder.Body.String() == "TEST", "unexpected response body %s", recorder.Body.String())
}

func TestSubRoutingParameters(t *testing.T) {
	router := NewRouter()
	router.Route(Path("/routing/{contextId}")).HandleFunc(Path("/sub/{id}"), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(418)
		writer.Write([]byte(NewSegments("/{contextId}/{id}").Print(GetParameters(request.Context()))))
	})

	req := httptest.NewRequest("GET", "http://example.com/routing/r1/sub/s1", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	go_http.Assert(t, recorder.Body.String() == "/r1/s1", "unexpected response body %s", recorder.Body.String())
	go_http.Assert(t, recorder.Code == 418, "unexpected response status %d", recorder.Code)
}

func TestRoutingWithFilterInSubrouting(t *testing.T) {
	router := NewRouter()

	var logs []string

	logFilter := func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		logs = append(logs, r.URL.String())
		next(w, r)
	}
	router.Route(Path("/routing/test").Filter(logFilter)).HandleFunc(Path("/sub/test"), func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("EXECUTED"))
	})

	req := httptest.NewRequest("GET", "http://example.com/routing/test/sub/test", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	go_http.Assert(t, recorder.Body.String() == "EXECUTED", "unexpected response body %s", recorder.Body.String())
	go_http.Assert(t, len(logs) == 1, "unexpected lenght of log %d", len(logs))
	go_http.Assert(t, logs[0] == "http://example.com/routing/test/sub/test", "unexpected frist log statement", len(logs[0]))
}

func TestRoutingWithFilter(t *testing.T) {
	router := NewRouter()

	var logs []string

	logFilter := func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		logs = append(logs, r.URL.String())
		next(w, r)
	}

	router.HandleFunc(Path("/routing/test").Filter(logFilter), func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("EXECUTED"))
	})

	req := httptest.NewRequest("GET", "http://example.com/routing/test", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	go_http.Assert(t, recorder.Body.String() == "EXECUTED", "unexpected response body %s", recorder.Body.String())
	go_http.Assert(t, len(logs) == 1, "unexpected lenght of log %d", len(logs))
	go_http.Assert(t, logs[0] == "http://example.com/routing/test", "unexpected frist log statement", len(logs[0]))
}

func TestRoutingWithMulltipleFilters(t *testing.T) {
	router := NewRouter()

	var logs []string

	logFilter := func(value string) Filter {
		return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			logs = append(logs, value)
			next(w, r)
		}
	}

	router.Route(Path("/routing/test").Filter(logFilter("a")).Filter(logFilter("b"))).HandleFunc(Path("/sub/test").Filter(logFilter("c")), func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("EXECUTED"))
	})

	req := httptest.NewRequest("GET", "http://example.com/routing/test/sub/test", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	go_http.Assert(t, recorder.Body.String() == "EXECUTED", "unexpected response body '%s'", recorder.Body.String())
	go_http.Assert(t, len(logs) == 3, "unexpected lenght of log %d", len(logs))
	go_http.Assert(t, logs[0] == "a", "unexpected 0th log statement", len(logs[0]))
	go_http.Assert(t, logs[1] == "b", "unexpected 1st log statement", len(logs[1]))
	go_http.Assert(t, logs[2] == "c", "unexpected 2nd log statement", len(logs[2]))
}

func TestRoutingWithAlternatives(t *testing.T) {
	router := NewRouter()

	handler := func(status int, responseValue string) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(status)
			writer.Write([]byte(responseValue))
		}
	}

	router.HandleFunc(Patch("/routing/{id}/sub1"), handler(210, "sub1"))
	router.HandleFunc(Patch("/routing/{id}/sub2"), handler(211, "sub2"))

	req := httptest.NewRequest("PATCH", "http://example.com/routing/abc/sub2", nil)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	go_http.Assert(t, recorder.Code == 211, "unexpected status code '%d'", recorder.Code)
	go_http.Assert(t, recorder.Body.String() == "sub2", "unexpected response body '%s'", recorder.Body.String())

}
