package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/joho/godotenv"
	carSeeder "github.com/zydhanlinnar11/hotel-train-car-booking-services/twophase/cmd/seeder/car"
	hotelSeeder "github.com/zydhanlinnar11/hotel-train-car-booking-services/twophase/cmd/seeder/hotel"
	trainSeeder "github.com/zydhanlinnar11/hotel-train-car-booking-services/twophase/cmd/seeder/train"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/twophase/internal/car"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/twophase/internal/hotel"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/twophase/internal/train"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/twophase/pkg/config"
)

func main() {
	log.Println("Starting database seeder...")

	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v, using system environment variables", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, cfg.GoogleProjectID)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	// Run car seeder
	log.Println("Seeding car data...")
	carRepo := car.NewRepository(client)
	if err := carSeeder.Seed(ctx, carRepo); err != nil {
		log.Printf("Error seeding car data: %v", err)
		os.Exit(1)
	}

	// Run hotel seeder
	log.Println("Seeding hotel room data...")
	hotelRepo := hotel.NewRepository(client)
	if err := hotelSeeder.Seed(ctx, hotelRepo); err != nil {
		log.Printf("Error seeding hotel room data: %v", err)
		os.Exit(1)
	}

	// Run train seeder
	log.Println("Seeding train data...")
	trainRepo := train.NewRepository(client)
	if err := trainSeeder.Seed(ctx, trainRepo); err != nil {
		log.Printf("Error seeding train data: %v", err)
		os.Exit(1)
	}

	log.Println("Database seeding completed successfully!")
}
