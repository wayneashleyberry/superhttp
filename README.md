Package `superhttp` provides a thin, idiomatic layer over Go's standard net/http, adding missing features such as type-safe HTTP method handlers, middleware chaining, route groups with scoped middleware, and request context enrichment with route metadata.

It is designed as a minimal, zero-dependency shim that complements and extends the [Go 1.22+ `http.ServeMux`](https://go.dev/blog/routing-enhancements), aiming to make structured routing and middleware composition more ergonomic.

This package exists to bridge functionality that will hopefully be included in the Go standard library one day. If and when those capabilities are added, this package should become unnecessary.

### Example

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/wayneashleyberry/superhttp"
)

func main() {
	r := superhttp.NewServeMux()

	r.Use(middleware.Logger) // <--<< Logger should come before Recoverer
	r.Use(middleware.Recoverer)

	r.GET("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	}))

	r.Group("/api", func(api *superhttp.ServeMux) {
		api.Use(middleware.NoCache)

		api.GET("/users/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pattern := superhttp.RoutePattern(r)
			fmt.Fprintf(w, "Pattern: %s", pattern)

			id := r.PathValue("id")
			fmt.Fprintf(w, "ID: %s", id)
		}))
	})

	http.ListenAndServe(":8080", r)
}
```
