package train

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/internal/train"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/utils"
)

var trainNames = []string{
	"Argo Bromo Anggrek",
	"Argo Lawu",
	"Argo Parahyangan",
	"Bima",
	"Gajayana",
	"Harina",
	"Kertajaya",
	"Lodaya",
	"Malabar",
	"Matarmaja",
}

func Seed(ctx context.Context, client *firestore.Client) error {
	log.Println("Starting train seeder...")

	collection := client.Collection("train_seats")

	for _, trainName := range trainNames {
		log.Printf("Seeding %s...", trainName)

		// 500 seats per train (10 trains × 500 seats = 5000 total)
		for seatNumber := 1; seatNumber <= 500; seatNumber++ {
			seatID := fmt.Sprintf("%d", seatNumber)
			seatDocID := utils.Slugify(fmt.Sprintf("%s-%s", trainName, seatID))

			trainSeat := train.TrainSeat{
				ID:        seatDocID,
				SeatID:    seatID,
				TrainName: trainName,
			}

			_, err := collection.Doc(seatDocID).Set(ctx, trainSeat)
			if err != nil {
				return fmt.Errorf("failed to seed train seat %s seat %s: %w", trainName, seatID, err)
			}
		}
	}

	log.Printf("Train seeder completed. Total seats: %d", len(trainNames)*500)
	return nil
}
