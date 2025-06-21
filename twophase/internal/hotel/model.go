package hotel

import (
	"time"
)

// Hotel represents a hotel entity
type Hotel struct {
	ID          string    `firestore:"id"`
	Name        string    `firestore:"name"`
	Location    string    `firestore:"location"`
	Description string    `firestore:"description"`
	Rating      float64   `firestore:"rating"`
	CreatedAt   time.Time `firestore:"created_at"`
	UpdatedAt   time.Time `firestore:"updated_at"`
}

// Room represents a hotel room
type Room struct {
	ID        string  `firestore:"id"`
	HotelID   string  `firestore:"hotel_id"`
	Number    string  `firestore:"number"`
	Type      string  `firestore:"type"`
	Capacity  int     `firestore:"capacity"`
	Price     float64 `firestore:"price"`
	Available bool    `firestore:"available"`
}

// Reservation represents a hotel reservation
type Reservation struct {
	ID           string    `firestore:"id"`
	OrderID      string    `firestore:"order_id"`
	HotelID      string    `firestore:"hotel_id"`
	RoomID       string    `firestore:"room_id"`
	UserID       string    `firestore:"user_id"`
	CheckInDate  time.Time `firestore:"check_in_date"`
	CheckOutDate time.Time `firestore:"check_out_date"`
	Status       string    `firestore:"status"` // "pending", "confirmed", "cancelled"
	TotalPrice   float64   `firestore:"total_price"`
	CreatedAt    time.Time `firestore:"created_at"`
	UpdatedAt    time.Time `firestore:"updated_at"`
}

// TwoPhaseTransaction represents a two-phase commit transaction for hotel
type TwoPhaseTransaction struct {
	ID            string    `firestore:"id"`
	TransactionID string    `firestore:"transaction_id"`
	OrderID       string    `firestore:"order_id"`
	Status        string    `firestore:"status"` // "prepared", "committed", "aborted"
	ReservationID string    `firestore:"reservation_id,omitempty"`
	CreatedAt     time.Time `firestore:"created_at"`
	UpdatedAt     time.Time `firestore:"updated_at"`
}

// CreateReservationRequest represents request to create reservation
type CreateReservationRequest struct {
	OrderID      string    `json:"order_id"`
	HotelID      string    `json:"hotel_id"`
	RoomID       string    `json:"room_id"`
	UserID       string    `json:"user_id"`
	CheckInDate  time.Time `json:"check_in_date"`
	CheckOutDate time.Time `json:"check_out_date"`
	TotalPrice   float64   `json:"total_price"`
}

// PrepareRequest represents prepare phase request
type PrepareRequest struct {
	TransactionID string `json:"transaction_id"`
	OrderID       string `json:"order_id"`
	ServiceName   string `json:"service_name"`
}

// PrepareResponse represents prepare phase response
type PrepareResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CommitRequest represents commit phase request
type CommitRequest struct {
	TransactionID string `json:"transaction_id"`
	OrderID       string `json:"order_id"`
	ServiceName   string `json:"service_name"`
}

// CommitResponse represents commit phase response
type CommitResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// AbortRequest represents abort phase request
type AbortRequest struct {
	TransactionID string `json:"transaction_id"`
	OrderID       string `json:"order_id"`
	ServiceName   string `json:"service_name"`
	Reason        string `json:"reason"`
}

// AbortResponse represents abort phase response
type AbortResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
