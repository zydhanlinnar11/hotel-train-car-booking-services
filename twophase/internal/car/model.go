package car

import (
	"time"
)

// Car represents a car entity
type Car struct {
	ID           string    `firestore:"id"`
	Brand        string    `firestore:"brand"`
	Model        string    `firestore:"model"`
	Year         int       `firestore:"year"`
	LicensePlate string    `firestore:"license_plate"`
	Color        string    `firestore:"color"`
	Price        float64   `firestore:"price"`
	Available    bool      `firestore:"available"`
	CreatedAt    time.Time `firestore:"created_at"`
	UpdatedAt    time.Time `firestore:"updated_at"`
}

// Reservation represents a car reservation
type Reservation struct {
	ID         string    `firestore:"id"`
	OrderID    string    `firestore:"order_id"`
	CarID      string    `firestore:"car_id"`
	UserID     string    `firestore:"user_id"`
	StartDate  time.Time `firestore:"start_date"`
	EndDate    time.Time `firestore:"end_date"`
	Status     string    `firestore:"status"` // "pending", "confirmed", "cancelled"
	TotalPrice float64   `firestore:"total_price"`
	CreatedAt  time.Time `firestore:"created_at"`
	UpdatedAt  time.Time `firestore:"updated_at"`
}

// TwoPhaseTransaction represents a two-phase commit transaction for car
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
	CarID      string    `json:"car_id"`
	UserID     string    `json:"user_id"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
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
