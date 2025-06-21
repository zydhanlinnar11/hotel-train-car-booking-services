package event

// EventName mendefinisikan tipe untuk nama event yang valid
type EventName string

const (
	// Commands dari Order Service ke Partisipan
	CommandReserveRoom EventName = "booking.command.reserve.room"
	CommandReserveCar  EventName = "booking.command.reserve.car"
	CommandReserveSeat EventName = "booking.command.reserve.seat"

	// Events dari Partisipan ke Order Service
	RoomReserved             EventName = "booking.event.room.reserved"
	RoomBookingConfirmed     EventName = "booking.event.room.booking.confirmed"
	RoomReservationCancelled EventName = "booking.event.room.reservation.cancelled"
	RoomReservationFailed    EventName = "booking.event.room.failed"
	CarReserved              EventName = "booking.event.car.reserved"
	CarBookingConfirmed      EventName = "booking.event.car.booking.confirmed"
	CarReservationCancelled  EventName = "booking.event.car.reservation.cancelled"
	CarReservationFailed     EventName = "booking.event.car.failed"
	SeatBookingConfirmed     EventName = "booking.event.seat.booking.confirmed"
	SeatReserved             EventName = "booking.event.seat.reserved"
	SeatReservationCancelled EventName = "booking.event.seat.reservation.cancelled"
	SeatReservationFailed    EventName = "booking.event.seat.failed"

	// Commands Kompensasi dari Order Service
	CommandCancelRoom EventName = "booking.command.cancel.room"
	CommandCancelCar  EventName = "booking.command.cancel.car"
	CommandCancelSeat EventName = "booking.command.cancel.seat"

	// Event Final
	OrderBooked    EventName = "booking.event.order.booked"
	OrderFailed    EventName = "booking.event.order.failed"
	OrderCompleted EventName = "booking.event.order.completed"
)

// Message adalah struktur dasar untuk setiap pesan di RabbitMQ
type Message struct {
	EventName     EventName `json:"event_name"`
	CorrelationID string    `json:"correlation_id"` // Menggunakan OrderID
	Payload       any       `json:"payload"`
}

// ReserveRoomPayload adalah contoh payload untuk command reserve room
type ReserveRoomPayload struct {
	RoomID    string `json:"hotel_room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type ReserveCarPayload struct {
	CarID     string `json:"car_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type ReserveSeatPayload struct {
	SeatID string `json:"seat_id"`
}

type CancelRoomPayload struct {
	RoomReservationID string `json:"room_reservation_id"`
}

type CancelCarPayload struct {
	CarReservationID string `json:"car_reservation_id"`
}

type CancelSeatPayload struct {
	SeatReservationID string `json:"seat_reservation_id"`
}

type OrderBookedPayload struct {
	OrderID string `json:"order_id"`
}

type OrderFailedPayload struct {
	OrderID string `json:"order_id"`
}

type RoomReservedPayload struct {
	RoomReservationID string `json:"room_reservation_id"`
}

type CarReservedPayload struct {
	CarReservationID string `json:"car_reservation_id"`
}

type SeatReservedPayload struct {
	SeatReservationID string `json:"seat_reservation_id"`
}

type RoomBookingConfirmedPayload struct {
	RoomReservationID string `json:"room_reservation_id"`
}

type CarReservationCancelledPayload struct {
	CarReservationID string `json:"car_reservation_id"`
}

type SeatReservationCancelledPayload struct {
	SeatReservationID string `json:"seat_reservation_id"`
}

type RoomReservationFailedPayload struct {
	RoomReservationID string `json:"room_reservation_id"`
	FailureReason     string `json:"failure_reason"`
}

type CarReservationFailedPayload struct {
	CarReservationID string `json:"car_reservation_id"`
	FailureReason    string `json:"failure_reason"`
}

type SeatReservationFailedPayload struct {
	SeatReservationID string `json:"seat_reservation_id"`
	FailureReason     string `json:"failure_reason"`
}
