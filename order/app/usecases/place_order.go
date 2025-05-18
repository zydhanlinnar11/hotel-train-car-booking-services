package usecases

import (
	"context"

	"github.com/zydhanlinnar11/sister/ec/order/domain/entities"
	"github.com/zydhanlinnar11/sister/ec/order/domain/repositories"
	"github.com/zydhanlinnar11/sister/ec/order/domain/values"
)

type PlaceOrderRequest struct {
	HotelRoomId        string
	HotelRoomStartDate string
	HotelRoomEndDate   string
	CarId              string
	CarStartDate       string
	CarEndDate         string
	TrainSeatId        string
}

type PlaceOrderResponse struct {
	OrderId string
}

type PlaceOrderUseCase struct {
	OrderRepo  repositories.Order
	OrderIdGen values.OrderIdGenerator
}

func (o *PlaceOrderUseCase) Execute(ctx context.Context, req PlaceOrderRequest) (PlaceOrderResponse, error) {
	hotelRoomStartDate, err := values.NewDate(req.HotelRoomStartDate)
	if err != nil {
		return PlaceOrderResponse{}, err
	}

	hotelRoomEndDate, err := values.NewDate(req.HotelRoomEndDate)
	if err != nil {
		return PlaceOrderResponse{}, err
	}

	carStartDate, err := values.NewDate(req.CarStartDate)
	if err != nil {
		return PlaceOrderResponse{}, err
	}

	carEndDate, err := values.NewDate(req.CarEndDate)
	if err != nil {
		return PlaceOrderResponse{}, err
	}

	order := entities.PlaceOrder(
		o.OrderIdGen,
		req.HotelRoomId,
		hotelRoomStartDate,
		hotelRoomEndDate,
		req.CarId,
		carStartDate,
		carEndDate,
		req.TrainSeatId,
	)

	if err := o.OrderRepo.Save(ctx, order); err != nil {
		return PlaceOrderResponse{}, err
	}

	return PlaceOrderResponse{
		OrderId: order.Id().String(),
	}, nil
}
