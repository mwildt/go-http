// Package routing provides a flexible HTTP request router with support for
// path parameters, wildcards, method-based routing, and middleware filters.
package routing

import (
	"context"
	"log"
	"net/http"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

const (
	// contextParams is the context key for storing URL parameters.
	contextParams = contextKey("router.http.params")
)

// WithParameters adds URL parameters to the context.
func WithParameters(c context.Context, parameters Parameters) context.Context {
	return context.WithValue(c, contextParams, parameters)
}

// GetParameters retrieves URL parameters from the context.
// Returns an empty Parameters map if none are present.
func GetParameters(c context.Context) Parameters {
	if value := c.Value(contextParams); value != nil {
		return value.(Parameters)
	}
	return make(Parameters)
}

// GetParameter retrieves a single URL parameter from the context.
// Returns the parameter value and a boolean indicating if it exists.
func GetParameter(c context.Context, key string) (string, bool) {
	params := GetParameters(c)
	value, exists := params[key]
	return value, exists
}

// Methods represents a list of HTTP methods (e.g., "GET", "POST").
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
	}
	if len(methods2) == 0 {
		return methods
	}
	// Deduplizierung: Alle Methoden aus beiden Slices zusammenführen
	seen := make(map[string]struct{})
	for _, m := range methods {
		seen[m] = struct{}{}
	}
	for _, m := range methods2 {
		seen[m] = struct{}{}
	}
	result := make(Methods, 0, len(seen))
	for m := range seen {
		result = append(result, m)
	}
	if len(result) == 0 {
		log.Fatal("illegal routing configuration: no methods left after extend")
	}
	return result
}

// matcher holds the path and method matching criteria for a route.
type matcher struct {
	path    Segments
	methods Methods
}

// Route represents a single route with its matcher and handler.
type Route struct {
	matcher     matcher
	handlerFunc http.HandlerFunc
}

// Router is an HTTP request multiplexer that matches requests to registered routes.
// It supports path parameters, wildcards, and method-based routing.
type Router struct {
	routes []Route
}

// Routing is the interface for configuring routes.
type Routing interface {
	HandleFunc(builder RouteBuilder, handlerFunc http.HandlerFunc)
	Handle(builder RouteBuilder, handler http.Handler)
	Route(builder RouteBuilder, configurations ...RoutingConsumer) Routing
}

// RoutingConsumer is a function that configures a Routing instance.
type RoutingConsumer func(router Routing)

// NewRouter creates a new Router and applies the given configurations.
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
		handlerFunc: routeBuilder.filterChain.Build(handlerFunc).ServeHTTP,
	})
}

func (r *Router) Handle(routeBuilder RouteBuilder, handler http.Handler) {
	r.addRoute(Route{matcher: routeBuilder.createMatcher(), handlerFunc: routeBuilder.filterChain.Build(handler).ServeHTTP})
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

// ServeHTTP implements http.Handler, matching the request to a registered route.
func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	for _, route := range r.routes {
		if match, _, params := route.matcher.path.Compare(NewUriPath(request.URL.Path)); !match {
			continue
		}
		if !route.matcher.methods.Compare(request.Method) {
			continue
		}
		route.handlerFunc.ServeHTTP(writer, request.WithContext(WithParameters(request.Context(), params)))
		return
	}
	// No route matched - return 404
	http.NotFound(writer, request)
}

// DefaultNotFound returns a RoutingConsumer that adds a default 404 handler for unmatched routes.
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
		handlerFunc: extendBuilder.filterChain.Build(handlerFunc).ServeHTTP,
	})
}

func (r *subrouter) Handle(builder RouteBuilder, handler http.Handler) {
	extendBuilder := r.routeBuilder.extend(builder)
	r.router.addRoute(Route{
		matcher:     extendBuilder.createMatcher(),
		handlerFunc: extendBuilder.filterChain.Build(handler).ServeHTTP,
	})
}

func (r *subrouter) Route(builder RouteBuilder, configurations ...RoutingConsumer) Routing {
	router := &subrouter{r.router, r.routeBuilder.extend(builder)}
	for _, configuration := range configurations {
		configuration(router)
	}
	return router
}
