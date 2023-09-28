package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"tech-school/api"
	db "tech-school/db/sqlc"
)

const (
	dbDriver string = "postgres"
	dbSource string = "postgresql://root:secret@localhost:5436/simple_bank?sslmode=disable"
	addr     string = "0.0.0.0:8888"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)
	if err := server.Run(addr); err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}
}
