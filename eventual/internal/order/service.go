package order

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/event"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/messagebus"
)

type CreateOrderPayload struct {
	HotelRoomID        string `json:"hotel_room_id"`
	HotelRoomStartDate string `json:"hotel_room_start_date"`
	HotelRoomEndDate   string `json:"hotel_room_end_date"`
	CarID              string `json:"car_id"`
	CarStartDate       string `json:"car_start_date"`
	CarEndDate         string `json:"car_end_date"`
	TrainSeatID        string `json:"train_seat_id"`
	UserID             string `json:"user_id"`
}

// Service mendefinisikan logika bisnis untuk Order Service
type Service interface {
	// StartSaga dipanggil oleh HTTP handler untuk memulai proses booking
	StartSaga(ctx context.Context, payload CreateOrderPayload) (*Order, error)

	// ProcessSagaEvent dipanggil oleh event handler saat menerima balasan dari service lain
	ProcessSagaEvent(ctx context.Context, msg event.Message) error
}

type service struct {
	repo      Repository
	publisher messagebus.Publisher
}

func (s *service) StartSaga(ctx context.Context, payload CreateOrderPayload) (*Order, error) {
	// 1. Buat Order baru dengan status PENDING
	order := &Order{
		ID:     uuid.NewString(),
		UserID: payload.UserID,
		Status: StatusPending,

		Version:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	// 2. Publish command untuk setiap layanan partisipan
	//    Gunakan CorrelationID yang sama dengan order.ID
	s.publisher.Publish(ctx, string(event.CommandReserveRoom), event.Message{
		EventName:     event.CommandReserveRoom,
		CorrelationID: order.ID,
		Payload: event.ReserveRoomPayload{
			RoomID:    payload.HotelRoomID,
			StartDate: payload.HotelRoomStartDate,
			EndDate:   payload.HotelRoomEndDate,
		},
	})
	s.publisher.Publish(ctx, string(event.CommandReserveCar), event.Message{
		EventName:     event.CommandReserveCar,
		CorrelationID: order.ID,
		Payload: event.ReserveCarPayload{
			CarID:     payload.CarID,
			StartDate: payload.CarStartDate,
			EndDate:   payload.CarEndDate,
		},
	})
	s.publisher.Publish(ctx, string(event.CommandReserveSeat), event.Message{
		EventName:     event.CommandReserveSeat,
		CorrelationID: order.ID,
		Payload: event.ReserveSeatPayload{
			SeatID: payload.TrainSeatID,
		},
	})

	// 3. Update status order menjadi AWAITING_CONFIRMATION
	order.Status = StatusAwaitingConfirmation
	s.repo.UpdateOrder(ctx, order)

	return order, nil
}

func (s *service) ProcessSagaEvent(ctx context.Context, msg event.Message) error {
	// 1. Ambil order dari DB menggunakan msg.CorrelationID
	order, err := s.repo.GetOrderByID(ctx, msg.CorrelationID)
	if err != nil {
		return err
	}

	// 2. State machine untuk memproses event
	switch payload := msg.Payload.(type) {
	case event.RoomReservedPayload:
		order.IsRoomReserved = true
		order.HotelReservationID = payload.RoomReservationID
	case event.CarReservedPayload:
		order.IsCarReserved = true
		order.CarReservationID = payload.CarReservationID
	case event.SeatReservedPayload:
		order.IsSeatReserved = true
		order.TrainReservationID = payload.SeatReservationID
	case event.RoomReservationFailedPayload:
		return s.startCompensation(ctx, order, payload.FailureReason)
	case event.CarReservationFailedPayload:
		return s.startCompensation(ctx, order, payload.FailureReason)
	case event.SeatReservationFailedPayload:
		return s.startCompensation(ctx, order, payload.FailureReason)
	}

	// 3. Cek apakah semua reservasi sudah berhasil
	if order.IsRoomReserved && order.IsCarReserved && order.IsSeatReserved {
		// Jika semua berhasil, finalisasi Saga
		order.Status = StatusBooked
		s.repo.UpdateOrder(ctx, order)
		// Publish event final ORDER_BOOKED
		s.publisher.Publish(ctx, string(event.OrderBooked), event.Message{
			EventName:     event.OrderBooked,
			CorrelationID: order.ID,
			Payload:       event.OrderBookedPayload{OrderID: order.ID},
		})
		return nil
	}
	// Jika belum semua, simpan state terbaru
	s.repo.UpdateOrder(ctx, order)

	return nil
}

func (s *service) startCompensation(ctx context.Context, order *Order, reason string) error {
	order.Status = StatusFailed
	order.FailureReason = reason
	s.repo.UpdateOrder(ctx, order)

	// Kirim command kompensasi untuk reservasi yang sudah berhasil
	if order.IsRoomReserved {
		s.publisher.Publish(ctx, string(event.CommandCancelRoom), event.Message{
			EventName:     event.CommandCancelRoom,
			CorrelationID: order.ID,
			Payload:       event.CancelRoomPayload{RoomReservationID: order.HotelReservationID},
		})
	}
	if order.IsCarReserved {
		s.publisher.Publish(ctx, string(event.CommandCancelCar), event.Message{
			EventName:     event.CommandCancelCar,
			CorrelationID: order.ID,
			Payload:       event.CancelCarPayload{CarReservationID: order.CarReservationID},
		})
	}
	if order.IsSeatReserved {
		s.publisher.Publish(ctx, string(event.CommandCancelSeat), event.Message{
			EventName:     event.CommandCancelSeat,
			CorrelationID: order.ID,
			Payload:       event.CancelSeatPayload{SeatReservationID: order.TrainReservationID},
		})
	}

	return nil
}
