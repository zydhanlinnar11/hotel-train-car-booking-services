package train

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/oklog/ulid/v2"
)

const (
	TrainSeatTicketCollection      = "twophase_train_seat_tickets"
	TrainSeatReservationCollection = "twophase_train_seat_reservations"
	TrainTransactionCollection     = "twophase_train_transactions"
)

var (
	ErrSeatNotAvailable = errors.New("seat not available")
)

// Repository handles Firestore operations for hotel service
type Repository struct {
	client *firestore.Client
}

// NewRepository creates a new repository instance
func NewRepository(client *firestore.Client) *Repository {
	return &Repository{
		client: client,
	}
}

// GetTwoPhaseTransaction retrieves a two-phase transaction
func (r *Repository) GetTwoPhaseTransaction(ctx context.Context, transactionID string) (*TwoPhaseTransaction, error) {
	ref := r.client.Collection(TrainTransactionCollection).Doc(transactionID)

	doc, err := ref.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	var transaction TwoPhaseTransaction
	if err := doc.DataTo(&transaction); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}

	return &transaction, nil
}

func (r *Repository) CommitSeatReservation(ctx context.Context, transactionID string) error {
	transactionRef := r.client.Collection(TrainTransactionCollection).Doc(transactionID)

	return r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		transactionDoc, err := tx.Get(transactionRef)
		if err != nil {
			return fmt.Errorf("failed to get transaction: %w", err)
		}

		var transaction TwoPhaseTransaction
		if err := transactionDoc.DataTo(&transaction); err != nil {
			return fmt.Errorf("failed to unmarshal transaction: %w", err)
		}

		if transaction.Status != TwoPhaseTransactionStatusPrepared {
			// Already committed or aborted
			return nil
		}

		if err := tx.Update(transactionRef, []firestore.Update{
			{Path: "status", Value: TwoPhaseTransactionStatusCommitted},
			{Path: "updated_at", Value: time.Now()},
		}); err != nil {
			return fmt.Errorf("failed to update transaction: %w", err)
		}

		return nil
	})
}

func (r *Repository) AbortSeatReservation(ctx context.Context, transactionID string) error {
	transactionRef := r.client.Collection(TrainTransactionCollection).Doc(transactionID)

	return r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		transactionDoc, err := tx.Get(transactionRef)
		if err != nil {
			return fmt.Errorf("failed to get transaction: %w", err)
		}

		var transaction TwoPhaseTransaction
		if err := transactionDoc.DataTo(&transaction); err != nil {
			return fmt.Errorf("failed to unmarshal transaction: %w", err)
		}

		if transaction.Status != TwoPhaseTransactionStatusPrepared {
			// Already committed or aborted
			return nil
		}

		reservationRef := r.client.Collection(TrainSeatReservationCollection).Doc(transaction.ReservationID)
		reservationDoc, err := tx.Get(reservationRef)
		if err != nil {
			return fmt.Errorf("failed to get reservation: %w", err)
		}

		var reservation TrainSeatReservation
		if err := reservationDoc.DataTo(&reservation); err != nil {
			return fmt.Errorf("failed to unmarshal reservation: %w", err)
		}

		ticketRef := r.client.Collection(TrainSeatTicketCollection).Doc(reservation.SeatID)
		ticketDoc, err := tx.Get(ticketRef)
		if err != nil {
			return fmt.Errorf("failed to get ticket: %w", err)
		}

		var ticket TrainSeatTicket
		if err := ticketDoc.DataTo(&ticket); err != nil {
			return fmt.Errorf("failed to unmarshal ticket: %w", err)
		}

		if !ticket.Available {
			return ErrSeatNotAvailable
		}

		if err := tx.Update(ticketRef, []firestore.Update{
			{Path: "available", Value: true},
		}); err != nil {
			return fmt.Errorf("failed to update ticket: %w", err)
		}

		if err := tx.Update(transactionRef, []firestore.Update{
			{Path: "status", Value: TwoPhaseTransactionStatusAborted},
			{Path: "updated_at", Value: time.Now()},
		}); err != nil {
			return fmt.Errorf("failed to update transaction: %w", err)
		}

		if err := tx.Update(reservationRef, []firestore.Update{
			{Path: "status", Value: TrainSeatReservationStatusCancelled},
			{Path: "updated_at", Value: time.Now()},
		}); err != nil {
			return fmt.Errorf("failed to update reservation: %w", err)
		}

		return nil
	})
}

// PrepareSeatReservation prepares a seat reservation
func (r *Repository) PrepareSeatReservation(ctx context.Context, transactionID, seatID string) error {
	ticketRef := r.client.Collection(TrainSeatTicketCollection).Doc(seatID)

	return r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		var ticket TrainSeatTicket
		ticketDoc, err := tx.Get(ticketRef)
		if err != nil {
			return fmt.Errorf("failed to get ticket: %w", err)
		}

		if err := ticketDoc.DataTo(&ticket); err != nil {
			return fmt.Errorf("failed to unmarshal ticket: %w", err)
		}

		if !ticket.Available {
			return ErrSeatNotAvailable
		}

		if err := tx.Update(ticketRef, []firestore.Update{
			{Path: "available", Value: false},
		}); err != nil {
			return fmt.Errorf("failed to update ticket: %w", err)
		}

		trainSeatReservation := &TrainSeatReservation{
			ID:            ulid.Make().String(),
			SeatID:        seatID,
			TrainName:     ticket.TrainName,
			TransactionID: transactionID,
			Status:        TrainSeatReservationStatusReserved,
		}

		trainSeatReservationRef := r.client.Collection(TrainSeatReservationCollection).Doc(trainSeatReservation.ID)
		if err := tx.Create(trainSeatReservationRef, trainSeatReservation); err != nil {
			return fmt.Errorf("failed to create train seat reservation: %w", err)
		}

		twoPhaseTransaction := &TwoPhaseTransaction{
			Id:            transactionID,
			Status:        TwoPhaseTransactionStatusPrepared,
			ReservationID: trainSeatReservation.ID,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		twoPhaseTransactionRef := r.client.Collection(TrainTransactionCollection).Doc(twoPhaseTransaction.Id)
		if err := tx.Create(twoPhaseTransactionRef, twoPhaseTransaction); err != nil {
			return fmt.Errorf("failed to create two-phase transaction: %w", err)
		}

		return nil
	})
}

func (r *Repository) BulkWriteTrainSeatTicket(ctx context.Context, trainSeatTickets []TrainSeatTicket) error {
	collection := r.client.Collection(TrainSeatTicketCollection)
	bw := r.client.BulkWriter(ctx)

	for _, trainSeatTicket := range trainSeatTickets {
		docRef := collection.Doc(trainSeatTicket.SeatID)
		bw.Set(docRef, trainSeatTicket)
	}

	// Flush all writes
	bw.Flush()

	return nil
}
