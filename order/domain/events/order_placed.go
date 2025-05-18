package events

import "time"

type OrderPlaced struct {
	Id                 string
	HotelRoomId        string
	HotelRoomStartDate string
	HotelRoomEndDate   string
	CarId              string
	CarStartDate       string
	CarEndDate         string
	TrainSeatId        string
	Time               time.Time
}

func NewOrderPlaced(
	id string,
	hotelRoomId string,
	hotelRoomStartDate string,
	hotelRoomEndDate string,
	carId string,
	carStartDate string,
	carEndDate string,
	trainSeatId string,
) *OrderPlaced {
	return &OrderPlaced{
		Id:                 id,
		HotelRoomId:        hotelRoomId,
		HotelRoomStartDate: hotelRoomStartDate,
		HotelRoomEndDate:   hotelRoomEndDate,
		CarId:              carId,
		CarStartDate:       carStartDate,
		CarEndDate:         carEndDate,
		TrainSeatId:        trainSeatId,
		Time:               time.Now(),
	}
}

func (o *OrderPlaced) Name() string {
	return "order_placed"
}

func (o *OrderPlaced) Timestamp() time.Time {
	return o.Time
}
