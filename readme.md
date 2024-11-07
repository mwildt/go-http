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
            http.FileServer(http.Dir("/static/")))

    })

    err = http.ListenAndServe(":8080", router)
}
```
