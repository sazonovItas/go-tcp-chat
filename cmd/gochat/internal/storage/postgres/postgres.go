package postgres

import (
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/config"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/internal/storage"
)

const (
	maxOpenConns        = 4
	maxConnLifetime     = time.Minute * 10
	maxIdleConns        = 2
	maxIdleConnLifetime = time.Minute
)

// TODO: Add optional database config
// TODO: Need to know more about conns and idle conns to database
func New(cfg *config.Storage) (*storage.Storage, error) {
	const op = "gochat.internal.storage.postgres.New"

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	connUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	// connect to database using pgx driver
	db, err := sqlx.Connect("pgx", connUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: error connecting to database: %s", op, err.Error())
	}

	// trying ping database
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: error ping database connection: %s", op, err.Error())
	}

	// Setting up db connections
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(maxConnLifetime)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(maxIdleConnLifetime)

	return &storage.Storage{DB: db}, nil
}
