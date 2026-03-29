package handlers

import (
	"net/http"
)

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		title := "404"
		p, err := loadPage(title)
		if err != nil {
			p = &Page{Title: title, Body: []byte("cannot find")}
		}
		w.WriteHeader(http.StatusNotFound)
		renderTemplate(w, title, p)
		return
	}

	title := "index"
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title, Body: []byte("cannot find")}
	}
	renderTemplate(w, title, p)
}
