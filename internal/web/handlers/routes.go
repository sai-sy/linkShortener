package handlers

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/sai-sy/linkShortener/internal/db"
	"github.com/sai-sy/linkShortener/internal/service/auth"
)

type Handler struct {
	auth *auth.Service
	db   *db.Queries
	conn *pgx.Conn
}

type Route struct {
	Pattern string
	Handler http.HandlerFunc
}

func NewHandler(db *db.Queries, conn *pgx.Conn) *Handler {
	return &Handler{
		auth: auth.NewService(db),
		db:   db,
		conn: conn,
	}
}

func (h *Handler) PrivateRoutes() []Route {
	return []Route{
		{Pattern: "POST /routemap", Handler: h.CreateRoutemapHandler},
		{Pattern: "POST /routemap/update", Handler: h.UpdateRoutemapHandler},
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
