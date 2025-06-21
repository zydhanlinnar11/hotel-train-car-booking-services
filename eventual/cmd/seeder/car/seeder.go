package car

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/internal/car"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/utils"
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

func Seed(ctx context.Context, client *firestore.Client) error {
	log.Println("Starting car seeder...")

	collection := client.Collection("cars")

	for _, brandData := range carBrands {
		for _, model := range brandData.models {
			log.Printf("Seeding %s %s...", brandData.brand, model)

			for unitNumber := 1; unitNumber <= 100; unitNumber++ {
				carName := fmt.Sprintf("%s %s - %03d", brandData.brand, model, unitNumber)
				carID := utils.Slugify(carName)

				carData := car.Car{
					ID:   carID,
					Name: carName,
				}

				_, err := collection.Doc(carID).Set(ctx, carData)
				if err != nil {
					return fmt.Errorf("failed to seed car %s: %w", carName, err)
				}
			}
		}
	}

	log.Printf("Car seeder completed. Total cars: %d", len(carBrands)*5*100)
	return nil
}
