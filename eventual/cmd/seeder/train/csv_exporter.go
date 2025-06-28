package train

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/utils"
)

func ExportToCSV(filename string) error {
	log.Println("Starting train CSV export...")

	// Create CSV file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"ID", "SeatID", "TrainName"}); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Generate the same data as seeder but write to CSV
	for _, trainName := range trainNames {
		log.Printf("Exporting %s...", trainName)

		// 500 seats per train (10 trains × 500 seats = 5000 total)
		for seatNumber := 1; seatNumber <= 500; seatNumber++ {
			seatID := fmt.Sprintf("%d", seatNumber)
			seatDocID := utils.Slugify(fmt.Sprintf("%s-%s", trainName, seatID))

			// Write to CSV
			if err := writer.Write([]string{seatDocID, seatID, trainName}); err != nil {
				return fmt.Errorf("failed to write row: %w", err)
			}
		}
	}

	log.Printf("Train CSV export completed. Total seats: %d", len(trainNames)*500)
	return nil
}
