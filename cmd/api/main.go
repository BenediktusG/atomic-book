package main

import (
	"atomic-book/internal/database"
	"atomic-book/internal/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
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

	fmt.Println("AtomicBook is running")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}