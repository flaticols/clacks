// Package clacks provides HTTP middleware for adding the X-Clacks-Overhead header
// to commemorate Terry Pratchett through the GNU Terry Pratchett protocol.
//
// See https://xclacksoverhead.org/ for more information.
package clacks

import "net/http"

const (
	// HeaderName is the name of the Clacks Overhead header
	HeaderName = "X-Clacks-Overhead"

	// HeaderValue is the standard GNU Terry Pratchett value
	HeaderValue = "GNU Terry Pratchett"
)

// Clacks adds the X-Clacks-Overhead header to all HTTP responses.
// It wraps the provided handler and ensures the header is set before
// calling the next handler in the chain.
//
// Example usage with http.ServeMux:
//
//	mux := http.NewServeMux()
//	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//		w.Write([]byte("Hello, World!"))
//	})
//
//	handler := clacks.Clacks(mux)
//	http.ListenAndServe(":8080", handler)
func Clacks(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(HeaderName, HeaderValue)
		next.ServeHTTP(w, r)
	})
}
