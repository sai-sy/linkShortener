package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}
		email := r.FormValue("email")
		fmt.Printf("login submissionsss: email=%s\n", email)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	title := "login"
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title, Body: []byte("cannot find")}
	}
	fmt.Println("rendering page", title)
	renderTemplate(w, title, p)
}
