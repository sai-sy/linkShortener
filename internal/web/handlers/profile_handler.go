package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
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

	profileID, err := strconv.ParseInt(profileIDStr, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	p, session := h.buildPageData(r, "profile")
	if session == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if session.UserID != profileID {
		err := errors.New("forbidden")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	renderTemplate(w, "profile", p)
}
