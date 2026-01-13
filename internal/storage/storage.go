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
	GetList(ctx context.Context) ([]types.Student , error)
	DeleteStudentById(ctx context.Context, id int64) error


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


func (s *PostgresStorage) GetList(
	ctx context.Context,
) ([]types.Student, error) {

	rows, err := s.DB.DB.Query(
		ctx,
		`SELECT id, name, email, age FROM students`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []types.Student

	for rows.Next() {
		var student types.Student
		err := rows.Scan(
			&student.Id,
			&student.Name,
			&student.Email,
			&student.Age,
		)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
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

func (s *PostgresStorage) DeleteStudentById(
	ctx context.Context,
	id int64,
) error {

	cmdTag, err := s.DB.DB.Exec(
		ctx,
		`DELETE FROM students WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	// Optional: check if a row was actually deleted
	if cmdTag.RowsAffected() == 0 {
		return postgres.ErrStudentNotFound // or define your own error
	}

	return nil
}