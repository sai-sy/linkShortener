package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sai-sy/linkShortener/internal/db"
)

func main() {
	ctx := context.Background()

	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@db:5432/linkshort?sslmode=disable")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hellooo World!\n")
	})

	http.HandleFunc("/routemap", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			path := r.URL.Query().Get("path")
			if path == "" {
				http.Error(w, "path is required", http.StatusBadRequest)
				return
			}

			routemap, err := queries.GetRoutemap(r.Context(), path)
			if err != nil {
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

	http.ListenAndServe(":8080", nil)
}
