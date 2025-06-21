package hotel

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/utils"
)

type HotelRoom struct {
	ID        string `firestore:"id" json:"id"`
	HotelID   string `firestore:"hotel_id" json:"hotel_id"`
	HotelName string `firestore:"hotel_name" json:"hotel_name"`
	Name      string `firestore:"name" json:"name"`
	Price     int    `firestore:"price" json:"price"`
}

var hotelBrands = []struct {
	hotelID   string
	hotelName string
	basePrice int
}{
	{"marriott-jakarta", "Marriott Jakarta", 1500000},
	{"ritz-carlton-jakarta", "Ritz-Carlton Jakarta", 2500000},
	{"mandarin-oriental-jakarta", "Mandarin Oriental Jakarta", 2000000},
	{"four-seasons-jakarta", "Four Seasons Jakarta", 2200000},
	{"grand-hyatt-jakarta", "Grand Hyatt Jakarta", 1800000},
	{"intercontinental-jakarta", "InterContinental Jakarta", 1600000},
	{"sheraton-jakarta", "Sheraton Jakarta", 1400000},
	{"pullman-jakarta", "Pullman Jakarta", 1200000},
	{"novotel-jakarta", "Novotel Jakarta", 1000000},
	{"ibis-jakarta", "Ibis Jakarta", 800000},
	{"marriott-bandung", "Marriott Bandung", 1200000},
	{"ritz-carlton-bandung", "Ritz-Carlton Bandung", 2000000},
	{"marriott-surabaya", "Marriott Surabaya", 1100000},
	{"ritz-carlton-surabaya", "Ritz-Carlton Surabaya", 1800000},
	{"marriott-medan", "Marriott Medan", 1000000},
}

func Seed(ctx context.Context, client *firestore.Client) error {
	log.Println("Starting hotel room seeder...")

	collection := client.Collection("hotel_rooms")

	for _, hotelData := range hotelBrands {
		log.Printf("Seeding %s...", hotelData.hotelName)

		// 5 floors, 20 units per floor
		for floor := 1; floor <= 5; floor++ {
			for unitNumber := 1; unitNumber <= 20; unitNumber++ {
				roomName := fmt.Sprintf("%d%02d", floor, unitNumber)
				roomID := utils.Slugify(fmt.Sprintf("%s-%s", hotelData.hotelName, roomName))

				// Calculate price based on floor (higher floor = higher price)
				floorMultiplier := 1.0 + (float64(floor-1) * 0.1) // 10% increase per floor
				price := int(float64(hotelData.basePrice) * floorMultiplier)

				// Add some variation based on unit number
				priceVariation := (unitNumber % 10) * 50000 // ±500k variation
				price += priceVariation

				hotelRoom := HotelRoom{
					ID:        roomID,
					HotelID:   hotelData.hotelID,
					HotelName: hotelData.hotelName,
					Name:      roomName,
					Price:     price,
				}

				_, err := collection.Doc(roomID).Set(ctx, hotelRoom)
				if err != nil {
					return fmt.Errorf("failed to seed hotel room %s: %w", roomName, err)
				}
			}
		}
	}

	log.Printf("Hotel room seeder completed. Total rooms: %d", len(hotelBrands)*5*20)
	return nil
}
