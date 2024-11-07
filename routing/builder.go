package routing

type RouteBuilder struct {
	path        Segments
	methods     Methods
	filterChain FilterChain
}

func NewRouteBuilder() RouteBuilder {
	return RouteBuilder{path: Segments{}, methods: make(Methods, 0), filterChain: FilterChain{}}
}

func (builder RouteBuilder) extend(extension RouteBuilder) RouteBuilder {

	return RouteBuilder{
		path:        builder.path.Extend(extension.path),
		methods:     builder.methods.Extend(extension.methods),
		filterChain: builder.filterChain.Extend(extension.filterChain),
	}
}

func Filtering(filter Filter) RouteBuilder {
	return NewRouteBuilder().Filter(filter)
}

func Path(path string) RouteBuilder {
	return NewRouteBuilder().Path(path)
}

func Method(methods ...string) RouteBuilder {
	return NewRouteBuilder().Method(methods...)
}

func Get(path string) RouteBuilder {
	return Method("GET").Path(path)
}

func Post(path string) RouteBuilder {
	return Method("POST").Path(path)
}

func Patch(path string) RouteBuilder {
	return Method("PATCH").Path(path)
}

func Put(path string) RouteBuilder {
	return Method("PUT").Path(path)
}

func Delete(path string) RouteBuilder {
	return Method("DELETE").Path(path)
}

func (builder RouteBuilder) Path(path string) RouteBuilder {
	builder.path = NewSegments(path)
	return builder
}

func (builder RouteBuilder) Method(methods ...string) RouteBuilder {
	builder.methods = methods
	return builder
}

func (builder RouteBuilder) Filter(filter Filter) RouteBuilder {
	builder.filterChain = append(builder.filterChain, filter)
	return builder
}

func (builder RouteBuilder) createMatcher() matcher {
	return matcher{path: builder.path, methods: builder.methods}
}
