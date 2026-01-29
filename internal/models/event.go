package models

import "time"

type Event struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	TotalSeats     int    `json:"total_seats"`
	AvailableSeats int    `json:"available_seats"`
}

type Booking struct {
	ID        int       `json:"id"`
	EventID   int       `json:"event_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateEventRequest struct {
	Name           string `json:"name"`
	TotalSeats     int    `json:"total_seats"`
}

type BookingRequest struct {
	EventID int    `json:"event_id"`
	UserID  string `json:"user_id"`
}

type BookingResponse struct {
	Message string `json:"message"`
	BookingID int `json:"booking_id"`
	RemainingSeats int `json:"remaining_seats"`
}

type EventResponse struct {
	EventId int `json:"event_id"`
	Name string `json:"name"`
	TotalSeats int `json:"total_seats"`
	AvailableSeats int `json:"available_seats"`
	BookedCount int `json:"booked_count"`
}