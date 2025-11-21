package postgres

import (
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

type Config struct {
	Hostname string
	Port     uint64
	Database string
	User     string
	Password string
}

func NewPostgres(cfg Config) *sql.DB {
	const driver = "pgx"

	pattern := fmt.Sprintf(
		"host=%s port=%d database=%s user=%s password=%s sslmode=disable",
		cfg.Hostname, cfg.Port, cfg.Database, cfg.User, cfg.Password,
	)

	config, err := pgx.ParseConfig(pattern)
	if err != nil {
		panic(err)
	}

	connection := stdlib.RegisterConnConfig(config)
	db, err := sql.Open(driver, connection)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return db
}
