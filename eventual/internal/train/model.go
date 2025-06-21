package train

type TrainReservationStatus string

const (
	TrainReservationStatusCancelled TrainReservationStatus = "CANCELLED"
	TrainReservationStatusReserved  TrainReservationStatus = "RESERVED"
)

type TrainSeat struct {
	ID        string `firestore:"id" json:"id"`
	SeatID    string `firestore:"seat_id" json:"seat_id"`
	TrainName string `firestore:"train_name" json:"train_name"`
}

type TrainReservation struct {
	ID        string                 `firestore:"id" json:"id"`
	SeatID    string                 `firestore:"seat_id" json:"seat_id"`
	TrainName string                 `firestore:"train_name" json:"train_name"`
	OrderID   string                 `firestore:"order_id" json:"order_id"`
	Status    TrainReservationStatus `firestore:"status" json:"status"`
}
