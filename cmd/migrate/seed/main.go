package main

import (
	"log"

	"github.com/damarteplok/social/internal/db"
	"github.com/damarteplok/social/internal/env"
	"github.com/damarteplok/social/internal/store"
)

func main() {
	addr := env.Envs.DbAddr
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	store := store.NewStorage(conn)
	db.Seed(store, conn)
}
