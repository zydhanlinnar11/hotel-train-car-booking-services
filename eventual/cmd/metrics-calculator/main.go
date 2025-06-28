package main

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/joho/godotenv"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/internal/order"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/config"
	"google.golang.org/api/iterator"
)

func main() {
	log.Println("Starting metrics calculator...")

	// Validasi argument untuk nama file export
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <output_filename.csv>", os.Args[0])
	}

	outputFilename := os.Args[1]
	if outputFilename == "" {
		log.Fatalf("Output filename cannot be empty")
	}

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

	orders := make([]order.Order, 0)

	iter := client.Collection("order_orders").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate orders: %v", err)
		}

		var o order.Order
		if err := doc.DataTo(&o); err != nil {
			log.Fatalf("Failed to convert document to order: %v", err)
		}
		orders = append(orders, o)
	}

	// Export to CSV dengan nama file yang dikustomisasi
	if err := exportOrdersToCSV(orders, outputFilename); err != nil {
		log.Fatalf("Failed to export orders to CSV: %v", err)
	}

	log.Printf("Successfully exported %d orders to %s", len(orders), outputFilename)
}

func exportOrdersToCSV(orders []order.Order, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	headers := []string{
		"ID",
		"UserID",
		"Status",
		"HotelRoomID",
		"CarID",
		"TrainSeatID",
		"HotelStartDate",
		"HotelEndDate",
		"CarStartDate",
		"CarEndDate",
		"HotelReservationID",
		"CarReservationID",
		"TrainReservationID",
		"HotelReservationStatus",
		"CarReservationStatus",
		"TrainReservationStatus",
		"HotelReservationFailureReason",
		"CarReservationFailureReason",
		"TrainReservationFailureReason",
		"CarDoneAt",
		"TrainDoneAt",
		"HotelDoneAt",
		"DoneAt",
		"CreatedAt",
		"UpdatedAt",
	}

	if err := writer.Write(headers); err != nil {
		return err
	}

	// Write data rows
	for _, o := range orders {
		row := []string{
			o.ID,
			o.UserID,
			string(o.Status),
			o.HotelRoomID,
			o.CarID,
			o.TrainSeatID,
			o.HotelStartDate,
			o.HotelEndDate,
			o.CarStartDate,
			o.CarEndDate,
			o.HotelReservationID,
			o.CarReservationID,
			o.TrainReservationID,
			string(o.HotelReservationStatus),
			string(o.CarReservationStatus),
			string(o.TrainReservationStatus),
			o.HotelReservationFailureReason,
			o.CarReservationFailureReason,
			o.TrainReservationFailureReason,
			strconv.FormatInt(formatTime(o.CarDoneAt), 10),
			strconv.FormatInt(formatTime(o.TrainDoneAt), 10),
			strconv.FormatInt(formatTime(o.HotelDoneAt), 10),
			strconv.FormatInt(formatTime(o.DoneAt), 10),
			strconv.FormatInt(formatTime(o.CreatedAt), 10),
			strconv.FormatInt(formatTime(o.UpdatedAt), 10),
		}

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func formatTime(t time.Time) int64 {
	return t.UnixMilli()
}
