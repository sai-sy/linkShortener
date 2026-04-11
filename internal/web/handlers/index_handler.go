package handlers

import (
	"net/http"
	"strings"
	"unicode"
)

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		path := r.URL.Path
		if path[0] == '/' {
			path = path[1:]
		}
		if path == "" || !isBase62(path) {
			title := "404"
			p, _ := h.buildPageData(r, title)
			w.WriteHeader(http.StatusNotFound)
			renderTemplate(w, title, p)
			return
		}

		routemap, err := h.db.GetRoutemap(r.Context(), path)
		if err != nil {
			title := "404"
			p, _ := h.buildPageData(r, title)
			w.WriteHeader(http.StatusNotFound)
			renderTemplate(w, title, p)
			return
		}

		destination := routemap.Destination
		if !strings.HasPrefix(destination, "http://") && !strings.HasPrefix(destination, "https://") {
			destination = "https://" + destination
		}
		http.Redirect(w, r, destination, http.StatusFound)
		return
	}

	title := "index"
	p, _ := h.buildPageData(r, title)
	renderTemplate(w, title, p)
}

func isBase62(value string) bool {
	for _, r := range value {
		if !unicode.IsDigit(r) && !(r >= 'a' && r <= 'z') && !(r >= 'A' && r <= 'Z') {
			return false
		}
	}
	return true
}
