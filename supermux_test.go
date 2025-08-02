package superhttp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wayneashleyberry/superhttp"
)

func TestRoutePatternAndPathParam(t *testing.T) {
	mux := superhttp.NewServeMux()

	var routePattern string
	var paramValue string

	mux.GET("/hello/{name}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routePattern = superhttp.RoutePattern(r)
		paramValue = r.PathValue("name")
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/hello/world", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	if routePattern != "/hello/{name}" {
		t.Errorf("expected route pattern '/hello/{name}', got '%s'", routePattern)
	}

	if paramValue != "world" {
		t.Errorf("expected path param 'world', got '%s'", paramValue)
	}
}
