package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect() {
	dbUrl := os.Getenv("DB_URL")
	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Fatal("Failed to parse Database config:", err)
	}

	config.MaxConns = 50
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour

	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.Ping(context.Background())
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Connected to PostgreSQL successfully")
}

func InitSchema() {
	query := `
		CREATE TABLE IF NOT EXISTS events (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			total_seats INT NOT NULL,
			available_seats INT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS bookings (
			id SERIAL PRIMARY KEY,
			event_id INT REFERENCES events(id) NOT NULL,
			user_id TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := DB.Exec(context.Background(), query)
	if err != nil {
		log.Fatal("Failed to create schema:", err)
	}

	fmt.Println("Database schema initialized.")
}