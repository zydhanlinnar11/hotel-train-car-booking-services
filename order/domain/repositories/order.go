package repositories

import (
	"context"
	"errors"

	"github.com/zydhanlinnar11/sister/ec/order/domain/entities"
)

var (
	ErrOrderNotFound        = errors.New("order not found")
	ErrOrderVersionMismatch = errors.New("order version mismatch")
)

type Order interface {
	Save(ctx context.Context, order *entities.Order) error
	FindById(ctx context.Context, id string) (*entities.Order, error)
}
