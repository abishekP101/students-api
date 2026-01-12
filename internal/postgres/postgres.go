package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/abishekP101/students-api/internal/config"
)

type Postgres struct {
	DB *pgxpool.Pool
}

func New(cfg *config.Config) (*Postgres, error) {
	db, err := pgxpool.New(context.Background(), cfg.Storage.DSN)
	if err != nil {
		return nil, err
	}

	// verify connection
	if err := db.Ping(context.Background()); err != nil {
		return nil, err
	}

	// create table
	_, err = db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS students (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			age INTEGER,
			email TEXT UNIQUE
		)
	`)
	if err != nil {
		return nil, err
	}

	return &Postgres{
		DB: db,
	}, nil
}
