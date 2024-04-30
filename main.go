package main

import (
	"database/sql"
	"ecom/pkg/server"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load("./.env")
	connStr := os.Getenv("POSTGRES_CONN")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Cannot connect to the database: %v", err)
	}
	server := server.NewServer(":3021", db)
	server.CreateTables()
	server.RegisterRoutes()
	if err := server.Start(); err != nil {
		log.Fatalf("Could not start the server: %v", err)
	}
}
