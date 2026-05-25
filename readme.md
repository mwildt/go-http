# mwildt/go-http

A simple, flexible HTTP routing library for Go with support for:
- Path parameters (e.g., `/users/{id}`)
- Wildcards (e.g., `/api/*` or `/api/**`)
- Method-based routing (GET, POST, etc.)
- Middleware filters
- Subrouting

## Installation

```bash
go get github.com/mwildt/go-http
```

## Basic Usage

### Simple Routing

```go
package main

import (
    "net/http"
    "github.com/mwildt/go-http/routing"
)

func main() {
    router := routing.NewRouter()
    
    router.HandleFunc(routing.Get("/hello"), func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
    
    http.ListenAndServe(":8080", router)
}
```

### Path Parameters

```go
router.HandleFunc(routing.Get("/users/{id}"), func(w http.ResponseWriter, r *http.Request) {
    params := routing.GetParameters(r.Context())
    id := params["id"]
    w.Write([]byte("User ID: " + id))
})
```

### Wildcards

```go
// Single wildcard matches one path segment
router.HandleFunc(routing.Get("/api/*"), func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("API endpoint"))
})

// Global wildcard matches the rest of the path
router.HandleFunc(routing.Get("/static/**"), func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Static file"))
})
```

### Middleware Filters

```go
// Define a filter
logFilter := func(w http.ResponseWriter, r *http.Request, next http.Handler) {
    log.Printf("Request: %s %s", r.Method, r.URL.Path)
    next.ServeHTTP(w, r)
}

// Apply filter to a route
router.HandleFunc(
    routing.Get("/api/users").Filter(logFilter),
    func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Users list"))
    },
)
```

### Subrouting

```go
router.Route(routing.Path("/api"), func(api routing.Routing) {
    api.HandleFunc(routing.Get("/users"), func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("API Users"))
    })
    
    api.HandleFunc(routing.Get("/posts"), func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("API Posts"))
    })
})
```

### Combining Path and Method

```go
router.HandleFunc(
    routing.Path("/api").Method("GET", "POST"),
    func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("API endpoint"))
    },
)
```

### Using with httputils

```go
import "github.com/mwildt/go-http/httputils"

router.HandleFunc(routing.Get("/api/data"), func(w http.ResponseWriter, r *http.Request) {
    data := map[string]string{"message": "Hello"}
    httputils.OkJson(w, r, data)
})
```

## Features

- **Fast**: No reflection, minimal allocations
- **Flexible**: Supports various routing patterns
- **Middleware**: Chainable filters for request processing
- **Context**: URL parameters available via context
- **Compatible**: Works with standard `http.Handler`

## License

MIT
