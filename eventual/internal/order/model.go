package order

import "time"

// OrderStatus merepresentasikan status dari Saga
type OrderStatus string

type ReservationStatus string

const (
	StatusPending              OrderStatus = "PENDING"
	StatusAwaitingConfirmation OrderStatus = "AWAITING_CONFIRMATION"
	StatusBooked               OrderStatus = "BOOKED"
	StatusFailed               OrderStatus = "FAILED"
)

const (
	ReservationStatusPending ReservationStatus = "PENDING"
	ReservationStatusBooked  ReservationStatus = "BOOKED"
	ReservationStatusFailed  ReservationStatus = "FAILED"
)

// Order adalah representasi data order di Firestore
type Order struct {
	ID             string      `firestore:"id" json:"id"`
	UserID         string      `firestore:"user_id" json:"user_id"`
	Status         OrderStatus `firestore:"status" json:"status"`
	HotelRoomID    string      `firestore:"hotel_room_id" json:"hotel_room_id"`
	CarID          string      `firestore:"car_id" json:"car_id"`
	TrainSeatID    string      `firestore:"train_seat_id" json:"train_seat_id"`
	HotelStartDate string      `firestore:"hotel_start_date" json:"hotel_start_date"`
	HotelEndDate   string      `firestore:"hotel_end_date" json:"hotel_end_date"`
	CarStartDate   string      `firestore:"car_start_date" json:"car_start_date"`
	CarEndDate     string      `firestore:"car_end_date" json:"car_end_date"`

	// Reservation ID dari sub-transaksi
	HotelReservationID string `firestore:"hotel_reservation_id,omitempty" json:"hotel_reservation_id,omitempty"`
	CarReservationID   string `firestore:"car_reservation_id,omitempty" json:"car_reservation_id,omitempty"`
	TrainReservationID string `firestore:"train_reservation_id,omitempty" json:"train_reservation_id,omitempty"`

	// Status untuk setiap sub-transaksi
	HotelReservationStatus        ReservationStatus `firestore:"hotel_reservation_status" json:"hotel_reservation_status"`
	CarReservationStatus          ReservationStatus `firestore:"car_reservation_status" json:"car_reservation_status"`
	TrainReservationStatus        ReservationStatus `firestore:"train_reservation_status" json:"train_reservation_status"`
	HotelReservationFailureReason string            `firestore:"hotel_reservation_failure_reason,omitempty" json:"hotel_reservation_failure_reason,omitempty"`
	CarReservationFailureReason   string            `firestore:"car_reservation_failure_reason,omitempty" json:"car_reservation_failure_reason,omitempty"`
	TrainReservationFailureReason string            `firestore:"train_reservation_failure_reason,omitempty" json:"train_reservation_failure_reason,omitempty"`

	CreatedAt time.Time `firestore:"created_at" json:"created_at"`
	UpdatedAt time.Time `firestore:"updated_at" json:"updated_at"`
}
