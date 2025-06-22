package hotel

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
	HotelRoomAvailabilityCollection = "twophase_hotel_room_availabilities"
	HotelRoomReservationCollection  = "twophase_hotel_room_reservations"
	HotelRoomTransactionCollection  = "twophase_hotel_transactions"
)

var (
	ErrRoomNotAvailable = errors.New("room not available")
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
	ref := r.client.Collection(HotelRoomTransactionCollection).Doc(transactionID)

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

func (r *Repository) getRoomAvailabilityId(roomID, date string) string {
	return fmt.Sprintf("%s-%s", roomID, date)
}

func (r *Repository) getRoomAvailabilityRefs(roomID string, checkInDate, checkOutDate string) ([]*firestore.DocumentRef, error) {
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

	var roomAvailabilityRefs []*firestore.DocumentRef
	for _, date := range dates {
		roomAvailabilityRefs = append(
			roomAvailabilityRefs,
			r.client.Collection(HotelRoomAvailabilityCollection).
				Doc(r.getRoomAvailabilityId(roomID, date)),
		)
	}

	return roomAvailabilityRefs, nil
}

func (r *Repository) CommitRoomReservation(ctx context.Context, transactionID string) error {
	transactionRef := r.client.Collection(HotelRoomTransactionCollection).Doc(transactionID)

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

func (r *Repository) AbortRoomReservation(ctx context.Context, transactionID string) error {
	transactionRef := r.client.Collection(HotelRoomTransactionCollection).Doc(transactionID)

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

		reservationRef := r.client.Collection(HotelRoomReservationCollection).Doc(transaction.ReservationID)
		reservationDoc, err := tx.Get(reservationRef)
		if err != nil {
			return fmt.Errorf("failed to get reservation: %w", err)
		}

		var reservation HotelReservation
		if err := reservationDoc.DataTo(&reservation); err != nil {
			return fmt.Errorf("failed to unmarshal reservation: %w", err)
		}

		roomAvailabilityRefs, err := r.getRoomAvailabilityRefs(reservation.HotelRoomID, reservation.HotelRoomStartDate, reservation.HotelRoomEndDate)
		if err != nil {
			return fmt.Errorf("failed to get room availability refs: %w", err)
		}

		for _, ref := range roomAvailabilityRefs {
			doc, err := tx.Get(ref)
			if err != nil {
				return fmt.Errorf("failed to get room availability: %w", err)
			}

			var roomAvailability HotelRoomAvailability
			if err := doc.DataTo(&roomAvailability); err != nil {
				return fmt.Errorf("failed to unmarshal room availability: %w", err)
			}

			if err := tx.Update(ref, []firestore.Update{
				{Path: "available", Value: true},
			}); err != nil {
				return fmt.Errorf("failed to update room availability: %w", err)
			}
		}

		if err := tx.Update(transactionRef, []firestore.Update{
			{Path: "status", Value: TwoPhaseTransactionStatusAborted},
			{Path: "updated_at", Value: time.Now()},
		}); err != nil {
			return fmt.Errorf("failed to update transaction: %w", err)
		}

		if err := tx.Update(reservationRef, []firestore.Update{
			{Path: "status", Value: HotelRoomReservationStatusCancelled},
			{Path: "updated_at", Value: time.Now()},
		}); err != nil {
			return fmt.Errorf("failed to update reservation: %w", err)
		}

		return nil
	})
}

// PrepareRoomReservation prepares a room reservation
func (r *Repository) PrepareRoomReservation(ctx context.Context, transactionID, roomID string, checkInDate, checkOutDate string) error {
	roomAvailabilityRefs, err := r.getRoomAvailabilityRefs(roomID, checkInDate, checkOutDate)
	if err != nil {
		return fmt.Errorf("failed to get room availability refs: %w", err)
	}

	if len(roomAvailabilityRefs) == 0 {
		return ErrRoomNotAvailable
	}

	return r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		var roomAvailability HotelRoomAvailability
		for _, ref := range roomAvailabilityRefs {
			doc, err := tx.Get(ref)
			if status.Code(err) == codes.NotFound {
				return ErrRoomNotAvailable
			}
			if err != nil {
				return fmt.Errorf("failed to get room availability: %w", err)
			}

			if err := doc.DataTo(&roomAvailability); err != nil {
				return fmt.Errorf("failed to unmarshal room availability: %w", err)
			}

			if !roomAvailability.Available {
				return ErrRoomNotAvailable
			}

			roomAvailability.Available = false

			if err := tx.Update(ref, []firestore.Update{
				{Path: "available", Value: false},
			}); err != nil {
				return fmt.Errorf("failed to update room availability: %w", err)
			}
		}

		hotelRoomReservation := &HotelReservation{
			ID:                 ulid.Make().String(),
			TransactionID:      transactionID,
			HotelRoomID:        roomID,
			HotelRoomName:      roomAvailability.RoomName,
			HotelName:          roomAvailability.HotelName,
			HotelRoomStartDate: checkInDate,
			HotelRoomEndDate:   checkOutDate,
			Status:             HotelRoomReservationStatusReserved,
		}

		hotelRoomReservationRef := r.client.Collection(HotelRoomReservationCollection).Doc(hotelRoomReservation.ID)
		if err := tx.Create(hotelRoomReservationRef, hotelRoomReservation); err != nil {
			return fmt.Errorf("failed to create hotel room reservation: %w", err)
		}

		twoPhaseTransaction := &TwoPhaseTransaction{
			Id:            transactionID,
			Status:        TwoPhaseTransactionStatusPrepared,
			ReservationID: hotelRoomReservation.ID,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		twoPhaseTransactionRef := r.client.Collection(HotelRoomTransactionCollection).Doc(twoPhaseTransaction.Id)
		if err := tx.Create(twoPhaseTransactionRef, twoPhaseTransaction); err != nil {
			return fmt.Errorf("failed to create two-phase transaction: %w", err)
		}

		return nil
	})
}

func (r *Repository) BulkWriteHotelRoomAvailability(ctx context.Context, hotelRoomAvailabilities []HotelRoomAvailability) error {
	collection := r.client.Collection(HotelRoomAvailabilityCollection)
	bw := r.client.BulkWriter(ctx)

	for _, hotelRoomAvailability := range hotelRoomAvailabilities {
		docRef := collection.Doc(r.getRoomAvailabilityId(hotelRoomAvailability.RoomID, hotelRoomAvailability.Date))
		bw.Set(docRef, hotelRoomAvailability)
	}

	// Flush all writes
	bw.Flush()

	return nil
}
