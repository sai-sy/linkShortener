package httpServer

import (
	"net/http"

	"github.com/sai-sy/linkShortener/internal/db"
	apiv1 "github.com/sai-sy/linkShortener/internal/service/api/v1"
	"github.com/sai-sy/linkShortener/internal/service/web"
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

	webHandler := web.NewHandler()
	webHandler.RegisterRoutes(mux)
	return http.ListenAndServe(s.addr, mux)
}
