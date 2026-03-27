package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handler) ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if _, err := h.auth.Authenticate(r.Context(), r); err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	title := "protected"
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title, Body: []byte("cannot find")}
	}
	fmt.Println("rendering page", title)
	renderTemplate(w, title, p)
}
