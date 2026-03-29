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

	publicWebMux := http.NewServeMux()
	publicWebHandler := handlers.NewHandler(s.db)
	publicWebHandler.RegisterRoutes(publicWebMux, publicWebHandler.PublicRoutes())

	privateWebMux := http.NewServeMux()
	privateWebHandler := handlers.NewHandler(s.db)
	privateWebHandler.RegisterRoutes(privateWebMux, privateWebHandler.PrivateRoutes())
	auth.SetDefaultService(s.db)
	privateWebMiddleware := middleware.CreateStack(middleware.RequireAuth)

	mux.Handle("/", privateWebMiddleware(privateWebMux))
	mux.Handle("/", publicWebMux)
	middleware := middleware.CreateStack(
		middleware.Logging,
	)

	return http.ListenAndServe(s.addr, middleware(mux))
}
