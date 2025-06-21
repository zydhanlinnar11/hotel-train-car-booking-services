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
	ID             string      `firestore:"id,omitempty"`
	UserID         string      `firestore:"user_id"`
	Status         OrderStatus `firestore:"status"`
	HotelID        string      `firestore:"hotel_id"`
	CarID          string      `firestore:"car_id"`
	TrainSeatID    string      `firestore:"train_seat_id"`
	HotelStartDate time.Time   `firestore:"hotel_start_date"`
	HotelEndDate   time.Time   `firestore:"hotel_end_date"`
	CarStartDate   time.Time   `firestore:"car_start_date"`
	CarEndDate     time.Time   `firestore:"car_end_date"`

	// Reservation ID dari sub-transaksi
	HotelReservationID string `firestore:"hotel_reservation_id,omitempty"`
	CarReservationID   string `firestore:"car_reservation_id,omitempty"`
	TrainReservationID string `firestore:"train_reservation_id,omitempty"`
	FailureReason      string `firestore:"failure_reason,omitempty"`

	// Status untuk setiap sub-transaksi
	IsRoomReserved bool `firestore:"is_room_reserved"`
	IsCarReserved  bool `firestore:"is_car_reserved"`
	IsSeatReserved bool `firestore:"is_seat_reserved"`

	Version   int       `firestore:"version"`
	CreatedAt time.Time `firestore:"created_at"`
	UpdatedAt time.Time `firestore:"updated_at"`
}
