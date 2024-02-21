package router

import "net/http"

// HealthCheckHandler returns an http.Handler that handles health check requests.
// It sets the Content-Type header to "text/plain; charset=utf-8" and writes "OK" as the response body.
func HealthCheckHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}
