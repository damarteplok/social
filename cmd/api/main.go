package main

import (
	"log"

	"github.com/damarteplok/social/internal/db"
	"github.com/damarteplok/social/internal/env"
	"github.com/damarteplok/social/internal/store"
)

const version = "0.0.1"

func main() {
	cfg := config{
		addr: env.Envs.Addr,
		db: dbConfig{
			addr:         env.Envs.DbAddr,
			maxOpenConns: env.Envs.MaxOpenConns,
			maxIdleConns: env.Envs.MaxIdleConns,
			maxIdleTime:  env.Envs.MaxIdleTime,
		},
		env: env.Envs.ENV,
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	log.Println("database connection pool established")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
