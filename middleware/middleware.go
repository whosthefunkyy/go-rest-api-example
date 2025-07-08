package middleware

import (
	"context"
	"net/http"
	"time"
)

func WithTimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			// Передаём новый контекст дальше по цепочке
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
