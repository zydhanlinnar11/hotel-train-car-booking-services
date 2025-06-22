package hotel

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zydhanlinnar11/hotel-train-car-booking-services/twophase/internal/hotel"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/twophase/pkg/utils"
)

var hotelBrands = []string{
	"Marriott Jakarta",
	"Ritz-Carlton Jakarta",
	"Mandarin Oriental Jakarta",
	"Four Seasons Jakarta",
	"Grand Hyatt Jakarta",
	// "InterContinental Jakarta",
	// "Sheraton Jakarta",
	// "Pullman Jakarta",
	// "Novotel Jakarta",
	// "Ibis Jakarta",
	// "Marriott Bandung",
	// "Ritz-Carlton Bandung",
	// "Marriott Surabaya",
	// "Ritz-Carlton Surabaya",
	// "Marriott Medan",
}

func Seed(ctx context.Context, repo *hotel.Repository) error {
	log.Println("Starting hotel room availability seeder...")

	hotelRoomAvailabilities := make([]hotel.HotelRoomAvailability, 0)

	for _, hotelName := range hotelBrands {
		log.Printf("Seeding %s...", hotelName)

		// 5 floors, 20 units per floor
		for floor := 1; floor <= 5; floor++ {
			for unitNumber := 1; unitNumber <= 20; unitNumber++ {
				roomName := fmt.Sprintf("%d%02d", floor, unitNumber)
				roomID := utils.Slugify(fmt.Sprintf("%s-%s", hotelName, roomName))

				// 10 days of availability
				for i := 0; i < 10; i++ {
					date := time.Now().AddDate(0, 0, i)
					hotelRoom := hotel.HotelRoomAvailability{
						RoomID:    roomID,
						HotelName: hotelName,
						RoomName:  roomName,
						Date:      date.Format(time.DateOnly),
						Available: true,
					}

					hotelRoomAvailabilities = append(hotelRoomAvailabilities, hotelRoom)
				}
			}
		}
	}

	if err := repo.BulkWriteHotelRoomAvailability(ctx, hotelRoomAvailabilities); err != nil {
		return fmt.Errorf("failed to bulk write hotel room availability: %w", err)
	}

	log.Printf("Hotel room availability seeder completed. Total rooms: %d", len(hotelRoomAvailabilities))
	return nil
}
