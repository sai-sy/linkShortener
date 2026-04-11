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
	password := r.FormValue("password")
	user, err := h.auth.Login(r.Context(), w, email, password)
	if err != nil {
		fmt.Printf("login error: %v\n", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	fmt.Printf("login submission: email=%s user_id=%v\n", user.Email, user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
