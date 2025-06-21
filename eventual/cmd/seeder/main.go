package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/joho/godotenv"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/cmd/seeder/car"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/cmd/seeder/hotel"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/cmd/seeder/train"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/config"
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
	if err := car.Seed(ctx, client); err != nil {
		log.Printf("Error seeding car data: %v", err)
		os.Exit(1)
	}

	// Run hotel seeder
	log.Println("Seeding hotel room data...")
	if err := hotel.Seed(ctx, client); err != nil {
		log.Printf("Error seeding hotel room data: %v", err)
		os.Exit(1)
	}

	// Run train seeder
	log.Println("Seeding train data...")
	if err := train.Seed(ctx, client); err != nil {
		log.Printf("Error seeding train data: %v", err)
		os.Exit(1)
	}

	log.Println("Database seeding completed successfully!")
}
