package handlers

import (
	"fmt"
	"net/http"

	"github.com/sai-sy/linkShortener/internal/db"
	"github.com/sai-sy/linkShortener/internal/service/auth"
)

type Handler struct {
	auth *auth.Service
}

type Route struct {
	Pattern string
	Handler http.HandlerFunc
}

func NewHandler(db *db.Queries) *Handler {
	return &Handler{auth: auth.NewService(db)}
}

func (h *Handler) PrivateRoutes() []Route {
	return []Route{
		{Pattern: "/protected", Handler: h.ProtectedHandler},
		{Pattern: "/profile/", Handler: h.ProfileHandler},
		{Pattern: "/logout", Handler: h.LogoutHandler},
	}
}

func (h *Handler) PublicRoutes() []Route {
	return []Route{
		{Pattern: "/", Handler: h.IndexHandler},
		{Pattern: "/hello-world", Handler: func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "hello world!!!\n") }},
		{Pattern: "GET /login", Handler: h.LoginHandlerGET},
		{Pattern: "POST /login", Handler: h.LoginHandlerPOST},
		{Pattern: "/register", Handler: h.RegisterHandler},
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, routes []Route) {
	for _, route := range routes {
		mux.HandleFunc(route.Pattern, route.Handler)
	}
}
