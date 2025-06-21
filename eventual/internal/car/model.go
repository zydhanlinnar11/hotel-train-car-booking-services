package car

type CarReservationStatus string

const (
	CarReservationStatusCancelled CarReservationStatus = "CANCELLED"
	CarReservationStatusReserved  CarReservationStatus = "RESERVED"
)

type Car struct {
	ID   string `firestore:"id" json:"id"`
	Name string `firestore:"name" json:"name"`
}

type CarReservation struct {
	ID        string               `firestore:"id" json:"id"`
	CarID     string               `firestore:"car_id" json:"car_id"`
	CarName   string               `firestore:"car_name" json:"car_name"`
	StartDate string               `firestore:"start_date" json:"start_date"`
	EndDate   string               `firestore:"end_date" json:"end_date"`
	OrderID   string               `firestore:"order_id" json:"order_id"`
	Status    CarReservationStatus `firestore:"status" json:"status"`
}
