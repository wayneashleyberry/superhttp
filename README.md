### Usage

```go
package main

import "github.com/wayneashleyberry/superhttp"

func main() {
    r := superhttp.NewServeMux()

    r.Use(loggingMiddleware)

    r.GET("/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello!")
    }))

    r.Group("/api", func(api *router.Router) {
        api.Use(apiMiddleware)

        api.GET("/users/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            pattern := router.RoutePattern(r)
            fmt.Fprintf(w, "Pattern: %s", pattern)
        }))
    })

    http.ListenAndServe(":8080", r)
}
```
