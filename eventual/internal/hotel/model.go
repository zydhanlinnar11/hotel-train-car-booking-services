package hotel

type HotelRoomReservationStatus string

const (
	HotelRoomReservationStatusCancelled HotelRoomReservationStatus = "CANCELLED"
	HotelRoomReservationStatusReserved  HotelRoomReservationStatus = "RESERVED"
)

type HotelRoom struct {
	ID        string `firestore:"id" json:"id"`
	HotelName string `firestore:"hotel_name" json:"hotel_name"`
	RoomName  string `firestore:"room_name" json:"room_name"`
}

type HotelReservation struct {
	ID                 string                     `firestore:"id" json:"id"`
	HotelRoomID        string                     `firestore:"hotel_room_id" json:"hotel_room_id"`
	HotelRoomName      string                     `firestore:"hotel_room_name" json:"hotel_room_name"`
	HotelName          string                     `firestore:"hotel_name" json:"hotel_name"`
	HotelRoomStartDate string                     `firestore:"hotel_room_start_date" json:"hotel_room_start_date"`
	HotelRoomEndDate   string                     `firestore:"hotel_room_end_date" json:"hotel_room_end_date"`
	OrderID            string                     `firestore:"order_id" json:"order_id"`
	Status             HotelRoomReservationStatus `firestore:"status" json:"status"`
}
