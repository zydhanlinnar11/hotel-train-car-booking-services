package car

type CarStatus string

const (
	CarStatusAvailable CarStatus = "AVAILABLE"
	CarStatusReserved  CarStatus = "RESERVED"
	CarStatusBooked    CarStatus = "BOOKED"
)

type Car struct {
	ID       string `firestore:"id" json:"id"`
	CarID    string `firestore:"car_id" json:"car_id"`
	CarName  string `firestore:"car_name" json:"car_name"`
	Brand    string `firestore:"brand" json:"brand"`
	Model    string `firestore:"model" json:"model"`
	Year     int    `firestore:"year" json:"year"`
	Price    int    `firestore:"price" json:"price"`
	Location string `firestore:"location" json:"location"`
}

type CarReservation struct {
	ID          string    `firestore:"id" json:"id"`
	CarID       string    `firestore:"car_id" json:"car_id"`
	CarName     string    `firestore:"car_name" json:"car_name"`
	CarBrand    string    `firestore:"car_brand" json:"car_brand"`
	CarModel    string    `firestore:"car_model" json:"car_model"`
	CarYear     int       `firestore:"car_year" json:"car_year"`
	CarPrice    int       `firestore:"car_price" json:"car_price"`
	CarLocation string    `firestore:"car_location" json:"car_location"`
	StartDate   string    `firestore:"start_date" json:"start_date"`
	EndDate     string    `firestore:"end_date" json:"end_date"`
	OrderID     string    `firestore:"order_id" json:"order_id"`
	Status      CarStatus `firestore:"status" json:"status"`
}
