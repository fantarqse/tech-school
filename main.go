package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"tech-school/api"
	db "tech-school/db/sqlc"
	"tech-school/util"
)

func main() {
	cfg, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("failed to load the config: %v", err)
	}

	conn, err := sql.Open(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)
	if err := server.Run(cfg.Address); err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}
}
