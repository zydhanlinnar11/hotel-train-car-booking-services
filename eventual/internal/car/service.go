package car

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
	case event.CommandReserveCar:
		return s.handleReserveCar(ctx, msg)
	case event.CommandCancelCar:
		return s.handleCancelCar(ctx, msg)
	}

	return nil
}

func (s *service) publishErrorEvent(ctx context.Context, msg event.Message, err error) error {
	s.publisher.Publish(ctx, string(event.CarReservationFailed), event.Message{
		EventName:     event.CarReservationFailed,
		CorrelationID: msg.CorrelationID,
		Payload:       event.CarReservationFailedPayload{FailureReason: err.Error()},
	})

	return err
}

func (s *service) handleReserveCar(ctx context.Context, msg event.Message) error {
	payload, err := mapToPayload[event.ReserveCarPayload](msg)
	if err != nil {
		return s.publishErrorEvent(ctx, msg, err)
	}

	isAvailable, err := s.repo.IsCarAvailable(ctx, payload.CarID, payload.StartDate, payload.EndDate)
	if err != nil {
		return s.publishErrorEvent(ctx, msg, err)
	}
	if !isAvailable {
		return s.publishErrorEvent(ctx, msg, errors.New("car is not available"))
	}

	car, err := s.repo.GetCarByID(ctx, payload.CarID)
	if err != nil {
		return s.publishErrorEvent(ctx, msg, err)
	}

	carReservation := &CarReservation{
		ID:          uuid.NewString(),
		CarID:       car.CarID,
		CarName:     car.CarName,
		CarBrand:    car.Brand,
		CarModel:    car.Model,
		CarYear:     car.Year,
		CarPrice:    car.Price,
		CarLocation: car.Location,
		StartDate:   payload.StartDate,
		EndDate:     payload.EndDate,
		OrderID:     msg.CorrelationID,
		Status:      CarStatusReserved,
	}

	if err := s.repo.CreateCarReservation(ctx, carReservation); err != nil {
		return s.publishErrorEvent(ctx, msg, err)
	}

	s.publisher.Publish(ctx, string(event.CarReserved), event.Message{
		EventName:     event.CarReserved,
		CorrelationID: msg.CorrelationID,
		Payload: event.CarReservedPayload{
			CarReservationID: carReservation.ID,
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

func (s *service) handleCancelCar(ctx context.Context, msg event.Message) error {
	payload, err := mapToPayload[event.CancelCarPayload](msg)
	if err != nil {
		return s.publishErrorEvent(ctx, msg, err)
	}

	carReservation, err := s.repo.GetCarReservationByID(ctx, payload.CarReservationID)
	if err != nil {
		return s.publishErrorEvent(ctx, msg, err)
	}

	carReservation.Status = CarStatusAvailable
	if err := s.repo.UpdateCarReservation(ctx, carReservation); err != nil {
		return s.publishErrorEvent(ctx, msg, err)
	}

	return nil
}
