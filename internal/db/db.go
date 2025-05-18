package db

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type DbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdlConns  int
	MaxIdlTime   string
}

func New(config DbConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.Addr)

	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(config.MaxIdlConns)
	db.SetMaxOpenConns(config.MaxOpenConns)

	idlTime, err := time.ParseDuration(config.MaxIdlTime)

	if err != nil {
		log.Panicf("Error parsing max ideal time")
	}
	db.SetConnMaxIdleTime(idlTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}
