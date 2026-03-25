package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	title := "index"
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title, Body: []byte("cannot find")}
	}
	fmt.Println("rendering page", title)
	renderTemplate(w, title, p)
}
