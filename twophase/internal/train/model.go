package train

import (
	"time"
)

// Train represents a train entity
type Train struct {
	ID            string    `firestore:"id"`
	Name          string    `firestore:"name"`
	Route         string    `firestore:"route"`
	Capacity      int       `firestore:"capacity"`
	DepartureTime time.Time `firestore:"departure_time"`
	ArrivalTime   time.Time `firestore:"arrival_time"`
	Price         float64   `firestore:"price"`
	CreatedAt     time.Time `firestore:"created_at"`
	UpdatedAt     time.Time `firestore:"updated_at"`
}

// Seat represents a train seat
type Seat struct {
	ID        string  `firestore:"id"`
	TrainID   string  `firestore:"train_id"`
	Number    string  `firestore:"number"`
	Class     string  `firestore:"class"` // "economy", "business", "first"
	Price     float64 `firestore:"price"`
	Available bool    `firestore:"available"`
}

// Reservation represents a train reservation
type Reservation struct {
	ID         string    `firestore:"id"`
	OrderID    string    `firestore:"order_id"`
	TrainID    string    `firestore:"train_id"`
	SeatID     string    `firestore:"seat_id"`
	UserID     string    `firestore:"user_id"`
	TravelDate time.Time `firestore:"travel_date"`
	Status     string    `firestore:"status"` // "pending", "confirmed", "cancelled"
	TotalPrice float64   `firestore:"total_price"`
	CreatedAt  time.Time `firestore:"created_at"`
	UpdatedAt  time.Time `firestore:"updated_at"`
}

// TwoPhaseTransaction represents a two-phase commit transaction for train
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
	OrderID    string    `json:"order_id"`
	TrainID    string    `json:"train_id"`
	SeatID     string    `json:"seat_id"`
	UserID     string    `json:"user_id"`
	TravelDate time.Time `json:"travel_date"`
	TotalPrice float64   `json:"total_price"`
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
