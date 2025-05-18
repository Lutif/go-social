package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/lutif/go-social/internal/config"
	"github.com/lutif/go-social/internal/db"
	"github.com/lutif/go-social/internal/env"
	"github.com/lutif/go-social/internal/store"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Printf("Error loading .env files: %v", err)
	}

	cfg := config.Config{
		Addr:         env.GetString("ADDR", config.DefaultConfig.Addr),
		WriteTimeout: config.DefaultConfig.WriteTimeout,
		ReadTimeout:  config.DefaultConfig.ReadTimeout,
		IdleTimeout:  config.DefaultConfig.IdleTimeout,
		DB: db.DbConfig{
			Addr:         env.GetString("DB_URL", config.DefaultConfig.DB.Addr),
			MaxOpenConns: env.GetInt("MAX_OPEN_CONNS", config.DefaultConfig.DB.MaxOpenConns),
			MaxIdlConns:  env.GetInt("MAX_IDEAL_CONNS", config.DefaultConfig.DB.MaxIdlConns),
			MaxIdlTime:   env.GetString("MAX_IDEAL_TIME", config.DefaultConfig.DB.MaxIdlTime),
		},
	}

	database, err := db.New(cfg.DB)

	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	storage := store.NewPostgresStorage(database)

	app := &application{
		config: cfg,
		store:  storage,
	}

	log.Fatal(app.run())
}
