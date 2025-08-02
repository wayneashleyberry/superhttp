package superhttp

import (
	"context"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

type ServeMux struct {
	mux        *http.ServeMux
	middleware []Middleware
	prefix     string
}

type contextKey string

const RoutePatternKey = contextKey("route-pattern")

func NewServeMux() *ServeMux {
	return &ServeMux{
		mux: http.NewServeMux(),
	}
}

func (r *ServeMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// methodHandler wraps the handler to set the route pattern in context
func methodHandler(pattern string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), RoutePatternKey, pattern)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (r *ServeMux) handle(method string, pattern string, handler http.Handler) {
	fullPattern := r.prefix + pattern
	wrapped := methodHandler(fullPattern, applyMiddleware(handler, r.middleware...))
	r.mux.Handle(method+" "+fullPattern, wrapped)
}

// Type-safe HTTP methods
func (r *ServeMux) GET(pattern string, handler http.Handler) {
	r.handle("GET", pattern, handler)
}

func (r *ServeMux) POST(pattern string, handler http.Handler) {
	r.handle("POST", pattern, handler)
}

func (r *ServeMux) PUT(pattern string, handler http.Handler) {
	r.handle("PUT", pattern, handler)
}

func (r *ServeMux) DELETE(pattern string, handler http.Handler) {
	r.handle("DELETE", pattern, handler)
}

// Add middleware to current router
func (r *ServeMux) Use(mw ...Middleware) {
	r.middleware = append(r.middleware, mw...)
}

// Create a route group with prefix and middleware
func (r *ServeMux) Group(prefix string, fn func(gr *ServeMux)) {
	group := &ServeMux{
		mux:        r.mux,
		prefix:     r.prefix + prefix,
		middleware: append([]Middleware{}, r.middleware...),
	}
	fn(group)
}

func applyMiddleware(h http.Handler, middleware ...Middleware) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

// Helper for getting route pattern from context
func RoutePattern(r *http.Request) string {
	if val, ok := r.Context().Value(RoutePatternKey).(string); ok {
		return val
	}
	return ""
}
