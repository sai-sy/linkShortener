package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/sai-sy/linkShortener/internal/db"
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

	queries, commit, rollback, err := h.withProfileContext(r.Context(), profileID)
	if err != nil {
		log.Printf("list workspace: start transaction failed: %v", err)
		http.Error(w, "failed to start transaction", http.StatusInternalServerError)
		return
	}

	page := 1
	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		if parsed, err := strconv.Atoi(pageParam); err == nil && parsed > 0 {
			page = parsed
		}
	}
	limit := int32(25)
	offset := int32((page - 1) * 25)

	workspaces, err := queries.GetWorkspacesPage(r.Context(), db.GetWorkspacesPageParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		_ = rollback()
		http.Error(w, "failed to load workspaces", http.StatusInternalServerError)
		return
	}
	if err := commit(); err != nil {
		log.Printf("list workspace: commit failed: %v", err)
		http.Error(w, "failed to commit", http.StatusInternalServerError)
		return
	}

	p.Workspaces = workspaces
	p.Page = page
	renderTemplate(w, "workspace", p)
}
