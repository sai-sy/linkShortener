package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sai-sy/linkShortener/internal/db"
)

func (h *Handler) CreateRoutemapHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	destination := r.FormValue("destination")
	if destination == "" {
		http.Error(w, "destination is required", http.StatusBadRequest)
		return
	}

	_, session := h.buildPageData(r, "index")
	profileID, ok := sessionProfileID(session)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	queries, commit, rollback, err := h.withProfileContext(r.Context(), profileID)
	if err != nil {
		http.Error(w, "failed to start transaction", http.StatusInternalServerError)
		return
	}

	workspace, err := queries.GetWorkspaceByProfile(r.Context(), pgtype.UUID{Bytes: profileID, Valid: true})
	if err != nil {
		_ = rollback()
		http.Error(w, "failed to load workspace", http.StatusInternalServerError)
		return
	}

	if err := queries.InsertRoutemap(r.Context(), db.InsertRoutemapParams{
		Destination: destination,
		WorkspaceID: workspace.ID,
	}); err != nil {
		_ = rollback()
		http.Error(w, "failed to create routemap", http.StatusInternalServerError)
		return
	}

	if err := commit(); err != nil {
		http.Error(w, "failed to commit", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
