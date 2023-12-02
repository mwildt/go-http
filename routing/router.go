package routing

import (
	"context"
	"net/http"
)

type Methods []string

func (methods Methods) Compare(method string) (match bool) {
	if len(methods) == 0 {
		return true
	}
	for _, m := range methods {
		if m == method {
			return true
		}
	}
	return false
}

type Matcher struct {
	path    Segments
	methods Methods
}

func NewMatcher() Matcher {
	return Matcher{path: NewSegments("/**"), methods: make(Methods, 0)}
}

func Path(path string) Matcher {
	return NewMatcher().Path(path)
}

func Method(methods ...string) Matcher {
	return NewMatcher().Method(methods...)
}

func (m Matcher) Path(path string) Matcher {
	m.path = NewSegments(path)
	return m
}

func (m Matcher) Method(methods ...string) Matcher {
	m.methods = methods
	return m
}

type Route struct {
	matcher     Matcher
	handlerFunc http.HandlerFunc
}

type Router struct {
	routes []Route
}

func (r Router) HandleFunc(matcher Matcher, handlerFunc http.HandlerFunc) {
	r.routes = append(r.routes, Route{matcher: matcher, handlerFunc: handlerFunc})
}

func (r Router) Handle(matcher Matcher, handler http.Handler) {
	r.routes = append(r.routes, Route{matcher: matcher, handlerFunc: handler.ServeHTTP})
}

func (r Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	for _, route := range r.routes {
		if match, _, params := route.matcher.path.Compare(NewUriPath(request.URL.Path)); !match {
			continue
		} else if !route.matcher.methods.Compare(request.Method) {
			continue
		} else {
			route.handlerFunc.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), "http.path.params", params)))
			return
		}
	}
}

type RouterConfigurer func(router Router)

func CreateRoutingHandler(configuration RouterConfigurer) http.Handler {
	router := Router{}
	configuration(router)
	return &router
}
