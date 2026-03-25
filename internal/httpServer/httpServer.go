package httpServer

import (
	"net/http"

	apiv1 "github.com/sai-sy/linkShortener/internal/api/v1"
	"github.com/sai-sy/linkShortener/internal/db"
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

	apiV1SubMux := http.NewServeMux()
	apiV1Handler := apiv1.NewHandler()
	apiV1Handler.RegisterRoutes(apiV1SubMux)
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiV1SubMux))

	webHandler := handlers.NewHandler(s.db)
	webHandler.RegisterRoutes(mux)
	return http.ListenAndServe(s.addr, mux)
}
