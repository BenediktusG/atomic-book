package main

import (
	"atomic-book/internal/database"
	"atomic-book/internal/handlers"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file (will run the app without using .env file):", err)
	}
	database.Connect()
	defer database.DB.Close()
	database.InitSchema()
	database.InitRedis()
	defer database.RedisClient.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /admin/event", handlers.CreateEvent)
	mux.HandleFunc("POST /book", handlers.BookEvent)
	mux.HandleFunc("GET /event/{id}", handlers.GetEvent)

	server := http.Server{
		Addr: ":8080",
		Handler: mux,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 10 * time.Second,
	}

	go func() {
		fmt.Println("AtomicBook is running")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down the server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exited properly")
}