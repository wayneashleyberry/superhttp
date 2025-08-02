// Package superhttp provides a thin, idiomatic layer over Go's standard
// net/http, adding missing features such as type-safe HTTP method handlers,
// middleware chaining, route groups with scoped middleware, and request context
// enrichment with route metadata.
//
// It is designed as a minimal, zero-dependency shim that complements and
// extends the Go 1.22+ http.ServeMux, aiming to make structured routing and
// middleware composition more ergonomic.
//
// This package exists to bridge functionality that may eventually be included
// in the Go standard library. If and when those capabilities are added, this
// package should become unnecessary.
package superhttp

import (
	"context"
	"net/http"
	"slices"
)

// Middleware defines a function type that takes an http.Handler and returns
// another http.Handler.
type Middleware func(http.Handler) http.Handler

// ServeMux is a custom HTTP request multiplexer that supports method-specific
// routing, middleware chaining, and route grouping with prefixes.
type ServeMux struct {
	mux        *http.ServeMux
	middleware []Middleware
	prefix     string
}

type contextKey string

// RoutePatternKey is a context key used to store the route pattern in the request context.
const RoutePatternKey = contextKey("route-pattern")

// NewServeMux creates a new ServeMux instance with an empty prefix and no middleware.
func NewServeMux() *ServeMux {
	return &ServeMux{
		mux:        http.NewServeMux(),
		prefix:     "",
		middleware: nil,
	}
}

// ServeHTTP implements the http.Handler interface for ServeMux.
func (r *ServeMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func methodHandler(pattern string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), RoutePatternKey, pattern)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GET registers a handler for HTTP GET requests with the given pattern.
func (r *ServeMux) GET(pattern string, handler http.Handler) {
	r.handle(http.MethodGet, pattern, handler)
}

// POST registers a handler for HTTP POST requests with the given pattern.
func (r *ServeMux) POST(pattern string, handler http.Handler) {
	r.handle(http.MethodPost, pattern, handler)
}

// PUT registers a handler for HTTP PUT requests with the given pattern.
func (r *ServeMux) PUT(pattern string, handler http.Handler) {
	r.handle(http.MethodPut, pattern, handler)
}

// DELETE registers a handler for HTTP DELETE requests with the given pattern.
func (r *ServeMux) DELETE(pattern string, handler http.Handler) {
	r.handle(http.MethodDelete, pattern, handler)
}

// HEAD registers a handler for HTTP HEAD requests with the given pattern.
func (r *ServeMux) HEAD(pattern string, handler http.Handler) {
	r.handle(http.MethodHead, pattern, handler)
}

// OPTIONS registers a handler for HTTP OPTIONS requests with the given pattern.
func (r *ServeMux) OPTIONS(pattern string, handler http.Handler) {
	r.handle(http.MethodOptions, pattern, handler)
}

// Use registers middleware to be applied to all routes handled by this ServeMux.
func (r *ServeMux) Use(mw ...Middleware) {
	r.middleware = append(r.middleware, mw...)
}

// Group creates a new ServeMux with a prefix and applies the provided
// middleware to all routes within that group.
func (r *ServeMux) Group(prefix string, mux func(gr *ServeMux)) {
	group := &ServeMux{
		mux:        r.mux,
		prefix:     r.prefix + prefix,
		middleware: slices.Clone(r.middleware),
	}
	mux(group)
}

func (r *ServeMux) handle(method string, pattern string, handler http.Handler) {
	fullPattern := r.prefix + pattern
	wrapped := methodHandler(fullPattern, applyMiddleware(handler, r.middleware...))
	r.mux.Handle(method+" "+fullPattern, wrapped)
}

func applyMiddleware(h http.Handler, middleware ...Middleware) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}

	return h
}

// RoutePattern retrieves the route pattern from the request context.
func RoutePattern(r *http.Request) string {
	if val, ok := r.Context().Value(RoutePatternKey).(string); ok {
		return val
	}

	return ""
}
