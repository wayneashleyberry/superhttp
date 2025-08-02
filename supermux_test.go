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

	mux.GET("/user/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routePattern = superhttp.RoutePattern(r)
		paramValue = r.PathValue("id")
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/user/123", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	if routePattern != "/user/{id}" {
		t.Errorf("expected route pattern '/user/{id}', got '%s'", routePattern)
	}

	if paramValue != "123" {
		t.Errorf("expected path param '123', got '%s'", paramValue)
	}
}
