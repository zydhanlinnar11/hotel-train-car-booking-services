package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/joho/godotenv"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/twophase/pkg/config"
	"google.golang.org/api/iterator"
)

func main() {
	log.Println("Starting status resetter...")

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

	bulkWriter := client.BulkWriter(ctx)

	iter := client.Collection("twophase_car_availabilities").Where("available", "==", false).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate car availabilities: %v", err)
		}

		bulkWriter.Update(doc.Ref, []firestore.Update{
			{
				Path:  "available",
				Value: true,
			},
		})
	}

	iter = client.Collection("twophase_hotel_room_availabilities").Where("available", "==", false).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate hotel room availabilities: %v", err)
		}

		bulkWriter.Update(doc.Ref, []firestore.Update{
			{
				Path:  "available",
				Value: true,
			},
		})
	}

	iter = client.Collection("twophase_train_seat_tickets").Where("available", "==", false).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate train seat tickets: %v", err)
		}

		bulkWriter.Update(doc.Ref, []firestore.Update{
			{
				Path:  "available",
				Value: true,
			},
		})
	}

	bulkWriter.Flush()
}
