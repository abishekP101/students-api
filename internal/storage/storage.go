package storage

import (
	"context"

	"github.com/abishekP101/students-api/internal/postgres"
	"github.com/abishekP101/students-api/internal/types"
)

type PostgresStorage struct {
	DB *postgres.Postgres
}
type Storage interface {
	CreateStudent(ctx context.Context, name, email string, age int) (int64, error)
	GetStudentById(ctx context.Context , id int64) (types.Student , error)
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



func (s *PostgresStorage) GetStudentById(
	ctx context.Context,
	id int64,
) (types.Student, error) {

	var student types.Student

	err := s.DB.DB.QueryRow(
		ctx,
		`SELECT id, name, email, age FROM students WHERE id = $1`,
		id,
	).Scan(
		&student.Id,
		&student.Name,
		&student.Email,
		&student.Age,
	)

	return student, err
}
