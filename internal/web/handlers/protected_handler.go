package handlers

import (
	"net/http"
)

func (h *Handler) ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	p, session := h.buildPageData(r, "protected")
	if session == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	renderTemplate(w, "protected", p)
}
