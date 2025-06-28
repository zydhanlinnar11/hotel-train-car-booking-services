package main

import (
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/cmd/seeder/car"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/cmd/seeder/hotel"
	"github.com/zydhanlinnar11/hotel-train-car-booking-services/eventual/cmd/seeder/train"
)

func main() {
	car.ExportToCSV("car.csv")
	hotel.ExportToCSV("hotel.csv")
	train.ExportToCSV("train.csv")
}
