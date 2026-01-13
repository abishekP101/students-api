package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/abishekP101/students-api/internal/postgres"
	"github.com/abishekP101/students-api/internal/storage"
	"github.com/abishekP101/students-api/internal/types"
	"github.com/abishekP101/students-api/internal/utils/response"

	"strconv"

	"github.com/go-playground/validator/v10"
)

func New(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student

		if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
			if errors.Is(err, io.EOF) {
				response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
				return
			}
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(student); err != nil {
			validateErrors := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrors))
			return
		}

		// save to database
		id, err := store.CreateStudent(
			r.Context(),
			student.Name,
			student.Email,
			student.Age,
		)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": id})
	}
}


func GetById(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := r.PathValue("id")
		slog.Info("Getting a student", slog.String("id", idParam))

		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			response.WriteJson(
				w,
				http.StatusBadRequest,
				response.GeneralError(errors.New("invalid student id")),
			)
			return
		}

		student, err := store.GetStudentById(r.Context(), id)
		if err != nil {
			// Optional: if you're using pgx, you can check pgx.ErrNoRows here
			response.WriteJson(
				w,
				http.StatusNotFound,
				response.GeneralError(errors.New("student not found")),
			)
			return
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}
func GetList(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Getting all students")

		students, err := store.GetList(r.Context())
		if err != nil {
			response.WriteJson(
				w,
				http.StatusInternalServerError,
				response.GeneralError(err),
			)
			return
		}

		response.WriteJson(w, http.StatusOK, students)
	}
}

func DeleteById(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := r.PathValue("id")
		slog.Info("Deleting a student", slog.String("id", idParam))

		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			response.WriteJson(
				w,
				http.StatusBadRequest,
				response.GeneralError(errors.New("invalid student id")),
			)
			return
		}

		err = store.DeleteStudentById(r.Context(), id)
		if err != nil {
			if errors.Is(err, postgres.ErrStudentNotFound) {
				response.WriteJson(
					w,
					http.StatusNotFound,
					response.GeneralError(err),
				)
				return
			}

			response.WriteJson(
				w,
				http.StatusInternalServerError,
				response.GeneralError(err),
			)
			return
		}

		// 204 No Content is standard for DELETE
		w.WriteHeader(http.StatusNoContent)
	}
}
