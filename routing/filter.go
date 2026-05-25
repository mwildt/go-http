package routing

import "net/http"

// Filter is a function that wraps an HTTP handler to provide middleware functionality.
// It receives the current writer, request, and the next handler in the chain.
type Filter func(w http.ResponseWriter, r *http.Request, next http.Handler)

type FilterChain []Filter

// Build constructs an HTTP handler by chaining all filters and the final handler.
// The filters are executed in the order they appear in the chain.
func (chain FilterChain) Build(handler http.Handler) http.Handler {
	if len(chain) <= 0 {
		return handler
	}
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		firstFilter := chain[0]
		remainingChain := chain[1:]
		firstFilter(writer, request, remainingChain.Build(handler))
	})
}

// Extend appends another FilterChain to this one, returning a new combined chain.
func (chain FilterChain) Extend(chain2 FilterChain) FilterChain {
	result := make(FilterChain, len(chain)+len(chain2))
	copy(result, chain)
	copy(result[len(chain):], chain2)
	return result
}
