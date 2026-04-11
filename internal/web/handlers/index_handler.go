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
	p, session := h.buildPageData(r, title)
	if session != nil {
		if profileID, ok := sessionProfileID(session); ok {
			queries, commit, rollback, err := h.withProfileContext(r.Context(), profileID)
			if err == nil {
				routemaps, queryErr := queries.GetAllRoutemaps(r.Context())
				if queryErr != nil {
					_ = rollback()
				} else {
					p.Routemaps = routemaps
					_ = commit()
				}
			}
		}
	}
	renderTemplate(w, title, p)
}
