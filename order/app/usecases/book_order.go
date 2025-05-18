package usecases

import (
	"context"

	"github.com/zydhanlinnar11/sister/ec/order/domain/repositories"
)

type BookOrderRequest struct {
	OrderId string
}

type BookOrderResponse struct{}

type BookOrderUseCase struct {
	OrderRepo repositories.Order
}

func (o *BookOrderUseCase) Execute(ctx context.Context, req BookOrderRequest) (BookOrderResponse, error) {
	order, err := o.OrderRepo.FindById(ctx, req.OrderId)
	if err != nil {
		return BookOrderResponse{}, err
	}

	order.Book()

	if err := o.OrderRepo.Save(ctx, order); err != nil {
		return BookOrderResponse{}, err
	}

	return BookOrderResponse{}, nil
}
