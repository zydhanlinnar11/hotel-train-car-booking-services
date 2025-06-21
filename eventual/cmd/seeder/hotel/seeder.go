package hotel

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/internal/hotel"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/utils"
)

var hotelBrands = []string{
	"Marriott Jakarta",
	"Ritz-Carlton Jakarta",
	"Mandarin Oriental Jakarta",
	"Four Seasons Jakarta",
	"Grand Hyatt Jakarta",
	"InterContinental Jakarta",
	"Sheraton Jakarta",
	"Pullman Jakarta",
	"Novotel Jakarta",
	"Ibis Jakarta",
	"Marriott Bandung",
	"Ritz-Carlton Bandung",
	"Marriott Surabaya",
	"Ritz-Carlton Surabaya",
	"Marriott Medan",
}

func Seed(ctx context.Context, client *firestore.Client) error {
	log.Println("Starting hotel room seeder...")

	collection := client.Collection("hotel_rooms")
	bw := client.BulkWriter(ctx)

	for _, hotelName := range hotelBrands {
		log.Printf("Seeding %s...", hotelName)

		// 5 floors, 20 units per floor
		for floor := 1; floor <= 5; floor++ {
			for unitNumber := 1; unitNumber <= 20; unitNumber++ {
				roomName := fmt.Sprintf("%d%02d", floor, unitNumber)
				roomID := utils.Slugify(fmt.Sprintf("%s-%s", hotelName, roomName))

				hotelRoom := hotel.HotelRoom{
					ID:        roomID,
					HotelName: hotelName,
					RoomName:  roomName,
				}

				docRef := collection.Doc(roomID)
				bw.Set(docRef, hotelRoom)
			}
		}
	}

	// Flush all writes
	bw.Flush()

	log.Printf("Hotel room seeder completed. Total rooms: %d", len(hotelBrands)*5*20)
	return nil
}
