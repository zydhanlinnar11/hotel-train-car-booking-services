package car

import (
	"time"
)

type CarReservationStatus string

const (
	CarReservationStatusCancelled CarReservationStatus = "CANCELLED"
	CarReservationStatusReserved  CarReservationStatus = "RESERVED"
)

type TwoPhaseTransactionStatus string

const (
	TwoPhaseTransactionStatusPrepared  TwoPhaseTransactionStatus = "PREPARED"
	TwoPhaseTransactionStatusCommitted TwoPhaseTransactionStatus = "COMMITTED"
	TwoPhaseTransactionStatusAborted   TwoPhaseTransactionStatus = "ABORTED"
)

type CarAvailability struct {
	CarID     string `firestore:"car_id" json:"car_id"`
	CarName   string `firestore:"car_name" json:"car_name"`
	Date      string `firestore:"date" json:"date"`
	Available bool   `firestore:"available" json:"available"`
}

type CarReservation struct {
	ID            string               `firestore:"id" json:"id"`
	CarID         string               `firestore:"car_id" json:"car_id"`
	CarName       string               `firestore:"car_name" json:"car_name"`
	CarStartDate  string               `firestore:"car_start_date" json:"car_start_date"`
	CarEndDate    string               `firestore:"car_end_date" json:"car_end_date"`
	TransactionID string               `firestore:"transaction_id" json:"transaction_id"`
	Status        CarReservationStatus `firestore:"status" json:"status"`
}

// TwoPhaseTransaction represents a two-phase commit transaction for hotel
type TwoPhaseTransaction struct {
	Id            string                    `firestore:"id"`
	Status        TwoPhaseTransactionStatus `firestore:"status"` // "prepared", "committed", "aborted"
	ReservationID string                    `firestore:"reservation_id,omitempty"`
	CreatedAt     time.Time                 `firestore:"created_at"`
	UpdatedAt     time.Time                 `firestore:"updated_at"`
}

type CarReservationPayload struct {
	CarID        string `json:"car_id" binding:"required"`
	CarStartDate string `json:"car_start_date" binding:"required"`
	CarEndDate   string `json:"car_end_date" binding:"required"`
}
