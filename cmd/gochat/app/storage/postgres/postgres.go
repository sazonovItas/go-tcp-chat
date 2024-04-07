package postgres

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/config"
	"github.com/sazonovItas/gochat-tcp/cmd/gochat/app/storage"
)

const (
	maxOpenConns        = 10
	maxConnLifetime     = time.Minute * 10
	maxIdleConns        = 5
	maxIdleConnLifetime = time.Minute
)

// TODO: Add optional postgres config
// TODO: Need to know more about conns and idle conns to database
func New(cfg *config.Storage) (*storage.Storage, error) {
	const op = "gochat.app.storage.postgres.New"

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

	// Setting up db connections
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(maxConnLifetime)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(maxIdleConnLifetime)

	return &storage.Storage{DB: db}, nil
}

func Migrate(cfg *config.Storage, pathToSqlFiles, dbName string) error {
	const op = "gochat.app.storage.postgres.Migrate"

	connUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db, err := sqlx.Connect("pgx", connUrl)
	if err != nil {
		return fmt.Errorf("%s: failed to connect %w", op, err)
	}

	driver, err := pgx.WithInstance(db.DB, &pgx.Config{})
	if err != nil {
		return fmt.Errorf("%s: with instance %w", op, err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+pathToSqlFiles, dbName, driver)
	if err != nil {
		return fmt.Errorf("%s: new with database instance %w", op, err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil {
		switch {
		case errors.Is(err, migrate.ErrNoChange):
			return nil
		default:
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}
