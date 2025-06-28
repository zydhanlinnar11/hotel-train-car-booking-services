package hotel

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/utils"
)

func ExportToCSV(filename string) error {
	log.Println("Starting hotel room CSV export...")

	// Create CSV file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"ID", "HotelName", "RoomName"}); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Generate the same data as seeder but write to CSV
	for _, hotelName := range hotelBrands {
		log.Printf("Exporting %s...", hotelName)

		// 5 floors, 20 units per floor
		for floor := 1; floor <= 5; floor++ {
			for unitNumber := 1; unitNumber <= 20; unitNumber++ {
				roomName := fmt.Sprintf("%d%02d", floor, unitNumber)
				roomID := utils.Slugify(fmt.Sprintf("%s-%s", hotelName, roomName))

				// Write to CSV
				if err := writer.Write([]string{roomID, hotelName, roomName}); err != nil {
					return fmt.Errorf("failed to write row: %w", err)
				}
			}
		}
	}

	log.Printf("Hotel room CSV export completed. Total rooms: %d", len(hotelBrands)*5*20)
	return nil
}
