package middleware

import (
	"net/http"

	"github.com/sai-sy/linkShortener/internal/service/auth"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth.DefaultService == nil {
			http.Error(w, "auth service not configured", http.StatusInternalServerError)
			return
		}
		if _, err := auth.DefaultService.Authenticate(r.Context(), r); err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
