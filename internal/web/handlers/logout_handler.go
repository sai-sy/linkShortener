package handlers

import (
	"net/http"
	"time"
)

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if csrfCookie, err := r.Cookie("csrf_token"); err == nil {
		r.Header.Set("X-CSRF-Token", csrfCookie.Value)
	}

	if _, err := h.auth.Authenticate(r.Context(), r); err == nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "csrf_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			Expires:  time.Unix(0, 0),
			SameSite: http.SameSiteLaxMode,
		})

		p := &Page{Title: "logout"}
		renderTemplate(w, "logout", p)
		return
	}

	p := &Page{Title: "logout_guest"}
	renderTemplate(w, "logout_guest", p)
}
