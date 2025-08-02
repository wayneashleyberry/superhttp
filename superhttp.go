// Package superhttp provides a thin, idiomatic layer over Go's standard
// net/http, adding missing features such as type-safe HTTP method handlers,
// middleware chaining, route groups with scoped middleware, and request context
// enrichment with route metadata.
//
// It is designed as a minimal, zero-dependency shim that complements and
// extends the Go 1.22+ http.ServeMux, aiming to make structured routing and
// middleware composition more ergonomic.
//
// This package exists to bridge functionality that will hopefully be included
// in the Go standard library one day. If and when those capabilities are added,
// this package should become unnecessary.
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

// GET registers a handler for HTTP GET requests with the given pattern.
func (r *ServeMux) GET(pattern string, handlerFn http.HandlerFunc) {
	r.handle(http.MethodGet, pattern, handlerFn)
}

// POST registers a handler for HTTP POST requests with the given pattern.
func (r *ServeMux) POST(pattern string, handlerFn http.HandlerFunc) {
	r.handle(http.MethodPost, pattern, handlerFn)
}

// PUT registers a handler for HTTP PUT requests with the given pattern.
func (r *ServeMux) PUT(pattern string, handlerFn http.HandlerFunc) {
	r.handle(http.MethodPut, pattern, handlerFn)
}

// PATCH registers a handler for HTTP PATCH requests with the given pattern.
func (r *ServeMux) PATCH(pattern string, handlerFn http.HandlerFunc) {
	r.handle(http.MethodPatch, pattern, handlerFn)
}

// DELETE registers a handler for HTTP DELETE requests with the given pattern.
func (r *ServeMux) DELETE(pattern string, handlerFn http.HandlerFunc) {
	r.handle(http.MethodDelete, pattern, handlerFn)
}

// HEAD registers a handler for HTTP HEAD requests with the given pattern.
func (r *ServeMux) HEAD(pattern string, handlerFn http.HandlerFunc) {
	r.handle(http.MethodHead, pattern, handlerFn)
}

// OPTIONS registers a handler for HTTP OPTIONS requests with the given pattern.
func (r *ServeMux) OPTIONS(pattern string, handlerFn http.HandlerFunc) {
	r.handle(http.MethodOptions, pattern, handlerFn)
}

// Use registers middleware to be applied to all routes handled by this ServeMux.
func (r *ServeMux) Use(mw ...Middleware) {
	r.middleware = append(r.middleware, mw...)
}

// Group creates a new ServeMux with a prefix and applies the provided
// middleware to all routes within that group.
func (r *ServeMux) Group(prefix string, fnGroup func(gr *ServeMux)) {
	group := &ServeMux{
		mux:        r.mux,
		prefix:     r.prefix + prefix,
		middleware: slices.Clone(r.middleware),
	}
	fnGroup(group)
}

func (r *ServeMux) handle(method string, pattern string, handlerFn http.HandlerFunc) {
	fullPattern := r.prefix + pattern
	mwHandler := applyMiddleware(handlerFn, r.middleware...)
	wrapped := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := context.WithValue(req.Context(), RoutePatternKey, fullPattern)
		mwHandler.ServeHTTP(w, req.WithContext(ctx))
	})
	r.mux.Handle(method+" "+fullPattern, wrapped)
}

func applyMiddleware(h http.HandlerFunc, middleware ...Middleware) http.Handler {
	var handler http.Handler = h
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}
	return handler
}

// RoutePattern retrieves the route pattern from the request context.
func RoutePattern(r *http.Request) string {
	if val, ok := r.Context().Value(RoutePatternKey).(string); ok {
		return val
	}

	return ""
}
