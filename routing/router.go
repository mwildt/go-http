package routing

import (
	"context"
	"log"
	"net/http"
)

type contextKey string

const (
	contextParams = contextKey("router.http.params")
)

func WithParameters(c context.Context, parameters Parameters) context.Context {
	return context.WithValue(c, contextParams, parameters)
}

func GetParameters(c context.Context) Parameters {
	if value := c.Value(contextParams); value != nil {
		return value.(Parameters)
	} else {
		return make(Parameters)
	}
}

func GetParameter(c context.Context, key string) (string, bool) {
	params := GetParameters(c)
	value, exists := params[key]
	return value, exists
}

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

func (methods Methods) Extend(methods2 Methods) Methods {
	if len(methods) == 0 {
		return methods2
	} else if len(methods2) == 0 {
		return methods
	} else {
		response := make(Methods, 0)
		for _, m1 := range methods {
			for _, m2 := range methods2 {
				if m1 == m2 {
					response = append(response, m1)
				}
			}
		}
		if len(response) == 0 {
			log.Fatal("illegal routing configuration")
		}
		return response
	}
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

func (m Matcher) Get(path string) Matcher {
	return Method("GET").Path(path)
}

func (m Matcher) Post(path string) Matcher {
	return Method("POST").Path(path)
}

func (m Matcher) Patch(path string) Matcher {
	return Method("PATCH").Path(path)
}

func (m Matcher) Put(path string) Matcher {
	return Method("PUT").Path(path)
}

func (m Matcher) Delete(path string) Matcher {
	return Method("DELETE").Path(path)
}

func (m Matcher) Path(path string) Matcher {
	m.path = NewSegments(path)
	return m
}

func (m Matcher) Method(methods ...string) Matcher {
	m.methods = methods
	return m
}

func (m Matcher) extend(matcher Matcher) Matcher {

	return Matcher{
		path:    m.path.Extend(matcher.path),
		methods: m.methods.Extend(matcher.methods),
	}
}

type Route struct {
	matcher     Matcher
	handlerFunc http.HandlerFunc
}

type Router struct {
	routes []Route
}

type Routing interface {
	HandleFunc(matcher Matcher, handlerFunc http.HandlerFunc)
	Handle(matcher Matcher, handler http.Handler)
	Route(matcher Matcher, configurations ...RoutingConsumer) Routing
}

type RoutingConsumer func(router Routing)

func NewRouter(configurations ...RoutingConsumer) *Router {
	router := &Router{}
	for _, configuration := range configurations {
		configuration(router)
	}
	return router
}

func (r *Router) HandleFunc(matcher Matcher, handlerFunc http.HandlerFunc) {
	r.addRoute(Route{matcher: matcher, handlerFunc: handlerFunc})
}

func (r *Router) Handle(matcher Matcher, handler http.Handler) {
	r.addRoute(Route{matcher: matcher, handlerFunc: handler.ServeHTTP})
}

func (r *Router) Route(matcher Matcher, conf ...RoutingConsumer) Routing {
	return &subrouter{
		router:  r,
		matcher: matcher,
	}
}

func (r *Router) addRoute(route Route) {
	r.routes = append(r.routes, route)
}

func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	for _, route := range r.routes {
		if match, _, params := route.matcher.path.Compare(NewUriPath(request.URL.Path)); !match {
			continue
		} else if !route.matcher.methods.Compare(request.Method) {
			continue
		} else {
			route.handlerFunc.ServeHTTP(writer, request.WithContext(WithParameters(request.Context(), params)))
			return
		}
	}
}

func DefaultNotFound() RoutingConsumer {
	return func(router Routing) {
		router.HandleFunc(Path("/**"), http.NotFound)
	}
}

type subrouter struct {
	router  *Router
	matcher Matcher
}

func (r *subrouter) HandleFunc(matcher Matcher, handlerFunc http.HandlerFunc) {
	r.router.addRoute(Route{matcher: r.matcher.extend(matcher), handlerFunc: handlerFunc})
}

func (r *subrouter) Handle(matcher Matcher, handler http.Handler) {
	r.router.addRoute(Route{matcher: r.matcher.extend(matcher), handlerFunc: handler.ServeHTTP})
}

func (r *subrouter) Route(matcher Matcher, configurations ...RoutingConsumer) Routing {
	router := &subrouter{r.router, r.matcher.extend(matcher)}
	for _, configuration := range configurations {
		configuration(router)
	}
	return router
}
