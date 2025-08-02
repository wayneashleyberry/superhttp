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
		}))
	})

	http.ListenAndServe(":8080", r)
}
```
