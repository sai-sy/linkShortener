package routes

import (
    "fmt"
    "net/http"

    "github.com/sai-sy/linkShortener/internal/db"
)

func RegisterRoutes(mux *http.ServeMux, queries *db.Queries) {
    mux.HandleFunc("/api/routemap", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            path := r.URL.Query().Get("path")
            if path == "" {
                http.Error(w, "path is required", http.StatusBadRequest)
                return
            }

            routemap, err := queries.GetRoutemap(r.Context(), path)
            if err != nil {
                fmt.Println("woah error", err)
                http.Error(w, "failed to load route map", http.StatusInternalServerError)
                return
            }

            fmt.Printf("Route map: %+v\n", routemap)
            fmt.Fprintf(w, "Route map for %s -> %s\n", routemap.Path, routemap.Destination)

        case http.MethodPost:
            path := r.URL.Query().Get("path")
            destination := r.URL.Query().Get("destination")
            fmt.Println(path)
            fmt.Println(destination)
            if path == "" || destination == "" {
                http.Error(w, "path and destination are required", http.StatusBadRequest)
                return
            }

            if err := queries.InsertRoutemap(r.Context(), db.InsertRoutemapParams{
                Path:        path,
                Destination: destination,
            }); err != nil {
                http.Error(w, "failed to create route map", http.StatusInternalServerError)
                return
            }

            fmt.Fprintf(w, "Route map created for %s -> %s\n", path, destination)

        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })
}
