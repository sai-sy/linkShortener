package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")
		if password != confirmPassword {
			http.Error(w, "passwords do not match", http.StatusBadRequest)
			return
		}

		user, err := h.auth.Register(r.Context(), w, email, password)
		if err != nil {
			fmt.Printf("register error: %v\n", err)
			http.Error(w, "failed to register", http.StatusBadRequest)
			return
		}
		fmt.Printf("register submission: email=%s user_id=%v\n", user.Email, user.ID)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	title := "register"
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title, Body: []byte("cannot find")}
	}
	fmt.Println("rendering page", title)
	renderTemplate(w, title, p)
}
