package hotel

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
	case event.CommandReserveRoom:
		return s.handleReserveRoom(ctx, msg)
	case event.CommandCancelRoom:
		return s.handleCancelRoom(ctx, msg)
	}

	return nil
}

func (s *service) publishErrorEvent(ctx context.Context, msg event.Message, err error) error {
	s.publisher.Publish(ctx, string(event.RoomReservationFailed), event.Message{
		EventName:     event.RoomReservationFailed,
		CorrelationID: msg.CorrelationID,
		Payload:       event.RoomReservationFailedPayload{FailureReason: err.Error()},
	})

	return err
}

func (s *service) handleReserveRoom(ctx context.Context, msg event.Message) error {
	payload, err := mapToPayload[event.ReserveRoomPayload](msg)
	if err != nil {
		return s.publishErrorEvent(ctx, msg, err)
	}

	isAvailable, err := s.repo.IsHotelRoomAvailable(ctx, payload.RoomID, payload.StartDate, payload.EndDate)
	if err != nil {
		return s.publishErrorEvent(ctx, msg, err)
	}
	if !isAvailable {
		return s.publishErrorEvent(ctx, msg, errors.New("hotel room is not available"))
	}

	hotelRoom, err := s.repo.GetHotelRoomByID(ctx, payload.RoomID)
	if err != nil {
		return s.publishErrorEvent(ctx, msg, err)
	}

	hotelReservation := &HotelReservation{
		ID:                 uuid.NewString(),
		HotelRoomID:        hotelRoom.ID,
		HotelRoomName:      hotelRoom.RoomName,
		HotelName:          hotelRoom.HotelName,
		HotelRoomStartDate: payload.StartDate,
		HotelRoomEndDate:   payload.EndDate,
		OrderID:            msg.CorrelationID,
		Status:             HotelRoomReservationStatusReserved,
	}

	if err := s.repo.CreateHotelReservation(ctx, hotelReservation); err != nil {
		return s.publishErrorEvent(ctx, msg, err)
	}

	s.publisher.Publish(ctx, string(event.RoomReserved), event.Message{
		EventName:     event.RoomReserved,
		CorrelationID: msg.CorrelationID,
		Payload: event.RoomReservedPayload{
			RoomReservationID: hotelReservation.ID,
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

func (s *service) handleCancelRoom(ctx context.Context, msg event.Message) error {
	payload, err := mapToPayload[event.CancelRoomPayload](msg)
	if err != nil {
		return s.publishErrorEvent(ctx, msg, err)
	}

	hotelReservation, err := s.repo.GetHotelReservationByOrderID(ctx, payload.OrderID)
	if errors.Is(err, ErrHotelReservationNotFound) {
		return nil
	}
	if err != nil {
		return s.publishErrorEvent(ctx, msg, err)
	}

	hotelReservation.Status = HotelRoomReservationStatusCancelled
	if err := s.repo.UpdateHotelReservation(ctx, hotelReservation); err != nil {
		return s.publishErrorEvent(ctx, msg, err)
	}

	return nil
}
