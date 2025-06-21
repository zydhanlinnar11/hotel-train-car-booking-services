package car

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/firestore/apiv1/firestorepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrCarNotFound            = errors.New("car not found")
	ErrCarReservationNotFound = errors.New("car reservation not found")
)

type Repository interface {
	GetCarByID(ctx context.Context, id string) (*Car, error)
	CreateCarReservation(ctx context.Context, carReservation *CarReservation) error
	GetCarReservationByID(ctx context.Context, id string) (*CarReservation, error)
	UpdateCarReservation(ctx context.Context, carReservation *CarReservation) error
	IsCarAvailable(ctx context.Context, carID string, startDate, endDate string) (bool, error)
}

const (
	carCollection            = "cars"
	carReservationCollection = "car_reservations"
)

type firestoreRepository struct {
	client *firestore.Client
}

func NewFirestoreRepository(client *firestore.Client) Repository {
	return &firestoreRepository{client: client}
}

func (r *firestoreRepository) GetCarByID(ctx context.Context, id string) (*Car, error) {
	doc, err := r.client.Collection(carCollection).Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil, ErrCarNotFound
	}
	if err != nil {
		return nil, err
	}

	var car Car
	if err := doc.DataTo(&car); err != nil {
		return nil, err
	}

	return &car, nil
}

func (r *firestoreRepository) CreateCarReservation(ctx context.Context, carReservation *CarReservation) error {
	_, err := r.client.Collection(carReservationCollection).Doc(carReservation.ID).Set(ctx, carReservation)
	return err
}

func (r *firestoreRepository) GetCarReservationByID(ctx context.Context, id string) (*CarReservation, error) {
	doc, err := r.client.Collection(carReservationCollection).Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil, ErrCarReservationNotFound
	}
	if err != nil {
		return nil, err
	}

	var carReservation CarReservation
	if err := doc.DataTo(&carReservation); err != nil {
		return nil, err
	}

	return &carReservation, nil
}

func (r *firestoreRepository) UpdateCarReservation(ctx context.Context, carReservation *CarReservation) error {
	_, err := r.client.Collection(carReservationCollection).Doc(carReservation.ID).Set(ctx, carReservation)
	return err
}

func (r *firestoreRepository) IsCarAvailable(ctx context.Context, carID string, startDate, endDate string) (bool, error) {
	query := r.client.Collection(carReservationCollection).
		Where("car_id", "==", carID).
		Where("start_date", "<=", endDate).
		Where("end_date", ">=", startDate)

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
