package values

import "errors"

type OrderId string

var (
	ErrInvalidOrderId = errors.New("invalid order id")
)

type OrderIdGenerator interface {
	Generate() OrderId
}

type OrderIdValidator interface {
	Validate(id string) error
}

func NewOrderId(id string, v OrderIdValidator) (OrderId, error) {
	if err := v.Validate(id); err != nil {
		return "", err
	}

	return OrderId(id), nil
}

func (o OrderId) String() string {
	return string(o)
}
