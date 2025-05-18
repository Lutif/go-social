package config

import (
	"time"

	"github.com/lutif/go-social/internal/db"
)

type Config struct {
	Addr         string
	DB           db.DbConfig
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	IdleTimeout  time.Duration
}

var DefaultConfig = Config{
	Addr:         ":8081",
	WriteTimeout: time.Second * 30,
	ReadTimeout:  time.Second * 10,
	IdleTimeout:  time.Minute * 1,

	DB: db.DbConfig{
		Addr:         "postgres://admin:password@localhost/social?sslmode=disable",
		MaxOpenConns: 30,
		MaxIdlConns:  30,
		MaxIdlTime:   "5m",
	},
}
