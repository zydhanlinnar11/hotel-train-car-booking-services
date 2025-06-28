package coordinator

import (
	"time"
)

// TransactionStatus represents the status of a two-phase commit transaction
type TransactionStatus string

const (
	StatusInitiated  TransactionStatus = "initiated"
	StatusPrepared   TransactionStatus = "prepared"
	StatusCommitted  TransactionStatus = "committed"
	StatusAborted    TransactionStatus = "aborted"
	StatusRolledBack TransactionStatus = "rolled_back"
	StatusTimedOut   TransactionStatus = "timed_out"
)

// TransactionLog represents a transaction log entry in Firestore
type TransactionLog struct {
	ID              string            `firestore:"id"`
	OrderID         string            `firestore:"order_id"`
	Status          TransactionStatus `firestore:"status"`
	Participants    []Participant     `firestore:"participants"`
	DoneAt          *time.Time        `firestore:"done_at,omitempty"`
	CreatedAt       time.Time         `firestore:"created_at"`
	UpdatedAt       time.Time         `firestore:"updated_at"`
	TimeoutAt       time.Time         `firestore:"timeout_at"`
	RetryCount      int               `firestore:"retry_count"`
	MaxRetries      int               `firestore:"max_retries"`
	LastRetryAt     *time.Time        `firestore:"last_retry_at,omitempty"`
	FailureReason   string            `firestore:"failure_reason,omitempty"`
	CommitTimestamp *time.Time        `firestore:"commit_timestamp,omitempty"`
}

// Participant represents a service participating in the transaction
type Participant struct {
	ServiceName string     `firestore:"service_name"`
	ServiceURL  string     `firestore:"service_url"`
	Status      string     `firestore:"status"` // "prepared", "committed", "aborted", "failed"
	DoneAt      *time.Time `firestore:"done_at,omitempty"`
	Error       string     `firestore:"error,omitempty"`
	RetryCount  int        `firestore:"retry_count"`
}

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	HotelRoomID        string `json:"hotel_room_id" binding:"required"`
	HotelRoomStartDate string `json:"hotel_room_start_date" binding:"required"`
	HotelRoomEndDate   string `json:"hotel_room_end_date" binding:"required"`
	CarID              string `json:"car_id" binding:"required"`
	CarStartDate       string `json:"car_start_date" binding:"required"`
	CarEndDate         string `json:"car_end_date" binding:"required"`
	TrainSeatID        string `json:"train_seat_id" binding:"required"`
	UserID             string `json:"user_id" binding:"required"`
}

// OrderResponse represents the response after order creation
type OrderResponse struct {
	OrderID       string            `json:"order_id"`
	TransactionID string            `json:"transaction_id"`
	Status        TransactionStatus `json:"status"`
	Message       string            `json:"message"`
}

// PrepareRequest represents the prepare phase request
type PrepareRequest struct {
	TransactionID string `json:"transaction_id"`
	OrderID       string `json:"order_id"`
	Payload       any    `json:"payload"`
}

// PrepareResponse represents the prepare phase response
type PrepareResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CommitRequest represents the commit phase request
type CommitRequest struct {
	TransactionID string `json:"transaction_id"`
	OrderID       string `json:"order_id"`
	ServiceName   string `json:"service_name"`
}

// CommitResponse represents the commit phase response
type CommitResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// AbortRequest represents the abort phase request
type AbortRequest struct {
	TransactionID string `json:"transaction_id"`
	OrderID       string `json:"order_id"`
	ServiceName   string `json:"service_name"`
	Reason        string `json:"reason"`
}

// AbortResponse represents the abort phase response
type AbortResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// TransactionStatusResponse represents the transaction status response
type TransactionStatusResponse struct {
	TransactionID string            `json:"transaction_id"`
	OrderID       string            `json:"order_id"`
	Status        TransactionStatus `json:"status"`
	Participants  []Participant     `json:"participants"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	TimeoutAt     time.Time         `json:"timeout_at"`
	RetryCount    int               `json:"retry_count"`
	FailureReason string            `json:"failure_reason,omitempty"`
}

// Config represents the coordinator configuration
type Config struct {
	TransactionTimeout time.Duration
	MaxRetries         int
	RetryDelay         time.Duration
	Services           map[string]string // service name -> service URL
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		TransactionTimeout: 30 * time.Second,
		MaxRetries:         3,
		RetryDelay:         2 * time.Second,
		Services: map[string]string{
			"hotel": "http://localhost:8081",
			"car":   "http://localhost:8082",
			"train": "http://localhost:8083",
		},
	}
}
