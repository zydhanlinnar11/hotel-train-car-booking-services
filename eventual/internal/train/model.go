package train

type TrainReservationStatus string

const (
	TrainReservationStatusCancelled TrainReservationStatus = "CANCELLED"
	TrainReservationStatusReserved  TrainReservationStatus = "RESERVED"
	TrainReservationStatusBooked    TrainReservationStatus = "BOOKED"
)

type TrainSeat struct {
	ID            string `firestore:"id" json:"id"`
	SeatID        string `firestore:"seat_id" json:"seat_id"`
	TrainNumber   string `firestore:"train_number" json:"train_number"`
	CarNumber     string `firestore:"car_number" json:"car_number"`
	SeatNumber    string `firestore:"seat_number" json:"seat_number"`
	Class         string `firestore:"class" json:"class"`
	Price         int    `firestore:"price" json:"price"`
	FromStation   string `firestore:"from_station" json:"from_station"`
	ToStation     string `firestore:"to_station" json:"to_station"`
	DepartureTime string `firestore:"departure_time" json:"departure_time"`
	ArrivalTime   string `firestore:"arrival_time" json:"arrival_time"`
}

type TrainReservation struct {
	ID            string                 `firestore:"id" json:"id"`
	SeatID        string                 `firestore:"seat_id" json:"seat_id"`
	TrainNumber   string                 `firestore:"train_number" json:"train_number"`
	CarNumber     string                 `firestore:"car_number" json:"car_number"`
	SeatNumber    string                 `firestore:"seat_number" json:"seat_number"`
	Class         string                 `firestore:"class" json:"class"`
	Price         int                    `firestore:"price" json:"price"`
	FromStation   string                 `firestore:"from_station" json:"from_station"`
	ToStation     string                 `firestore:"to_station" json:"to_station"`
	DepartureTime string                 `firestore:"departure_time" json:"departure_time"`
	ArrivalTime   string                 `firestore:"arrival_time" json:"arrival_time"`
	OrderID       string                 `firestore:"order_id" json:"order_id"`
	Status        TrainReservationStatus `firestore:"status" json:"status"`
}
