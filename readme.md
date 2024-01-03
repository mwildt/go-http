# mwildt/go-http

This is a simple utility-package for handling http requests.

## Router

```go
func main() {
    router := routing.NewRouter(func(router routing.Routing) {

        router.Handle(
            routing.Path("/api/**"),
            httputil.NewSingleHostReverseProxy(serviceLocation))

        router.Handle(
            routing.Path("/**").Method("GET"),
            http.FileServer(http.Dir("frontend/build")))

    })

    err = http.ListenAndServe(
        GetEnvOrDefault("LISTEN_ADDRESS", ":3010"),
        router)
}
```
