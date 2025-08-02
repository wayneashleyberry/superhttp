> Package `superhttp` provides a thin, idiomatic layer over Go's standard net/http, adding missing features such as type-safe HTTP method handlers, middleware chaining, route groups with scoped middleware, and request context enrichment with route metadata.

[![Go Reference](https://pkg.go.dev/badge/github.com/wayneashleyberry/superhttp.svg)](https://pkg.go.dev/github.com/wayneashleyberry/superhttp)
[![Test](https://github.com/wayneashleyberry/superhttp/actions/workflows/test.yaml/badge.svg)](https://github.com/wayneashleyberry/superhttp/actions/workflows/test.yaml)
[![Lint](https://github.com/wayneashleyberry/superhttp/actions/workflows/lint.yaml/badge.svg)](https://github.com/wayneashleyberry/superhttp/actions/workflows/lint.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/wayneashleyberry/superhttp)](https://goreportcard.com/report/github.com/wayneashleyberry/superhttp)

### Why?

Package `superhttp` is designed as a minimal, zero-dependency shim that complements and extends the [Go 1.22+ `http.ServeMux`](https://go.dev/blog/routing-enhancements), aiming to make structured routing and middleware composition more ergonomic. The aim is to have a developer experience more like [`chi`](https://github.com/go-chi/chi) or [`gin`](https://github.com/gin-gonic/gin), but with the underlying implementation still being the standard library.

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
