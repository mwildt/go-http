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

type matcher struct {
	path    Segments
	methods Methods
}

type Route struct {
	matcher     matcher
	handlerFunc http.HandlerFunc
}

type Router struct {
	routes []Route
}

type Routing interface {
	HandleFunc(builder RouteBuilder, handlerFunc http.HandlerFunc)
	Handle(builder RouteBuilder, handler http.Handler)
	Route(builder RouteBuilder, configurations ...RoutingConsumer) Routing
}

type RoutingConsumer func(router Routing)

func NewRouter(configurations ...RoutingConsumer) *Router {
	router := &Router{}
	for _, configuration := range configurations {
		configuration(router)
	}
	return router
}

func (r *Router) HandleFunc(routeBuilder RouteBuilder, handlerFunc http.HandlerFunc) {
	r.addRoute(Route{
		matcher:     routeBuilder.createMatcher(),
		handlerFunc: routeBuilder.filterChain.Build(handlerFunc),
	})
}

func (r *Router) Handle(routeBuilder RouteBuilder, handler http.Handler) {
	r.addRoute(Route{matcher: routeBuilder.createMatcher(), handlerFunc: handler.ServeHTTP})
}

func (r *Router) Route(matcher RouteBuilder, configurations ...RoutingConsumer) Routing {
	router := &subrouter{
		router:       r,
		routeBuilder: matcher,
	}
	for _, configuration := range configurations {
		configuration(router)
	}
	return router
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
	router       *Router
	routeBuilder RouteBuilder
}

func (r *subrouter) HandleFunc(builder RouteBuilder, handlerFunc http.HandlerFunc) {
	extendBuilder := r.routeBuilder.extend(builder)
	r.router.addRoute(Route{
		matcher:     extendBuilder.createMatcher(),
		handlerFunc: extendBuilder.filterChain.Build(handlerFunc),
	})
}

func (r *subrouter) Handle(builder RouteBuilder, handler http.Handler) {
	r.HandleFunc(builder, handler.ServeHTTP)
}

func (r *subrouter) Route(builder RouteBuilder, configurations ...RoutingConsumer) Routing {
	router := &subrouter{r.router, r.routeBuilder.extend(builder)}
	for _, configuration := range configurations {
		configuration(router)
	}
	return router
}
