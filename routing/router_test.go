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
