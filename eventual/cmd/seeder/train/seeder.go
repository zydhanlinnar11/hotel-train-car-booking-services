package train

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/pkg/utils"
)

type TrainSeat struct {
	ID        string `firestore:"id" json:"id"`
	SeatID    string `firestore:"seat_id" json:"seat_id"`
	TrainName string `firestore:"train_name" json:"train_name"`
}

var trainNames = []string{
	"Argo Bromo Anggrek",
	"Argo Lawu",
	"Argo Parahyangan",
	"Argo Sindoro",
	"Argo Wilis",
	"Bima",
	"Gajayana",
	"Harina",
	"Kertajaya",
	"Kutojaya",
	"Lodaya",
	"Malabar",
	"Matarmaja",
	"Mutiara Selatan",
	"Purwojaya",
	"Sancaka",
	"Sembrani",
	"Senja Utama",
	"Serayu",
	"Taksaka",
}

func Seed(ctx context.Context, client *firestore.Client) error {
	log.Println("Starting train seeder...")

	collection := client.Collection("train_seats")

	for _, trainName := range trainNames {
		log.Printf("Seeding %s...", trainName)

		// 1000 seats per train
		for seatNumber := 1; seatNumber <= 1000; seatNumber++ {
			seatID := fmt.Sprintf("%d", seatNumber)
			seatDocID := utils.Slugify(fmt.Sprintf("%s-%s", trainName, seatID))

			trainSeat := TrainSeat{
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

	log.Printf("Train seeder completed. Total seats: %d", len(trainNames)*1000)
	return nil
}
