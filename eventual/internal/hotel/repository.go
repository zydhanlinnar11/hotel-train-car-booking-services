package hotel

import (
	"context"

	"cloud.google.com/go/firestore"
)

type Repository interface {
	GetHotelRoomByID(ctx context.Context, id string) (*HotelRoom, error)
	CreateHotelReservation(ctx context.Context, hotelReservation *HotelReservation) error
	GetHotelReservationByID(ctx context.Context, id string) (*HotelReservation, error)
	UpdateHotelReservation(ctx context.Context, hotelReservation *HotelReservation) error
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
