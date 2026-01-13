package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abishekP101/students-api/internal/config"
	"github.com/abishekP101/students-api/internal/http/handlers/student"
	"github.com/abishekP101/students-api/internal/postgres"
	"github.com/abishekP101/students-api/internal/storage"
)

func main() {
	cfg := config.MustLoad()

	db, err := postgres.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.DB.Close() // if using pgxpool, use db.DB.Close()

	slog.Info("Storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	router := http.NewServeMux()
	store := storage.NewPostgres(db)           // create storage
	router.HandleFunc("POST /api/students", student.New(store))
	router.HandleFunc("GET /api/students/{id}" , student.GetById(store))
	router.HandleFunc("GET /api/students" , student.GetList(store))	
	router.HandleFunc("DELETE /api/students/{id}" , student.DeleteById(store))

		

	server := http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: router,
	}

	slog.Info("Server starting", slog.String("address", cfg.HTTPServer.Address))
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	<-done

	slog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")
}
