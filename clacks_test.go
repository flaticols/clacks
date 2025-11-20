package clacks_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flaticols/clacks"
)

func TestMiddleware(t *testing.T) {
	// Create a simple handler that returns 200 OK
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	// Wrap it with the Clacks middleware
	wrappedHandler := clacks.Clacks(handler)

	// Create a test request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// Serve the request
	wrappedHandler.ServeHTTP(rec, req)

	// Check the response status
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	// Check the X-Clacks-Overhead header is present
	header := rec.Header().Get(clacks.HeaderName)
	if header != clacks.HeaderValue {
		t.Errorf("expected header %q to be %q, got %q", clacks.HeaderName, clacks.HeaderValue, header)
	}

	// Check the response body
	body, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	expectedBody := "Hello, World!"
	if string(body) != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, string(body))
	}
}

func TestMiddlewareWithServeMux(t *testing.T) {
	// Create a ServeMux with multiple routes
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	mux.HandleFunc("/world", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("world"))
	})

	// Wrap the entire mux with the middleware
	wrappedMux := clacks.Clacks(mux)

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "hello route",
			path:     "/hello",
			expected: "hello",
		},
		{
			name:     "world route",
			path:     "/world",
			expected: "world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()

			wrappedMux.ServeHTTP(rec, req)

			// Check the header is present on all routes
			header := rec.Header().Get(clacks.HeaderName)
			if header != clacks.HeaderValue {
				t.Errorf("expected header %q to be %q, got %q", clacks.HeaderName, clacks.HeaderValue, header)
			}

			// Check the response body
			body, err := io.ReadAll(rec.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}

			if string(body) != tt.expected {
				t.Errorf("expected body %q, got %q", tt.expected, string(body))
			}
		})
	}
}

func TestMiddlewarePreservesExistingHeaders(t *testing.T) {
	// Create a handler that sets its own headers
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Custom-Header", "custom-value")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"ok"}`))
	})

	wrappedHandler := clacks.Clacks(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rec, req)

	// Check the Clacks header
	if rec.Header().Get(clacks.HeaderName) != clacks.HeaderValue {
		t.Errorf("Clacks header not set correctly")
	}

	// Check existing headers are preserved
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type header not preserved")
	}

	if rec.Header().Get("X-Custom-Header") != "custom-value" {
		t.Errorf("Custom header not preserved")
	}
}
