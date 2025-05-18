package usecases

import (
	"context"

	"github.com/zydhanlinnar11/sister/ec/order/domain/repositories"
)

type CompleteOrderRequest struct {
	OrderId string
}

type CompleteOrderResponse struct{}

type CompleteOrderUseCase struct {
	OrderRepo repositories.Order
}

func (o *CompleteOrderUseCase) Execute(ctx context.Context, req CompleteOrderRequest) (CompleteOrderResponse, error) {
	order, err := o.OrderRepo.FindById(ctx, req.OrderId)
	if err != nil {
		return CompleteOrderResponse{}, err
	}

	order.Complete()

	if err := o.OrderRepo.Save(ctx, order); err != nil {
		return CompleteOrderResponse{}, err
	}

	return CompleteOrderResponse{}, nil
}
