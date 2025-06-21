package hotel

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/firestore/apiv1/firestorepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrHotelRoomNotFound        = errors.New("hotel room not found")
	ErrHotelReservationNotFound = errors.New("hotel reservation not found")
)

type Repository interface {
	GetHotelRoomByID(ctx context.Context, id string) (*HotelRoom, error)
	CreateHotelReservation(ctx context.Context, hotelReservation *HotelReservation) error
	GetHotelReservationByID(ctx context.Context, id string) (*HotelReservation, error)
	UpdateHotelReservation(ctx context.Context, hotelReservation *HotelReservation) error
	IsHotelRoomAvailable(ctx context.Context, hotelRoomID string, startDate, endDate string) (bool, error)
}

const (
	hotelRoomCollection        = "hotel_rooms"
	hotelReservationCollection = "hotel_reservations"
)

type firestoreRepository struct {
	client *firestore.Client
}

func NewFirestoreRepository(client *firestore.Client) Repository {
	return &firestoreRepository{client: client}
}

func (r *firestoreRepository) GetHotelRoomByID(ctx context.Context, id string) (*HotelRoom, error) {
	doc, err := r.client.Collection(hotelRoomCollection).Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil, ErrHotelRoomNotFound
	}
	if err != nil {
		return nil, err
	}

	var hotelRoom HotelRoom
	if err := doc.DataTo(&hotelRoom); err != nil {
		return nil, err
	}

	return &hotelRoom, nil
}

func (r *firestoreRepository) CreateHotelReservation(ctx context.Context, hotelReservation *HotelReservation) error {
	_, err := r.client.Collection(hotelReservationCollection).Doc(hotelReservation.ID).Set(ctx, hotelReservation)
	return err
}

func (r *firestoreRepository) GetHotelReservationByID(ctx context.Context, id string) (*HotelReservation, error) {
	doc, err := r.client.Collection(hotelReservationCollection).Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil, ErrHotelReservationNotFound
	}
	if err != nil {
		return nil, err
	}

	var hotelReservation HotelReservation
	if err := doc.DataTo(&hotelReservation); err != nil {
		return nil, err
	}

	return &hotelReservation, nil
}

func (r *firestoreRepository) UpdateHotelReservation(ctx context.Context, hotelReservation *HotelReservation) error {
	_, err := r.client.Collection(hotelReservationCollection).Doc(hotelReservation.ID).Set(ctx, hotelReservation)
	return err
}

func (r *firestoreRepository) IsHotelRoomAvailable(ctx context.Context, hotelRoomID string, startDate, endDate string) (bool, error) {
	query := r.client.Collection(hotelReservationCollection).
		Where("hotel_room_id", "==", hotelRoomID).
		Where("hotel_room_start_date", "<=", endDate).
		Where("hotel_room_end_date", ">=", startDate).
		Where("status", "!=", HotelRoomReservationStatusCancelled)

	aggregationQuery := query.NewAggregationQuery().WithCount("all")
	results, err := aggregationQuery.Get(ctx)
	if err != nil {
		return false, err
	}

	count, ok := results["all"]
	if !ok {
		return false, errors.New("firestore: couldn't get alias for COUNT from results")
	}

	countValue := count.(*firestorepb.Value)

	return countValue.GetIntegerValue() == 0, nil
}
