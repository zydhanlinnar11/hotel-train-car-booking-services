package car

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// GetCar retrieves a car by ID
func (r *Repository) GetCar(ctx context.Context, carID string) (*Car, error) {
	collection := r.client.Collection("twophase_cars")
	doc := collection.Doc(carID)

	docSnap, err := doc.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("car not found: %s", carID)
		}
		return nil, fmt.Errorf("failed to get car: %w", err)
	}

	var car Car
	if err := docSnap.DataTo(&car); err != nil {
		return nil, fmt.Errorf("failed to unmarshal car: %w", err)
	}

	return &car, nil
}

// UpdateCarAvailability updates car availability
func (r *Repository) UpdateCarAvailability(ctx context.Context, carID string, available bool) error {
	collection := r.client.Collection("twophase_cars")
	doc := collection.Doc(carID)

	_, err := doc.Update(ctx, []firestore.Update{
		{Path: "available", Value: available},
	})
	if err != nil {
		return fmt.Errorf("failed to update car availability: %w", err)
	}

	return nil
}

// CreateReservation creates a new reservation
func (r *Repository) CreateReservation(ctx context.Context, reservation *Reservation) error {
	collection := r.client.Collection("twophase_car_reservations")
	doc := collection.Doc(reservation.ID)

	_, err := doc.Create(ctx, reservation)
	if err != nil {
		return fmt.Errorf("failed to create reservation: %w", err)
	}

	return nil
}

// GetReservationByOrderID retrieves reservation by order ID
func (r *Repository) GetReservationByOrderID(ctx context.Context, orderID string) (*Reservation, error) {
	collection := r.client.Collection("twophase_car_reservations")

	query := collection.Where("order_id", "==", orderID)
	iter := query.Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, fmt.Errorf("reservation not found for order: %s", orderID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to iterate reservations: %w", err)
	}

	var reservation Reservation
	if err := doc.DataTo(&reservation); err != nil {
		return nil, fmt.Errorf("failed to unmarshal reservation: %w", err)
	}

	return &reservation, nil
}

// UpdateReservationStatus updates reservation status
func (r *Repository) UpdateReservationStatus(ctx context.Context, reservationID, status string) error {
	collection := r.client.Collection("twophase_car_reservations")
	doc := collection.Doc(reservationID)

	_, err := doc.Update(ctx, []firestore.Update{
		{Path: "status", Value: status},
		{Path: "updated_at", Value: time.Now()},
	})
	if err != nil {
		return fmt.Errorf("failed to update reservation status: %w", err)
	}

	return nil
}

// DeleteReservation deletes a reservation
func (r *Repository) DeleteReservation(ctx context.Context, reservationID string) error {
	collection := r.client.Collection("twophase_car_reservations")
	doc := collection.Doc(reservationID)

	_, err := doc.Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete reservation: %w", err)
	}

	return nil
}

// CreateTwoPhaseTransaction creates a two-phase transaction log
func (r *Repository) CreateTwoPhaseTransaction(ctx context.Context, transaction *TwoPhaseTransaction) error {
	collection := r.client.Collection("twophase_car_transactions")
	doc := collection.Doc(transaction.ID)

	_, err := doc.Create(ctx, transaction)
	if err != nil {
		return fmt.Errorf("failed to create two-phase transaction: %w", err)
	}

	return nil
}

// GetTwoPhaseTransaction retrieves a two-phase transaction
func (r *Repository) GetTwoPhaseTransaction(ctx context.Context, transactionID string) (*TwoPhaseTransaction, error) {
	collection := r.client.Collection("twophase_car_transactions")

	query := collection.Where("transaction_id", "==", transactionID)
	iter := query.Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, fmt.Errorf("two-phase transaction not found: %s", transactionID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to iterate transactions: %w", err)
	}

	var transaction TwoPhaseTransaction
	if err := doc.DataTo(&transaction); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}

	return &transaction, nil
}

// UpdateTwoPhaseTransactionStatus updates two-phase transaction status
func (r *Repository) UpdateTwoPhaseTransactionStatus(ctx context.Context, transactionID, status string) error {
	collection := r.client.Collection("twophase_car_transactions")

	query := collection.Where("transaction_id", "==", transactionID)
	iter := query.Documents(context.Background())
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return fmt.Errorf("two-phase transaction not found: %s", transactionID)
	}
	if err != nil {
		return fmt.Errorf("failed to iterate transactions: %w", err)
	}

	_, err = doc.Ref.Update(ctx, []firestore.Update{
		{Path: "status", Value: status},
		{Path: "updated_at", Value: time.Now()},
	})
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	return nil
}

// CheckCarAvailability checks if a car is available for the given dates
func (r *Repository) CheckCarAvailability(ctx context.Context, carID string, startDate, endDate time.Time) (bool, error) {
	collection := r.client.Collection("twophase_car_reservations")

	query := collection.Where("car_id", "==", carID).
		Where("status", "in", []string{"pending", "confirmed"}).
		Where("start_date", "<", endDate).
		Where("end_date", ">", startDate)

	iter := query.Documents(ctx)
	defer iter.Stop()

	// If there are any overlapping reservations, car is not available
	_, err := iter.Next()
	if err == iterator.Done {
		return true, nil // No conflicts, car is available
	}
	if err != nil {
		return false, fmt.Errorf("failed to check car availability: %w", err)
	}

	return false, nil // Car is not available
}
