package handlers

import (
	"log"
	"net/http"
	"strconv"
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

	if err := h.routemapSvc.Create(r.Context(), profileID, destination); err != nil {
		log.Printf("create routemap: %v", err)
		http.Error(w, "failed to create routemap", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/routemap", http.StatusSeeOther)
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

	_, session := h.buildPageData(r, "index")
	profileID, ok := sessionProfileID(session)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := h.routemapSvc.UpdateDestination(r.Context(), profileID, id, destination); err != nil {
		log.Printf("update routemap: %v", err)
		http.Error(w, "failed to update routemap", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/routemap", http.StatusSeeOther)
}

func (h *Handler) ListRoutemapHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	p, session := h.buildPageData(r, "routemap")
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

	routemaps, err := h.routemapSvc.List(r.Context(), profileID, page)
	if err != nil {
		log.Printf("list routemap: %v", err)
		http.Error(w, "failed to load routemaps", http.StatusInternalServerError)
		return
	}

	p.Routemaps = routemaps
	p.Page = page
	renderTemplate(w, "routemap", p)
}
