package handlers

import (
	"net/http"
	"unicode"

	"github.com/jackc/pgx/v5/pgtype"
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

		routemap, err := h.db.GetRoutemap(r.Context(), pgtype.Text{String: path, Valid: true})
		if err != nil {
			title := "404"
			p, _ := h.buildPageData(r, title)
			w.WriteHeader(http.StatusNotFound)
			renderTemplate(w, title, p)
			return
		}

		http.Redirect(w, r, routemap.Destination, http.StatusFound)
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

func isBase62(value string) bool {
	for _, r := range value {
		if !unicode.IsDigit(r) && !(r >= 'a' && r <= 'z') && !(r >= 'A' && r <= 'Z') {
			return false
		}
	}
	return true
}
