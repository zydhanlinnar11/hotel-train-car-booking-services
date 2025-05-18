package infrastructures

import "time"

const (
	OrderCollection = "order_orders"
	EventCollection = "order_events"
)

type Order struct {
	Id                 string `firestore:"id"`
	HotelRoomId        string `firestore:"hotel_room_id"`
	HotelRoomStartDate string `firestore:"hotel_room_start_date"`
	HotelRoomEndDate   string `firestore:"hotel_room_end_date"`
	CarId              string `firestore:"car_id"`
	CarStartDate       string `firestore:"car_start_date"`
	CarEndDate         string `firestore:"car_end_date"`
	TrainSeatId        string `firestore:"train_seat_id"`
	FailureCause       int    `firestore:"failure_cause"`
	Version            int    `firestore:"version"`
	Status             int    `firestore:"status"`
}

type OrderPlacedEventPayload struct {
	OrderId            string `firestore:"id"`
	HotelRoomId        string `firestore:"hotel_room_id"`
	HotelRoomStartDate string `firestore:"hotel_room_start_date"`
	HotelRoomEndDate   string `firestore:"hotel_room_end_date"`
	CarId              string `firestore:"car_id"`
	CarStartDate       string `firestore:"car_start_date"`
	CarEndDate         string `firestore:"car_end_date"`
	TrainSeatId        string `firestore:"train_seat_id"`
}

type OrderBookedEventPayload struct {
	OrderId string `firestore:"id"`
}

type OrderCompletedEventPayload struct {
	OrderId string `firestore:"id"`
}

type OrderReleasedEventPayload struct {
	OrderId string `firestore:"id"`
	Cause   int    `firestore:"cause"`
}

type OrderFailedEventPayload struct {
	OrderId string `firestore:"id"`
	Cause   int    `firestore:"cause"`
}

type Event struct {
	Name    string    `firestore:"name"`
	Payload any       `firestore:"payload"`
	Time    time.Time `firestore:"time"`
}
