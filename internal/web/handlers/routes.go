package handlers

import (
	"net/http"

	"github.com/sai-sy/linkShortener/internal/db"
	"github.com/sai-sy/linkShortener/internal/service/auth"
)

type Handler struct {
	auth *auth.Service
}

func NewHandler(db *db.Queries) *Handler {
	return &Handler{auth: auth.NewService(db)}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.IndexHandler)
	mux.HandleFunc("/login", h.LoginHandler)
	mux.HandleFunc("/register", h.RegisterHandler)
}
