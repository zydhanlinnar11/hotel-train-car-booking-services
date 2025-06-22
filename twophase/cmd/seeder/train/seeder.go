package train

import (
	"context"
	"fmt"
	"log"

	"github.com/zydhanlinnar11/hotel-train-car-booking-services/twophase/internal/train"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/twophase/pkg/utils"
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

func Seed(ctx context.Context, repo *train.Repository) error {
	log.Println("Starting train seeder...")

	var trainSeats []train.TrainSeatTicket

	for _, trainName := range trainNames {
		log.Printf("Seeding %s...", trainName)

		// 500 seats per train
		for seatNumber := 1; seatNumber <= 500; seatNumber++ {
			seatID := fmt.Sprintf("%d", seatNumber)
			seatDocID := utils.Slugify(fmt.Sprintf("%s-%s", trainName, seatID))

			trainSeats = append(trainSeats, train.TrainSeatTicket{
				SeatID:    seatDocID,
				TrainName: trainName,
				Available: true,
			})
		}
	}

	if err := repo.BulkWriteTrainSeatTicket(ctx, trainSeats); err != nil {
		return fmt.Errorf("failed to bulk write train seats: %w", err)
	}

	log.Printf("Train seeder completed. Total seats: %d", len(trainSeats))
	return nil
}
