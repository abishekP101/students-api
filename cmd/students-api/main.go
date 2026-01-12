package main

import (
	"context"
	"fmt"
	"github/abishekP101/students-api/internal/config"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustLoad()

	router := http.NewServeMux()

	router.HandleFunc("GET /" , func(w http.ResponseWriter , r *http.Request) {
		w.Write([]byte("Welcome to the students api"))

	})

	server := http.Server {
		Addr : cfg.HTTPServer.ADDRESS,
		Handler: router,
	}

	slog.Info("Server started ",slog.String("Address" , cfg.ADDRESS))
	done := make(chan os.Signal , 1)

	signal.Notify(done , os.Interrupt , syscall.SIGINT , syscall.SIGTERM)


	go func(){
		fmt.Println("Server started")
		err := server.ListenAndServe()
		if err != nil{
			log.Fatal("Failed to start server")
	}
	}()

	<-done

	slog.Info("Shutting down the server.....")
	ctx , cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err:= server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server" , slog.String("error" , err.Error()))
	}

	slog.Info("Server shutdown Successfully")

}