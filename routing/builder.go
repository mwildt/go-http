package routing

// RouteBuilder is used to construct route definitions with path, methods, and filters.
type RouteBuilder struct {
	path        Segments
	methods     Methods
	filterChain FilterChain
}

// NewRouteBuilder creates a new empty RouteBuilder.
func NewRouteBuilder() RouteBuilder {
	return RouteBuilder{path: Segments{}, methods: make(Methods, 0), filterChain: FilterChain{}}
}

// extend combines this RouteBuilder with another, merging paths, methods, and filters.
func (builder RouteBuilder) extend(extension RouteBuilder) RouteBuilder {
	return RouteBuilder{
		path:        builder.path.Extend(extension.path),
		methods:     builder.methods.Extend(extension.methods),
		filterChain: builder.filterChain.Extend(extension.filterChain),
	}
}

// Filtering creates a new RouteBuilder with the given filter.
func Filtering(filter Filter) RouteBuilder {
	return NewRouteBuilder().Filter(filter)
}

// Path creates a new RouteBuilder with the given path.
func Path(path string) RouteBuilder {
	return NewRouteBuilder().Path(path)
}

// Method creates a new RouteBuilder with the given HTTP methods.
func Method(methods ...string) RouteBuilder {
	return NewRouteBuilder().Method(methods...)
}

// Get creates a new RouteBuilder for GET requests to the given path.
func Get(path string) RouteBuilder {
	return Method("GET").Path(path)
}

// Post creates a new RouteBuilder for POST requests to the given path.
func Post(path string) RouteBuilder {
	return Method("POST").Path(path)
}

// Patch creates a new RouteBuilder for PATCH requests to the given path.
func Patch(path string) RouteBuilder {
	return Method("PATCH").Path(path)
}

// Put creates a new RouteBuilder for PUT requests to the given path.
func Put(path string) RouteBuilder {
	return Method("PUT").Path(path)
}

// Delete creates a new RouteBuilder for DELETE requests to the given path.
func Delete(path string) RouteBuilder {
	return Method("DELETE").Path(path)
}

// Path extends the current path with the given path string.
func (builder RouteBuilder) Path(path string) RouteBuilder {
	builder.path = builder.path.Extend(NewSegments(path))
	return builder
}

// Method sets the HTTP methods for this route.
func (builder RouteBuilder) Method(methods ...string) RouteBuilder {
	builder.methods = methods
	return builder
}

// Filter adds a filter to the filter chain.
func (builder RouteBuilder) Filter(filter Filter) RouteBuilder {
	builder.filterChain = append(builder.filterChain, filter)
	return builder
}

// createMatcher creates a matcher from the current builder state.
func (builder RouteBuilder) createMatcher() matcher {
	return matcher{path: builder.path, methods: builder.methods}
}
