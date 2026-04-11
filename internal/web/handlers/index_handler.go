package handlers

import (
	"net/http"
)

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		title := "404"
		p, _ := h.buildPageData(r, title)
		w.WriteHeader(http.StatusNotFound)
		renderTemplate(w, title, p)
		return
	}

	title := "index"
	p, _ := h.buildPageData(r, title)
	renderTemplate(w, title, p)
}
