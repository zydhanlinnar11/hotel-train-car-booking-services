package car

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/utils"
)

type Car struct {
	ID       string `firestore:"id" json:"id"`
	CarID    string `firestore:"car_id" json:"car_id"`
	Name     string `firestore:"name" json:"name"`
	Brand    string `firestore:"brand" json:"brand"`
	Model    string `firestore:"model" json:"model"`
	Year     int    `firestore:"year" json:"year"`
	Price    int    `firestore:"price" json:"price"`
	Location string `firestore:"location" json:"location"`
}

var carBrands = []struct {
	brand     string
	models    []string
	basePrice int
	year      int
}{
	{"Toyota", []string{"Avanza", "Innova", "Fortuner", "Camry", "Corolla"}, 250000000, 2023},
	{"Honda", []string{"Brio", "Jazz", "HR-V", "CR-V", "Civic"}, 200000000, 2023},
	{"Suzuki", []string{"Ertiga", "XL7", "Ignis", "Baleno", "Swift"}, 180000000, 2023},
	{"Daihatsu", []string{"Ayla", "Calya", "Xenia", "Terios", "Rocky"}, 150000000, 2023},
	{"Mitsubishi", []string{"Xpander", "Pajero", "L300", "Colt", "Mirage"}, 220000000, 2023},
	{"Nissan", []string{"Livina", "Grand Livina", "X-Trail", "Serena", "March"}, 200000000, 2023},
	{"Hyundai", []string{"Brio", "Creta", "Santa Fe", "Stargazer", "Palisade"}, 250000000, 2023},
	{"Kia", []string{"Picanto", "Rio", "Seltos", "Sportage", "Carnival"}, 230000000, 2023},
	{"Wuling", []string{"Almaz", "Cortez", "Confero", "Air ev", "Alvez"}, 200000000, 2023},
	{"MG", []string{"ZS", "HS", "RX5", "5", "3"}, 220000000, 2023},
}

var locations = []string{
	"Jakarta Pusat", "Jakarta Selatan", "Jakarta Barat", "Jakarta Timur", "Jakarta Utara",
	"Bandung", "Surabaya", "Medan", "Semarang", "Yogyakarta",
	"Makassar", "Palembang", "Denpasar", "Manado", "Balikpapan",
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

				// Calculate price with some variation
				priceVariation := (unitNumber % 20) * 5000000 // ±50 juta variation
				price := brandData.basePrice + priceVariation

				// Rotate through locations
				location := locations[unitNumber%len(locations)]

				car := Car{
					ID:       carID,
					CarID:    carID,
					Name:     carName,
					Brand:    brandData.brand,
					Model:    model,
					Year:     brandData.year,
					Price:    price,
					Location: location,
				}

				_, err := collection.Doc(carID).Set(ctx, car)
				if err != nil {
					return fmt.Errorf("failed to seed car %s: %w", carName, err)
				}
			}
		}
	}

	log.Printf("Car seeder completed. Total cars: %d", len(carBrands)*5*100)
	return nil
}
