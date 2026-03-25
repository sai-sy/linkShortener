package handlers

import "net/http"

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.rootHandler)
}

func (h *Handler) rootHandler(w http.ResponseWriter, r *http.Request) {
}
