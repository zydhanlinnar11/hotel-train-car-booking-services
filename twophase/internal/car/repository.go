package car

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/oklog/ulid/v2"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/twophase/pkg/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	CarAvailabilityCollection = "twophase_car_availabilities"
	CarReservationCollection  = "twophase_car_reservations"
	CarTransactionCollection  = "twophase_car_transactions"
)

var (
	ErrCarNotAvailable = errors.New("car not available")
)

// Repository handles Firestore operations for car service
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
	ref := r.client.Collection(CarTransactionCollection).Doc(transactionID)

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

func (r *Repository) getCarAvailabilityId(carID, date string) string {
	return fmt.Sprintf("%s-%s", carID, date)
}

func (r *Repository) getCarAvailabilityRefs(carID string, checkInDate, checkOutDate string) ([]*firestore.DocumentRef, error) {
	startDate, err := time.Parse(config.DateFormat, checkInDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse check-in date: %w", err)
	}

	endDate, err := time.Parse(config.DateFormat, checkOutDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse check-out date: %w", err)
	}

	// Iterate through the dates (inclusive)
	var dates []string
	for date := startDate; date.Before(endDate.AddDate(0, 0, 1)); date = date.AddDate(0, 0, 1) {
		dates = append(dates, date.Format(config.DateFormat))
	}

	var carAvailabilityRefs []*firestore.DocumentRef
	for _, date := range dates {
		carAvailabilityRefs = append(
			carAvailabilityRefs,
			r.client.Collection(CarAvailabilityCollection).
				Doc(r.getCarAvailabilityId(carID, date)),
		)
	}

	return carAvailabilityRefs, nil
}

func (r *Repository) CommitCarReservation(ctx context.Context, transactionID string) error {
	transactionRef := r.client.Collection(CarTransactionCollection).Doc(transactionID)

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

func (r *Repository) AbortCarReservation(ctx context.Context, transactionID string) error {
	transactionRef := r.client.Collection(CarTransactionCollection).Doc(transactionID)

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

		reservationRef := r.client.Collection(CarReservationCollection).Doc(transaction.ReservationID)
		reservationDoc, err := tx.Get(reservationRef)
		if err != nil {
			return fmt.Errorf("failed to get reservation: %w", err)
		}

		var reservation CarReservation
		if err := reservationDoc.DataTo(&reservation); err != nil {
			return fmt.Errorf("failed to unmarshal reservation: %w", err)
		}

		carAvailabilityRefs, err := r.getCarAvailabilityRefs(reservation.CarID, reservation.CarStartDate, reservation.CarEndDate)
		if err != nil {
			return fmt.Errorf("failed to get car availability refs: %w", err)
		}

		for _, ref := range carAvailabilityRefs {
			doc, err := tx.Get(ref)
			if err != nil {
				return fmt.Errorf("failed to get car availability: %w", err)
			}

			var carAvailability CarAvailability
			if err := doc.DataTo(&carAvailability); err != nil {
				return fmt.Errorf("failed to unmarshal car availability: %w", err)
			}

			if err := tx.Update(ref, []firestore.Update{
				{Path: "available", Value: true},
			}); err != nil {
				return fmt.Errorf("failed to update car availability: %w", err)
			}
		}

		if err := tx.Update(transactionRef, []firestore.Update{
			{Path: "status", Value: TwoPhaseTransactionStatusAborted},
			{Path: "updated_at", Value: time.Now()},
		}); err != nil {
			return fmt.Errorf("failed to update transaction: %w", err)
		}

		if err := tx.Update(reservationRef, []firestore.Update{
			{Path: "status", Value: CarReservationStatusCancelled},
			{Path: "updated_at", Value: time.Now()},
		}); err != nil {
			return fmt.Errorf("failed to update reservation: %w", err)
		}

		return nil
	})
}

// PrepareCarReservation prepares a car reservation
func (r *Repository) PrepareCarReservation(ctx context.Context, transactionID, carID string, checkInDate, checkOutDate string) error {
	carAvailabilityRefs, err := r.getCarAvailabilityRefs(carID, checkInDate, checkOutDate)
	if err != nil {
		return fmt.Errorf("failed to get car availability refs: %w", err)
	}

	if len(carAvailabilityRefs) == 0 {
		return ErrCarNotAvailable
	}

	return r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		var carAvailability CarAvailability
		for _, ref := range carAvailabilityRefs {
			doc, err := tx.Get(ref)
			if status.Code(err) == codes.NotFound {
				return ErrCarNotAvailable
			}
			if err != nil {
				return fmt.Errorf("failed to get car availability: %w", err)
			}

			if err := doc.DataTo(&carAvailability); err != nil {
				return fmt.Errorf("failed to unmarshal car availability: %w", err)
			}

			if !carAvailability.Available {
				return ErrCarNotAvailable
			}

			carAvailability.Available = false

			if err := tx.Update(ref, []firestore.Update{
				{Path: "available", Value: false},
			}); err != nil {
				return fmt.Errorf("failed to update car availability: %w", err)
			}
		}

		carReservation := &CarReservation{
			ID:            ulid.Make().String(),
			TransactionID: transactionID,
			CarID:         carID,
			CarName:       carAvailability.CarName,
			CarStartDate:  checkInDate,
			CarEndDate:    checkOutDate,
			Status:        CarReservationStatusReserved,
		}

		carReservationRef := r.client.Collection(CarReservationCollection).Doc(carReservation.ID)
		if err := tx.Create(carReservationRef, carReservation); err != nil {
			return fmt.Errorf("failed to create car reservation: %w", err)
		}

		twoPhaseTransaction := &TwoPhaseTransaction{
			Id:            transactionID,
			Status:        TwoPhaseTransactionStatusPrepared,
			ReservationID: carReservation.ID,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		twoPhaseTransactionRef := r.client.Collection(CarTransactionCollection).Doc(twoPhaseTransaction.Id)
		if err := tx.Create(twoPhaseTransactionRef, twoPhaseTransaction); err != nil {
			return fmt.Errorf("failed to create two-phase transaction: %w", err)
		}

		return nil
	})
}

func (r *Repository) BulkWriteCarAvailability(ctx context.Context, carAvailabilities []CarAvailability) error {
	collection := r.client.Collection(CarAvailabilityCollection)
	bw := r.client.BulkWriter(ctx)

	for _, carAvailability := range carAvailabilities {
		docRef := collection.Doc(r.getCarAvailabilityId(carAvailability.CarID, carAvailability.Date))
		bw.Set(docRef, carAvailability)
	}

	// Flush all writes
	bw.Flush()

	return nil
}
