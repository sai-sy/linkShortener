package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handler) LoginHandlerGET(w http.ResponseWriter, r *http.Request) {
	title := "login"
	p, _ := h.buildPageData(r, title)
	renderTemplate(w, title, p)
}

func (h *Handler) LoginHandlerPOST(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	email := r.FormValue("email")
	fmt.Printf("login submissionsss: email=%s\n", email)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
