package car

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zydhanlinnar11/hotel-train-car-booking-services/twophase/internal/car"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/twophase/pkg/config"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/twophase/pkg/utils"
)

var carBrands = []struct {
	brand  string
	models []string
}{
	{"Toyota", []string{"Avanza", "Innova", "Fortuner", "Camry", "Corolla"}},
	{"Honda", []string{"Brio", "Jazz", "HR-V", "CR-V", "Civic"}},
	{"Suzuki", []string{"Ertiga", "XL7", "Ignis", "Baleno", "Swift"}},
	{"Daihatsu", []string{"Ayla", "Calya", "Xenia", "Terios", "Rocky"}},
	{"Mitsubishi", []string{"Xpander", "Pajero", "L300", "Colt", "Mirage"}},
	{"Nissan", []string{"Livina", "Grand Livina", "X-Trail", "Serena", "March"}},
	{"Hyundai", []string{"Brio", "Creta", "Santa Fe", "Stargazer", "Palisade"}},
	{"Kia", []string{"Picanto", "Rio", "Seltos", "Sportage", "Carnival"}},
	{"Wuling", []string{"Almaz", "Cortez", "Confero", "Air ev", "Alvez"}},
	{"MG", []string{"ZS", "HS", "RX5", "5", "3"}},
}

func Seed(ctx context.Context, repo *car.Repository) error {
	log.Println("Starting car seeder...")

	var carAvailabilities []car.CarAvailability

	for _, brandData := range carBrands {
		for _, model := range brandData.models {
			log.Printf("Seeding %s %s...", brandData.brand, model)

			for unitNumber := 1; unitNumber <= 100; unitNumber++ {
				carName := fmt.Sprintf("%s %s - %03d", brandData.brand, model, unitNumber)
				carID := utils.Slugify(carName)

				date := time.Now().AddDate(0, 0, -1).Format(config.DateFormat)
				carAvailabilities = append(carAvailabilities, car.CarAvailability{
					CarID:     carID,
					CarName:   carName,
					Date:      date,
					Available: true,
				})
			}
		}
	}

	if err := repo.BulkWriteCarAvailability(ctx, carAvailabilities); err != nil {
		return fmt.Errorf("failed to bulk write car availability: %w", err)
	}

	log.Printf("Car seeder completed. Total cars: %d", len(carAvailabilities))
	return nil
}
