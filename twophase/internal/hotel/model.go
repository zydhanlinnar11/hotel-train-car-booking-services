package hotel

import (
	"time"
)

type HotelRoomReservationStatus string

const (
	HotelRoomReservationStatusCancelled HotelRoomReservationStatus = "CANCELLED"
	HotelRoomReservationStatusReserved  HotelRoomReservationStatus = "RESERVED"
)

type TwoPhaseTransactionStatus string

const (
	TwoPhaseTransactionStatusPrepared  TwoPhaseTransactionStatus = "PREPARED"
	TwoPhaseTransactionStatusCommitted TwoPhaseTransactionStatus = "COMMITTED"
	TwoPhaseTransactionStatusAborted   TwoPhaseTransactionStatus = "ABORTED"
)

type HotelRoomAvailability struct {
	RoomID    string `firestore:"room_id" json:"room_id"`
	HotelName string `firestore:"hotel_name" json:"hotel_name"`
	RoomName  string `firestore:"room_name" json:"room_name"`
	Date      string `firestore:"date" json:"date"`
	Available bool   `firestore:"available" json:"available"`
}

type HotelReservation struct {
	ID                 string                     `firestore:"id" json:"id"`
	HotelRoomID        string                     `firestore:"hotel_room_id" json:"hotel_room_id"`
	HotelRoomName      string                     `firestore:"hotel_room_name" json:"hotel_room_name"`
	HotelName          string                     `firestore:"hotel_name" json:"hotel_name"`
	HotelRoomStartDate string                     `firestore:"hotel_room_start_date" json:"hotel_room_start_date"`
	HotelRoomEndDate   string                     `firestore:"hotel_room_end_date" json:"hotel_room_end_date"`
	TransactionID      string                     `firestore:"transaction_id" json:"transaction_id"`
	Status             HotelRoomReservationStatus `firestore:"status" json:"status"`
}

// TwoPhaseTransaction represents a two-phase commit transaction for hotel
type TwoPhaseTransaction struct {
	Id            string                    `firestore:"id"`
	Status        TwoPhaseTransactionStatus `firestore:"status"` // "prepared", "committed", "aborted"
	ReservationID string                    `firestore:"reservation_id,omitempty"`
	CreatedAt     time.Time                 `firestore:"created_at"`
	UpdatedAt     time.Time                 `firestore:"updated_at"`
}

// PrepareRequest represents prepare phase request
type PrepareRequest struct {
	TransactionID string `json:"transaction_id"`
	OrderID       string `json:"order_id"`
	Payload       struct {
		HotelRoomID        string `json:"hotel_room_id" binding:"required"`
		HotelRoomStartDate string `json:"hotel_room_start_date" binding:"required"`
		HotelRoomEndDate   string `json:"hotel_room_end_date" binding:"required"`
	} `json:"payload"`
}

// PrepareResponse represents prepare phase response
type PrepareResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CommitRequest represents commit phase request
type CommitRequest struct {
	TransactionID string `json:"transaction_id" binding:"required"`
}

// CommitResponse represents commit phase response
type CommitResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// AbortRequest represents abort phase request
type AbortRequest struct {
	TransactionID string `json:"transaction_id" binding:"required"`
}

// AbortResponse represents abort phase response
type AbortResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
