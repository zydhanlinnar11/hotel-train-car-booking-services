package usecases

import (
	"context"

	"github.com/zydhanlinnar11/sister/ec/order/domain/entities"
	"github.com/zydhanlinnar11/sister/ec/order/domain/repositories"
)

type ReleaseOrderRequest struct {
	OrderId string
	Cause   int
}

type ReleaseOrderResponse struct{}

type ReleaseOrderUseCase struct {
	OrderRepo repositories.Order
}

func (o *ReleaseOrderUseCase) Execute(ctx context.Context, req ReleaseOrderRequest) (ReleaseOrderResponse, error) {
	order, err := o.OrderRepo.FindById(ctx, req.OrderId)
	if err != nil {
		return ReleaseOrderResponse{}, err
	}

	order.Release(entities.NewOrderReleaseCause(req.Cause))

	if err := o.OrderRepo.Save(ctx, order); err != nil {
		return ReleaseOrderResponse{}, err
	}

	return ReleaseOrderResponse{}, nil
}
