package handlers

import (
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) ListWorkspaceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	p, session := h.buildPageData(r, "workspace")
	profileID, ok := sessionProfileID(session)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	page := 1
	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		if parsed, err := strconv.Atoi(pageParam); err == nil && parsed > 0 {
			page = parsed
		}
	}

	workspaces, err := h.workspaceSvc.List(r.Context(), profileID, page)
	if err != nil {
		log.Printf("list workspace: %v", err)
		http.Error(w, "failed to load workspaces", http.StatusInternalServerError)
		return
	}

	p.Workspaces = workspaces
	p.Page = page
	renderTemplate(w, "workspace", p)
}
