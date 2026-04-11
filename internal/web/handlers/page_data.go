package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/sai-sy/linkShortener/internal/db"
)

func (h *Handler) buildPageData(r *http.Request, title string) (*Page, *db.GetAuthSessionByTokenRow) {
	p := &Page{Title: title}

	session, err := h.auth.Authenticate(r.Context(), r)
	if err != nil {
		return p, nil
	}

	p.Authenticated = true
	if session.UserID.Valid {
		userUUID, err := uuid.FromBytes(session.UserID.Bytes[:])
		if err == nil {
			p.ProfileURL = "/profile/" + userUUID.String()
		}
	}

	return p, &session
}
