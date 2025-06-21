package train

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/event"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/messagebus"
)

type Service interface {
	ProcessSagaEvent(ctx context.Context, msg event.Message) error
}

type service struct {
	repo      Repository
	publisher messagebus.Publisher
}

func NewService(repo Repository, publisher messagebus.Publisher) Service {
	return &service{repo: repo, publisher: publisher}
}

func (s *service) ProcessSagaEvent(ctx context.Context, msg event.Message) error {
	log.Println("Received saga event", msg.EventName)
	switch msg.EventName {
	case event.CommandReserveSeat:
		return s.handleReserveSeat(ctx, msg)
	case event.CommandCancelSeat:
		return s.handleCancelSeat(ctx, msg)
	}

	return nil
}

func (s *service) publishErrorEvent(ctx context.Context, msg event.Message, err error) error {
	s.publisher.Publish(ctx, string(event.SeatReservationFailed), event.Message{
		EventName:     event.SeatReservationFailed,
		CorrelationID: msg.CorrelationID,
		Payload:       event.SeatReservationFailedPayload{FailureReason: err.Error()},
	})

	return err
}

func (s *service) handleReserveSeat(ctx context.Context, msg event.Message) error {
	payload, err := mapToPayload[event.ReserveSeatPayload](msg)
	if err != nil {
		return s.publishErrorEvent(ctx, msg, err)
	}

	isAvailable, err := s.repo.IsTrainSeatAvailable(ctx, payload.SeatID)
	if err != nil {
		return s.publishErrorEvent(ctx, msg, err)
	}
	if !isAvailable {
		return s.publishErrorEvent(ctx, msg, errors.New("train seat is not available"))
	}

	trainSeat, err := s.repo.GetTrainSeatByID(ctx, payload.SeatID)
	if err != nil {
		return s.publishErrorEvent(ctx, msg, err)
	}

	trainReservation := &TrainReservation{
		ID:            uuid.NewString(),
		SeatID:        trainSeat.SeatID,
		TrainNumber:   trainSeat.TrainNumber,
		CarNumber:     trainSeat.CarNumber,
		SeatNumber:    trainSeat.SeatNumber,
		Class:         trainSeat.Class,
		Price:         trainSeat.Price,
		FromStation:   trainSeat.FromStation,
		ToStation:     trainSeat.ToStation,
		DepartureTime: trainSeat.DepartureTime,
		ArrivalTime:   trainSeat.ArrivalTime,
		OrderID:       msg.CorrelationID,
		Status:        TrainReservationStatusReserved,
	}

	if err := s.repo.CreateTrainReservation(ctx, trainReservation); err != nil {
		return s.publishErrorEvent(ctx, msg, err)
	}

	s.publisher.Publish(ctx, string(event.SeatReserved), event.Message{
		EventName:     event.SeatReserved,
		CorrelationID: msg.CorrelationID,
		Payload: event.SeatReservedPayload{
			SeatReservationID: trainReservation.ID,
		},
	})

	return nil
}

func mapToPayload[T any](msg event.Message) (T, error) {
	var payload T
	marshalledPayload, err := json.Marshal(msg.Payload)
	if err != nil {
		return payload, err
	}
	if err := json.Unmarshal(marshalledPayload, &payload); err != nil {
		return payload, err
	}
	return payload, nil
}

func (s *service) handleCancelSeat(ctx context.Context, msg event.Message) error {
	payload, err := mapToPayload[event.CancelSeatPayload](msg)
	if err != nil {
		return s.publishErrorEvent(ctx, msg, err)
	}

	trainReservation, err := s.repo.GetTrainReservationByOrderID(ctx, payload.OrderID)
	if err != nil {
		return s.publishErrorEvent(ctx, msg, err)
	}

	trainReservation.Status = TrainReservationStatusCancelled
	if err := s.repo.UpdateTrainReservation(ctx, trainReservation); err != nil {
		return s.publishErrorEvent(ctx, msg, err)
	}

	return nil
}
