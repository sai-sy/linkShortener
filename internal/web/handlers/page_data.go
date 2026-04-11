package handlers

import (
	"net/http"
	"strconv"

	"github.com/sai-sy/linkShortener/internal/db"
)

func (h *Handler) buildPageData(r *http.Request, title string) (*Page, *db.GetAuthSessionByTokenRow) {
	p := &Page{Title: title}

	session, err := h.auth.Authenticate(r.Context(), r)
	if err != nil {
		return p, nil
	}

	p.Authenticated = true
	p.ProfileURL = "/profile/" + strconv.FormatInt(session.UserID, 10)

	if csrfCookie, err := r.Cookie("csrf_token"); err == nil {
		p.CSRFToken = csrfCookie.Value
	}

	return p, &session
}

func sessionProfileID(session *db.GetAuthSessionByTokenRow) (int64, bool) {
	if session == nil {
		return 0, false
	}

	return session.UserID, true
}
