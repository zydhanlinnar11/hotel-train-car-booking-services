package car

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/utils"
)

func ExportToCSV(filename string) error {
	log.Println("Starting car CSV export...")

	// Create CSV file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"ID", "Name"}); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Generate the same data as seeder but write to CSV
	for _, brandData := range carBrands {
		for _, model := range brandData.models {
			log.Printf("Exporting %s %s...", brandData.brand, model)

			for unitNumber := 1; unitNumber <= 100; unitNumber++ {
				carName := fmt.Sprintf("%s %s - %03d", brandData.brand, model, unitNumber)
				carID := utils.Slugify(carName)

				// Write to CSV
				if err := writer.Write([]string{carID, carName}); err != nil {
					return fmt.Errorf("failed to write row: %w", err)
				}
			}
		}
	}

	log.Printf("Car CSV export completed. Total cars: %d", len(carBrands)*5*100)
	return nil
}
