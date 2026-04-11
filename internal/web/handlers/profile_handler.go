package handlers

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (h *Handler) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !strings.HasPrefix(r.URL.Path, "/profile/") {
		http.NotFound(w, r)
		return
	}

	profileIDStr := strings.TrimPrefix(r.URL.Path, "/profile/")
	if profileIDStr == "" || strings.Contains(profileIDStr, "/") {
		http.NotFound(w, r)
		return
	}

	profileID, err := uuid.Parse(profileIDStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	p, session := h.buildPageData(r, "profile")
	if session == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !session.UserID.Valid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	sessionID, err := uuid.FromBytes(session.UserID.Bytes[:])
	if err != nil || sessionID != profileID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	renderTemplate(w, "profile", p)
}
