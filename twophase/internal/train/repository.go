package train

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Repository handles Firestore operations for train service
type Repository struct {
	client *firestore.Client
}

// NewRepository creates a new repository instance
func NewRepository(client *firestore.Client) *Repository {
	return &Repository{
		client: client,
	}
}

// GetSeat retrieves a seat by ID
func (r *Repository) GetSeat(ctx context.Context, seatID string) (*Seat, error) {
	collection := r.client.Collection("twophase_train_seats")
	doc := collection.Doc(seatID)

	docSnap, err := doc.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("seat not found: %s", seatID)
		}
		return nil, fmt.Errorf("failed to get seat: %w", err)
	}

	var seat Seat
	if err := docSnap.DataTo(&seat); err != nil {
		return nil, fmt.Errorf("failed to unmarshal seat: %w", err)
	}

	return &seat, nil
}

// UpdateSeatAvailability updates seat availability
func (r *Repository) UpdateSeatAvailability(ctx context.Context, seatID string, available bool) error {
	collection := r.client.Collection("twophase_train_seats")
	doc := collection.Doc(seatID)

	_, err := doc.Update(ctx, []firestore.Update{
		{Path: "available", Value: available},
	})
	if err != nil {
		return fmt.Errorf("failed to update seat availability: %w", err)
	}

	return nil
}

// CreateReservation creates a new reservation
func (r *Repository) CreateReservation(ctx context.Context, reservation *Reservation) error {
	collection := r.client.Collection("twophase_train_reservations")
	doc := collection.Doc(reservation.ID)

	_, err := doc.Create(ctx, reservation)
	if err != nil {
		return fmt.Errorf("failed to create reservation: %w", err)
	}

	return nil
}

// GetReservationByOrderID retrieves reservation by order ID
func (r *Repository) GetReservationByOrderID(ctx context.Context, orderID string) (*Reservation, error) {
	collection := r.client.Collection("twophase_train_reservations")

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
	collection := r.client.Collection("twophase_train_reservations")
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
	collection := r.client.Collection("twophase_train_reservations")
	doc := collection.Doc(reservationID)

	_, err := doc.Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete reservation: %w", err)
	}

	return nil
}

// CreateTwoPhaseTransaction creates a two-phase transaction log
func (r *Repository) CreateTwoPhaseTransaction(ctx context.Context, transaction *TwoPhaseTransaction) error {
	collection := r.client.Collection("twophase_train_transactions")
	doc := collection.Doc(transaction.ID)

	_, err := doc.Create(ctx, transaction)
	if err != nil {
		return fmt.Errorf("failed to create two-phase transaction: %w", err)
	}

	return nil
}

// GetTwoPhaseTransaction retrieves a two-phase transaction
func (r *Repository) GetTwoPhaseTransaction(ctx context.Context, transactionID string) (*TwoPhaseTransaction, error) {
	collection := r.client.Collection("twophase_train_transactions")

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
	collection := r.client.Collection("twophase_train_transactions")

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

// CheckSeatAvailability checks if a seat is available for the given date
func (r *Repository) CheckSeatAvailability(ctx context.Context, seatID string, travelDate time.Time) (bool, error) {
	collection := r.client.Collection("twophase_train_reservations")

	query := collection.Where("seat_id", "==", seatID).
		Where("status", "in", []string{"pending", "confirmed"}).
		Where("travel_date", "==", travelDate)

	iter := query.Documents(ctx)
	defer iter.Stop()

	// If there are any reservations for this seat on this date, seat is not available
	_, err := iter.Next()
	if err == iterator.Done {
		return true, nil // No conflicts, seat is available
	}
	if err != nil {
		return false, fmt.Errorf("failed to check seat availability: %w", err)
	}

	return false, nil // Seat is not available
}
