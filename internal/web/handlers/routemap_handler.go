package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

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

	if !hasScheme(destination) {
		destination = "https://" + destination
	}

	_, session := h.buildPageData(r, "index")
	profileID, ok := sessionProfileID(session)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	queries, commit, rollback, err := h.withProfileContext(r.Context(), profileID)
	if err != nil {
		log.Printf("create routemap: start transaction failed: %v", err)
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
		log.Printf("create routemap: commit failed: %v", err)
		http.Error(w, "failed to commit", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) UpdateRoutemapHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	idParam := r.FormValue("id")
	destination := r.FormValue("destination")
	if idParam == "" || destination == "" {
		http.Error(w, "id and destination are required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if !hasScheme(destination) {
		destination = "https://" + destination
	}

	_, session := h.buildPageData(r, "index")
	profileID, ok := sessionProfileID(session)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	queries, commit, rollback, err := h.withProfileContext(r.Context(), profileID)
	if err != nil {
		log.Printf("update routemap: start transaction failed: %v", err)
		http.Error(w, "failed to start transaction", http.StatusInternalServerError)
		return
	}

	if err := queries.UpdateRoutemapDestination(r.Context(), db.UpdateRoutemapDestinationParams{
		ID:          id,
		Destination: destination,
	}); err != nil {
		_ = rollback()
		http.Error(w, "failed to update routemap", http.StatusInternalServerError)
		return
	}

	if err := commit(); err != nil {
		log.Printf("update routemap: commit failed: %v", err)
		http.Error(w, "failed to commit", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func hasScheme(value string) bool {
	return strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://")
}
