package train

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/firestore/apiv1/firestorepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrTrainSeatNotFound        = errors.New("train seat not found")
	ErrTrainReservationNotFound = errors.New("train reservation not found")
)

type Repository interface {
	GetTrainSeatByID(ctx context.Context, id string) (*TrainSeat, error)
	CreateTrainReservation(ctx context.Context, trainReservation *TrainReservation) error
	GetTrainReservationByID(ctx context.Context, id string) (*TrainReservation, error)
	GetTrainReservationByOrderID(ctx context.Context, orderID string) (*TrainReservation, error)
	UpdateTrainReservation(ctx context.Context, trainReservation *TrainReservation) error
	IsTrainSeatAvailable(ctx context.Context, seatID string) (bool, error)
}

const (
	trainSeatCollection        = "train_seats"
	trainReservationCollection = "train_reservations"
)

type firestoreRepository struct {
	client *firestore.Client
}

func NewFirestoreRepository(client *firestore.Client) Repository {
	return &firestoreRepository{client: client}
}

func (r *firestoreRepository) GetTrainSeatByID(ctx context.Context, id string) (*TrainSeat, error) {
	doc, err := r.client.Collection(trainSeatCollection).Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil, ErrTrainSeatNotFound
	}
	if err != nil {
		return nil, err
	}

	var trainSeat TrainSeat
	if err := doc.DataTo(&trainSeat); err != nil {
		return nil, err
	}

	return &trainSeat, nil
}

func (r *firestoreRepository) CreateTrainReservation(ctx context.Context, trainReservation *TrainReservation) error {
	_, err := r.client.Collection(trainReservationCollection).Doc(trainReservation.ID).Set(ctx, trainReservation)
	return err
}

func (r *firestoreRepository) GetTrainReservationByID(ctx context.Context, id string) (*TrainReservation, error) {
	doc, err := r.client.Collection(trainReservationCollection).Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return nil, ErrTrainReservationNotFound
	}
	if err != nil {
		return nil, err
	}

	var trainReservation TrainReservation
	if err := doc.DataTo(&trainReservation); err != nil {
		return nil, err
	}

	return &trainReservation, nil
}

func (r *firestoreRepository) GetTrainReservationByOrderID(ctx context.Context, orderID string) (*TrainReservation, error) {
	query := r.client.Collection(trainReservationCollection).Where("order_id", "==", orderID)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, ErrTrainReservationNotFound
	}

	var trainReservation TrainReservation
	if err := docs[0].DataTo(&trainReservation); err != nil {
		return nil, err
	}

	return &trainReservation, nil
}

func (r *firestoreRepository) UpdateTrainReservation(ctx context.Context, trainReservation *TrainReservation) error {
	_, err := r.client.Collection(trainReservationCollection).Doc(trainReservation.ID).Set(ctx, trainReservation)
	return err
}

func (r *firestoreRepository) IsTrainSeatAvailable(ctx context.Context, seatID string) (bool, error) {
	query := r.client.Collection(trainReservationCollection).
		Where("seat_id", "==", seatID).
		Where("status", "!=", TrainReservationStatusCancelled)

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
