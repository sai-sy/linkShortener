package middleware

import (
	"log"
	"net/http"

	"github.com/sai-sy/linkShortener/internal/service/auth"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("auth middleware entered: %s %s", r.Method, r.URL.Path)
		if auth.DefaultService == nil {
			http.Error(w, "auth service not configured", http.StatusInternalServerError)
			return
		}
		if _, err := auth.DefaultService.Authenticate(r.Context(), r); err != nil {
			log.Printf("auth failed: %s %s: %v", r.Method, r.URL.Path, err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
