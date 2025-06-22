package train

import (
	"time"
)

type TrainSeatReservationStatus string

const (
	TrainSeatReservationStatusCancelled TrainSeatReservationStatus = "CANCELLED"
	TrainSeatReservationStatusReserved  TrainSeatReservationStatus = "RESERVED"
)

type TwoPhaseTransactionStatus string

const (
	TwoPhaseTransactionStatusPrepared  TwoPhaseTransactionStatus = "PREPARED"
	TwoPhaseTransactionStatusCommitted TwoPhaseTransactionStatus = "COMMITTED"
	TwoPhaseTransactionStatusAborted   TwoPhaseTransactionStatus = "ABORTED"
)

type TrainSeatTicket struct {
	SeatID    string `firestore:"seat_id" json:"seat_id"`
	TrainName string `firestore:"train_name" json:"train_name"`
	Available bool   `firestore:"available" json:"available"`
}

type TrainSeatReservation struct {
	ID            string                     `firestore:"id" json:"id"`
	SeatID        string                     `firestore:"seat_id" json:"seat_id"`
	TrainName     string                     `firestore:"train_name" json:"train_name"`
	TransactionID string                     `firestore:"transaction_id" json:"transaction_id"`
	Status        TrainSeatReservationStatus `firestore:"status" json:"status"`
}

// TwoPhaseTransaction represents a two-phase commit transaction for train seat reservation
type TwoPhaseTransaction struct {
	Id            string                    `firestore:"id"`
	Status        TwoPhaseTransactionStatus `firestore:"status"` // "prepared", "committed", "aborted"
	ReservationID string                    `firestore:"reservation_id,omitempty"`
	CreatedAt     time.Time                 `firestore:"created_at"`
	UpdatedAt     time.Time                 `firestore:"updated_at"`
}

type TrainSeatReservationPayload struct {
	SeatID    string `json:"seat_id" binding:"required"`
	TrainName string `json:"train_name" binding:"required"`
}
