package httpServer

import (
	"net/http"

	"github.com/sai-sy/linkShortener/internal/db"
	"github.com/sai-sy/linkShortener/internal/middleware"
	"github.com/sai-sy/linkShortener/internal/service/auth"
	"github.com/sai-sy/linkShortener/internal/web"
	"github.com/sai-sy/linkShortener/internal/web/handlers"
)

type HTTPServer struct {
	addr string
	db   *db.Queries
}

func NewHTTPServer(addr string, db *db.Queries) *HTTPServer {
	return &HTTPServer{
		addr: addr,
		db:   db,
	}
}

func (s *HTTPServer) Run() error {
	mux := http.NewServeMux()

	staticRoot, err := web.StaticFS()
	if err != nil {
		return err
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticRoot))))

	publicWebHandler := handlers.NewHandler(s.db)
	publicWebHandler.RegisterRoutes(mux, publicWebHandler.PublicRoutes())

	privateWebHandler := handlers.NewHandler(s.db)
	auth.SetDefaultService(s.db)
	privateWebMiddleware := middleware.CreateStack(middleware.RequireAuth)
	for _, route := range privateWebHandler.PrivateRoutes() {
		mux.Handle(route.Pattern, privateWebMiddleware(http.HandlerFunc(route.Handler)))
	}
	middleware := middleware.CreateStack(
		middleware.Logging,
	)

	return http.ListenAndServe(s.addr, middleware(mux))
}
