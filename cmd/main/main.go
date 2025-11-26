package main

import (
	"fmt"
	"net/http"
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close(ctx)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprint(w, "Hellooo World!\n")
	})
	
	http.ListenAndServe(":8080", nil)
}
