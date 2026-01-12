package storage

import (
	"context"

	"github.com/abishekP101/students-api/internal/postgres"
)

type PostgresStorage struct {
	DB *postgres.Postgres
}
type Storage interface {
	CreateStudent(ctx context.Context, name, email string, age int) (int64, error)
}

func NewPostgres(db *postgres.Postgres) *PostgresStorage {
	return &PostgresStorage{DB: db}
}

func (s *PostgresStorage) CreateStudent(
	ctx context.Context,
	name, email string,
	age int,
) (int64, error) {

	var id int64
	err := s.DB.DB.QueryRow(
		ctx,
		`INSERT INTO students(name, email, age) VALUES($1, $2, $3) RETURNING id`,
		name,
		email,
		age,
	).Scan(&id)

	return id, err
}


