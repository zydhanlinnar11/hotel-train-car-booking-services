package usecases

import (
	"context"

	"github.com/zydhanlinnar11/sister/ec/order/domain/entities"
	"github.com/zydhanlinnar11/sister/ec/order/domain/repositories"
)

type MarkOrderAsFailureRequest struct {
	OrderId string
	Cause   int
}

type MarkOrderAsFailureResponse struct{}

type MarkOrderAsFailureUseCase struct {
	OrderRepo repositories.Order
}

func (o *MarkOrderAsFailureUseCase) Execute(ctx context.Context, req MarkOrderAsFailureRequest) (MarkOrderAsFailureResponse, error) {
	order, err := o.OrderRepo.FindById(ctx, req.OrderId)
	if err != nil {
		return MarkOrderAsFailureResponse{}, err
	}

	order.Fail(entities.NewOrderReleaseCause(req.Cause))

	if err := o.OrderRepo.Save(ctx, order); err != nil {
		return MarkOrderAsFailureResponse{}, err
	}

	return MarkOrderAsFailureResponse{}, nil
}
