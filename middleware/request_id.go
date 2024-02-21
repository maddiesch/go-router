package middleware

import (
	"net/http"

	"github.com/oklog/ulid/v2"
)

func ULIDRequestID() string {
	return "req_" + ulid.Make().String()
}

func RequestID(provider ...func() string) func(http.Handler) http.Handler {
	fn := ULIDRequestID
	if len(provider) > 0 {
		fn = provider[0]
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := fn()

			r.Header.Set("X-Request-ID", id)

			next.ServeHTTP(w, r)
		})
	}
}
