package order

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/config"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/event"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/messagebus"
)

type CreateOrderPayload struct {
	HotelRoomID        string `json:"hotel_room_id" binding:"required"`
	HotelRoomStartDate string `json:"hotel_room_start_date" binding:"required"`
	HotelRoomEndDate   string `json:"hotel_room_end_date" binding:"required"`
	CarID              string `json:"car_id" binding:"required"`
	CarStartDate       string `json:"car_start_date" binding:"required"`
	CarEndDate         string `json:"car_end_date" binding:"required"`
	TrainSeatID        string `json:"train_seat_id" binding:"required"`
	UserID             string `json:"user_id" binding:"required"`
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

func NewService(repo Repository, publisher messagebus.Publisher) Service {
	return &service{repo: repo, publisher: publisher}
}

func (s *service) parseDate(hotelStartDateStr, hotelEndDateStr, carStartDateStr, carEndDateStr string) (time.Time, time.Time, time.Time, time.Time, error) {
	hotelStartDate, err := time.Parse(config.DateFormat, hotelStartDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}, err
	}
	hotelEndDate, err := time.Parse(config.DateFormat, hotelEndDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}, err
	}
	carStartDate, err := time.Parse(config.DateFormat, carStartDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}, err
	}
	carEndDate, err := time.Parse(config.DateFormat, carEndDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}, err
	}

	return hotelStartDate, hotelEndDate, carStartDate, carEndDate, nil
}

func (s *service) StartSaga(ctx context.Context, payload CreateOrderPayload) (*Order, error) {
	hotelStartDate, hotelEndDate, carStartDate, carEndDate, err := s.parseDate(
		payload.HotelRoomStartDate,
		payload.HotelRoomEndDate,
		payload.CarStartDate,
		payload.CarEndDate,
	)
	if err != nil {
		return nil, err
	}

	// 1. Buat Order baru dengan status PENDING
	order := &Order{
		ID:     ulid.Make().String(),
		UserID: payload.UserID,
		Status: StatusPending,

		HotelRoomID:    payload.HotelRoomID,
		CarID:          payload.CarID,
		TrainSeatID:    payload.TrainSeatID,
		HotelStartDate: hotelStartDate.Format(config.DateFormat),
		HotelEndDate:   hotelEndDate.Format(config.DateFormat),
		CarStartDate:   carStartDate.Format(config.DateFormat),
		CarEndDate:     carEndDate.Format(config.DateFormat),

		HotelReservationStatus: ReservationStatusPending,
		CarReservationStatus:   ReservationStatusPending,
		TrainReservationStatus: ReservationStatusPending,

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
	log.Println("Received saga event", msg.EventName)
	// 1. Ambil order dari DB menggunakan msg.CorrelationID
	order, err := s.repo.GetOrderByID(ctx, msg.CorrelationID)
	if err != nil {
		return err
	}

	// 2. State machine untuk memproses event berdasarkan EventName
	switch msg.EventName {
	case event.RoomReserved:
		var payload event.RoomReservedPayload
		if err := s.unmarshalPayload(msg.Payload, &payload); err != nil {
			return err
		}
		order.HotelReservationStatus = ReservationStatusBooked
		order.HotelReservationID = payload.RoomReservationID

	case event.CarReserved:
		var payload event.CarReservedPayload
		if err := s.unmarshalPayload(msg.Payload, &payload); err != nil {
			return err
		}
		order.CarReservationStatus = ReservationStatusBooked
		order.CarReservationID = payload.CarReservationID

	case event.SeatReserved:
		var payload event.SeatReservedPayload
		if err := s.unmarshalPayload(msg.Payload, &payload); err != nil {
			return err
		}
		order.TrainReservationStatus = ReservationStatusBooked
		order.TrainReservationID = payload.SeatReservationID

	case event.RoomReservationFailed:
		var payload event.RoomReservationFailedPayload
		if err := s.unmarshalPayload(msg.Payload, &payload); err != nil {
			return err
		}
		order.HotelReservationStatus = ReservationStatusFailed
		order.HotelReservationFailureReason = payload.FailureReason

	case event.CarReservationFailed:
		var payload event.CarReservationFailedPayload
		if err := s.unmarshalPayload(msg.Payload, &payload); err != nil {
			return err
		}
		order.CarReservationStatus = ReservationStatusFailed
		order.CarReservationFailureReason = payload.FailureReason

	case event.SeatReservationFailed:
		var payload event.SeatReservationFailedPayload
		if err := s.unmarshalPayload(msg.Payload, &payload); err != nil {
			return err
		}
		order.TrainReservationStatus = ReservationStatusFailed
		order.TrainReservationFailureReason = payload.FailureReason
	}

	// 3. Cek apakah ada yang pending
	if order.HotelReservationStatus == ReservationStatusPending || order.CarReservationStatus == ReservationStatusPending || order.TrainReservationStatus == ReservationStatusPending {
		return s.repo.UpdateOrder(ctx, order)
	}

	// 4. Cek apakah semua reservasi sudah berhasil
	if order.HotelReservationStatus == ReservationStatusBooked && order.CarReservationStatus == ReservationStatusBooked && order.TrainReservationStatus == ReservationStatusBooked {
		order.Status = StatusBooked
		s.repo.UpdateOrder(ctx, order)
		s.publisher.Publish(ctx, string(event.OrderBooked), event.Message{
			EventName:     event.OrderBooked,
			CorrelationID: order.ID,
			Payload:       event.OrderBookedPayload{OrderID: order.ID},
		})
		return nil
	}

	return s.startCompensation(ctx, order)
}

// unmarshalPayload adalah helper function untuk unmarshal JSON payload
func (s *service) unmarshalPayload(payload any, target any) error {
	// Convert payload to JSON bytes first
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	// Then unmarshal to target struct
	return json.Unmarshal(jsonBytes, target)
}

func (s *service) startCompensation(ctx context.Context, order *Order) error {
	order.Status = StatusFailed
	s.repo.UpdateOrder(ctx, order)

	// Kirim command kompensasi untuk command yang sudah dikirim
	// Menggunakan OrderID saja karena relasi one-to-one
	s.publisher.Publish(ctx, string(event.CommandCancelRoom), event.Message{
		EventName:     event.CommandCancelRoom,
		CorrelationID: order.ID,
		Payload:       event.CancelRoomPayload{OrderID: order.ID},
	})

	s.publisher.Publish(ctx, string(event.CommandCancelCar), event.Message{
		EventName:     event.CommandCancelCar,
		CorrelationID: order.ID,
		Payload:       event.CancelCarPayload{OrderID: order.ID},
	})
	s.publisher.Publish(ctx, string(event.CommandCancelSeat), event.Message{
		EventName:     event.CommandCancelSeat,
		CorrelationID: order.ID,
		Payload:       event.CancelSeatPayload{OrderID: order.ID},
	})

	// Publish event final ORDER_FAILED
	s.publisher.Publish(ctx, string(event.OrderFailed), event.Message{
		EventName:     event.OrderFailed,
		CorrelationID: order.ID,
		Payload:       event.OrderFailedPayload{OrderID: order.ID},
	})

	return nil
}
