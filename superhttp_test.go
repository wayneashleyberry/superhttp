package superhttp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wayneashleyberry/superhttp"
)

func TestRoutePatternStoredInContext(t *testing.T) {
	t.Parallel()

	router := superhttp.NewServeMux()
	expectedPattern := "/hello"

	router.GET(expectedPattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := superhttp.RoutePattern(r)
		if got != expectedPattern {
			t.Errorf("expected route pattern %q, got %q", expectedPattern, got)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Result().StatusCode)
	}
}

func TestPathParamsFromStandardLibrary(t *testing.T) {
	t.Parallel()

	router := superhttp.NewServeMux()

	router.GET("/users/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id != "42" {
			t.Errorf("expected path param id = 42, got %q", id)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/users/42", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Result().StatusCode)
	}
}

func TestMiddlewareIsApplied(t *testing.T) {
	t.Parallel()

	router := superhttp.NewServeMux()
	middlewareCalled := false

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middlewareCalled = true
			next.ServeHTTP(w, r)
		})
	})

	router.GET("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if !middlewareCalled {
		t.Errorf("expected middleware to be called")
	}
	if rec.Result().StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Result().StatusCode)
	}
}

func TestGroupScopedMiddleware(t *testing.T) {
	t.Parallel()

	router := superhttp.NewServeMux()
	var order []string

	router.Group("/api", func(api *superhttp.ServeMux) {
		api.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, "group")
				next.ServeHTTP(w, r)
			})
		})

		api.GET("/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "handler")
			w.WriteHeader(http.StatusOK)
		}))
	})

	req := httptest.NewRequest(http.MethodGet, "/api/hello", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if got, want := order, []string{"group", "handler"}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("expected middleware and handler order %v, got %v", want, got)
	}
	if rec.Result().StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Result().StatusCode)
	}
}
