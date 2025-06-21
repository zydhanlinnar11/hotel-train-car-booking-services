package order

import "time"

// OrderStatus merepresentasikan status dari Saga
type OrderStatus string

const (
	StatusPending              OrderStatus = "PENDING"
	StatusAwaitingConfirmation OrderStatus = "AWAITING_CONFIRMATION"
	StatusBooked               OrderStatus = "BOOKED"
	StatusFailed               OrderStatus = "FAILED"
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
	FailureReason      string `firestore:"failure_reason,omitempty" json:"failure_reason,omitempty"`

	// Status untuk setiap sub-transaksi
	IsRoomReserved bool `firestore:"is_room_reserved" json:"is_room_reserved"`
	IsCarReserved  bool `firestore:"is_car_reserved" json:"is_car_reserved"`
	IsSeatReserved bool `firestore:"is_seat_reserved" json:"is_seat_reserved"`

	CreatedAt time.Time `firestore:"created_at" json:"created_at"`
	UpdatedAt time.Time `firestore:"updated_at" json:"updated_at"`
}
