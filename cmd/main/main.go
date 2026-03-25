package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sai-sy/linkShortener/internal/db"
	"github.com/sai-sy/linkShortener/internal/httpServer"
)

func connect(ctx context.Context) *db.Queries {
	database_url := os.Getenv("DATABASE_URL")
	fmt.Println(database_url)

	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	database := db.New(conn)
	return database
}

func main() {
	ctx := context.Background()
	database := connect(ctx)
	server := httpServer.NewHTTPServer(":8080", database)
	if err := server.Run(); err != nil {
		fmt.Println("server error:", err)
		os.Exit(1)
	}
	fmt.Println("Listening on :8080")
}
