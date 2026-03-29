package httpServer

import (
	"net/http"

	"github.com/sai-sy/linkShortener/internal/db"
	"github.com/sai-sy/linkShortener/internal/middleware"
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

	webHandler := handlers.NewHandler(s.db)
	webHandler.RegisterRoutes(mux)

	middleware := middleware.CreateStack(
		middleware.Logging,
	)

	return http.ListenAndServe(s.addr, middleware(mux))
}
