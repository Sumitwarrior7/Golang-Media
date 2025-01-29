package main

import (
	"log"

	"github.com/Sumitwarrior7/social/internal/db"
	"github.com/Sumitwarrior7/social/internal/env"
	"github.com/Sumitwarrior7/social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://sumit_user:pg_pass_key@localhost/social_network?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	// version := 1
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	store := store.NewPostgresStorage(conn)
	db.Seed(store, conn)
}
