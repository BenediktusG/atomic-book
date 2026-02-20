package handlers

import (
	"atomic-book/internal/database"
	"atomic-book/internal/models"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.TotalSeats < 1 {
		http.Error(w, "Total seats must be positive", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO events (name, total_seats, available_seats)
		VALUES ($1, $2, $2)
		RETURNING id, name, total_seats, available_seats
	`	
	var event models.Event
	err := database.DB.QueryRow(context.Background(), query, req.Name, req.TotalSeats).Scan(&event.ID, &event.Name, &event.TotalSeats, &event.AvailableSeats)
	if err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

func BookEvent(w http.ResponseWriter, r *http.Request) {
	var req models.BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.EventID < 1 {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	tx, err := database.DB.Begin(context.Background())
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(context.Background())
	var event models.Event
	err = tx.QueryRow(context.Background(), "SELECT available_seats FROM events WHERE id = $1 FOR UPDATE", req.EventID).Scan(&event.AvailableSeats)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	if event.AvailableSeats < 1 {
		http.Error(w, "Event is sold out", http.StatusBadRequest)
		return
	}

	newAvailableSeats := event.AvailableSeats - 1
	query := `
		UPDATE events
		SET available_seats = $1
		WHERE id = $2
		RETURNING available_seats
	`

	var bookingResponse models.BookingResponse
	err = tx.QueryRow(context.Background(), query, newAvailableSeats, req.EventID).Scan(&bookingResponse.RemainingSeats)
	if err != nil {
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	query = `
		INSERT INTO bookings (event_id, user_id)
		VALUES ($1, $2)
		RETURNING id
	`

	
	err = tx.QueryRow(context.Background(), query, req.EventID, req.UserID).Scan(&bookingResponse.BookingID)
	if err != nil {
		http.Error(w, "Failed to book event", http.StatusInternalServerError)
		return
	}

	cacheKey := "event:" + strconv.Itoa(req.EventID)
	database.RedisClient.Del(context.Background(), cacheKey)

	bookingResponse.Message = "Booking successful"

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bookingResponse)
	tx.Commit(context.Background())
}

func GetEvent(w http.ResponseWriter, r *http.Request) {
	eventID := r.PathValue("id")
	cacheKey := "event:" + eventID
	val, err := database.RedisClient.Get(context.Background(), cacheKey).Result()

	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache", "HIT")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(val))
		return
	}
	var eventResponse models.EventResponse
	query := `
		SELECT *
		FROM events
		WHERE id = $1;
	`
	err = database.DB.QueryRow(context.Background(), query, eventID).Scan(&eventResponse.EventId, &eventResponse.Name, &eventResponse.TotalSeats, &eventResponse.AvailableSeats)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	query = `
		SELECT COUNT(*)
		FROM bookings
		WHERE event_id = $1
	`

	err = database.DB.QueryRow(context.Background(), query, eventID).Scan(&eventResponse.BookedCount)
	if err != nil {
		http.Error(w, "Failed to get event", http.StatusInternalServerError)
		return
	}

	JSONBytes, _ := json.Marshal(eventResponse)

	database.RedisClient.Set(context.Background(), cacheKey, JSONBytes, 10*time.Second)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cache", "MISS")
	w.WriteHeader(http.StatusOK)
	w.Write(JSONBytes)
}