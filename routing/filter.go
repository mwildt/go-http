package routing

import "net/http"

type Filter func(w http.ResponseWriter, r *http.Request, handlerFunc http.HandlerFunc)

type FilterChain []Filter

func (chain FilterChain) Build(handler http.HandlerFunc) http.HandlerFunc {
	if len(chain) <= 0 {
		return handler
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		firstFilter := chain[0]
		remainingChain := chain[1:]
		firstFilter(writer, request, remainingChain.Build(handler))
	}
}

func (chain FilterChain) Extend(chain2 FilterChain) FilterChain {
	return append(chain, chain2...)
}
